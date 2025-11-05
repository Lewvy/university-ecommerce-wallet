package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"ecommerce/internal/cache"
	"ecommerce/internal/mailer"
)

type MailJob struct {
	Recipient    string `json:"recipient"`
	TemplateFile string `json:"template_file"`
	TemplateData any    `json:"template_data"`
}

const EmailQueueKey = "queue:emails"
const MaxWorkers = 50
const ScaleUpThreshold = 20

type WorkerPool struct {
	Mailer        *mailer.Mailer
	Cache         *cache.ValkeyCache
	cancelWorkers []chan struct{}
	mu            sync.Mutex
	ActiveWorkers int64
	stopSignal    chan struct{}
	Logger        *slog.Logger
}

func NewWorkerPool(m *mailer.Mailer, c *cache.ValkeyCache) *WorkerPool {
	return &WorkerPool{
		Mailer:        m,
		mu:            sync.Mutex{},
		Cache:         c,
		ActiveWorkers: 0,
		stopSignal:    make(chan struct{}),
		cancelWorkers: make([]chan struct{}, 0),
	}
}

func (p *WorkerPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case <-p.stopSignal:
		return
	default:
		close(p.stopSignal)
	}

	for _, ch := range p.cancelWorkers {
		close(ch)
	}
	p.cancelWorkers = nil
	p.Logger.Info("WorkerPool stopped cleanly.")
}

func (p *WorkerPool) StartQueueMonitor() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		monitorBaseCtx := context.Background()

		for {
			select {
			case <-p.stopSignal:
				return

			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				queueLength, err := p.Cache.Client.Do(ctx, p.Cache.Client.B().Llen().Key(EmailQueueKey).Build()).AsInt64()
				cancel()

				if err != nil {
					p.Logger.Log(monitorBaseCtx, slog.LevelError, "Error fetching length of the queue", "err", err)
					continue
				}

				currentWorkers := atomic.LoadInt64(&p.ActiveWorkers)

				if queueLength > ScaleUpThreshold && currentWorkers < MaxWorkers {
					workersToDeploy := 5

					p.Logger.Log(monitorBaseCtx, slog.LevelInfo, "Scaling UP triggered.",
						"queue_size", queueLength,
						"current_workers", currentWorkers,
						"deploying", workersToDeploy)

					p.StartEmailWorkers(workersToDeploy)
				}

				if queueLength <= 5 && currentWorkers > 5 {
					workersToStop := int(currentWorkers - 5)
					p.mu.Lock()
					channelsToClose := p.cancelWorkers[len(p.cancelWorkers)-workersToStop:]
					p.cancelWorkers = p.cancelWorkers[:len(p.cancelWorkers)-workersToStop]
					p.mu.Unlock()

					p.Logger.Info("Scaling down:", "Queue size", queueLength, "Stopping workers...", workersToStop)

					for _, ch := range channelsToClose {
						close(ch)
					}
				}

			}
		}
	}()
}

func (p *WorkerPool) StartEmailWorkers(numWorkers int) {
	for range numWorkers {
		workerCancelChan := make(chan struct{})

		p.mu.Lock()
		p.cancelWorkers = append(p.cancelWorkers, workerCancelChan)
		p.mu.Unlock()
		go p.RunEmailWorker(workerCancelChan)

	}
}

func (p *WorkerPool) RunEmailWorker(workerCancelChan chan struct{}) {
	id := atomic.AddInt64(&p.ActiveWorkers, 1)
	defer func() {
		atomic.AddInt64(&p.ActiveWorkers, -1)
		p.removeWorkerCancelChan(workerCancelChan)
	}()

	baseCtx, baseCancel := context.WithCancel(context.Background())
	defer baseCancel()
	go func() {
		select {
		case <-workerCancelChan:
			baseCancel()
			p.Logger.Info("Worker received SCALE-DOWN signal. Exiting.", "worker_id", id)
			return
		case <-p.stopSignal:
			baseCancel()
			p.Logger.Warn("Worker received shutdown signal.", "worker_id", id)
			return
		}
	}()
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	for {
		ctx, cancel := context.WithTimeout(baseCtx, 10*time.Second)
		res := p.Cache.Client.Do(ctx, p.Cache.Client.B().Brpop().Key(EmailQueueKey).Timeout(0).Build())
		cancel()

		if err := res.Error(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			p.Logger.Error("Valkey BRPOP error. Retrying...", "error", err, "worker_id", id)
			time.Sleep(backoff * time.Second)
			if backoff < maxBackoff {
				backoff *= 2
			}
			continue
		}

		elements, err := res.AsStrSlice()
		if err != nil || len(elements) < 2 {
			p.Logger.Error("Failed to extract job data from BRPOP", "error", err, "worker_id", id)
			continue
		}

		jobJSON := elements[1]

		var job MailJob
		if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
			p.Logger.Error("Failed to unmarshal job JSON.",
				"job skipped", jobJSON, "worker_id", id)
			continue
		}

		if err := p.Mailer.Send(job.Recipient, job.TemplateFile, job.TemplateData); err != nil {
			p.Logger.Error("Final delivery mail failed", "recipient",
				job.Recipient, "template", job.TemplateFile, "error", err)
		} else {
			p.Logger.Info("Successfully sent mail", "recipient", job.Recipient)
		}
	}
}

func (p *WorkerPool) removeWorkerCancelChan(c chan struct{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, ch := range p.cancelWorkers {
		if ch == c {
			p.cancelWorkers = append(p.cancelWorkers[:i], p.cancelWorkers[i+1:]...)
			break
		}
	}
}
