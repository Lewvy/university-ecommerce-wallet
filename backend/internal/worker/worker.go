package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
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

const (
	EmailQueueKey     = "queue:emails"
	MaxWorkers        = 50
	ScaleUpThreshold  = 20
	ScaleDownCooldown = 10 * time.Second
)

type WorkerPool struct {
	Mailer        *mailer.Mailer
	Cache         *cache.ValkeyCache
	cancelWorkers []chan struct{}
	mu            sync.Mutex
	ActiveWorkers int64
	stopSignal    chan struct{}
	Logger        *slog.Logger
	TestMode      bool
}

func NewWorkerPool(m *mailer.Mailer, c *cache.ValkeyCache, l *slog.Logger, t bool) *WorkerPool {
	return &WorkerPool{
		Mailer:        m,
		Cache:         c,
		stopSignal:    make(chan struct{}),
		cancelWorkers: make([]chan struct{}, 0),
		Logger:        l,
		TestMode:      t,
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
	p.Logger.Info("WorkerPool stopped")
}

func (p *WorkerPool) StartQueueMonitor() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		var lastScaleDown time.Time
		bgctx := context.Background()

		prevQueueLength := 0
		prevActiveWorkers := 1

		for {
			select {
			case <-p.stopSignal:
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(bgctx, 5*time.Second)
				queueLength, err := p.Cache.Client.Do(
					ctx, p.Cache.Client.B().Llen().Key(EmailQueueKey).Build()).AsInt64()
				cancel()

				active := atomic.LoadInt64(&p.ActiveWorkers)
				if prevQueueLength != int(queueLength) || prevActiveWorkers != int(active) {
					p.Logger.Info("Monitoring queue length", "length", queueLength, "current_active_workers", active)
					prevQueueLength = int(queueLength)
					prevActiveWorkers = int(active)
				}

				if err != nil {
					p.Logger.Error("Error fetching queue length", "err", err)
					continue
				}

				if queueLength > ScaleUpThreshold && active < MaxWorkers {
					workersToDeploy := 5
					p.Logger.Info("Scaling UP", "queue_size", queueLength, "current_workers", active, "deploying", workersToDeploy)
					p.StartEmailWorkers(workersToDeploy)
				}

				if queueLength <= 5 && active > 5 {
					if time.Since(lastScaleDown) < ScaleDownCooldown {
						continue
					}
					lastScaleDown = time.Now()

					p.mu.Lock()
					workersToStop := min(int(active-5), len(p.cancelWorkers))
					if workersToStop > 0 {
						channelsToClose := p.cancelWorkers[len(p.cancelWorkers)-workersToStop:]
						p.cancelWorkers = p.cancelWorkers[:len(p.cancelWorkers)-workersToStop]
						p.mu.Unlock()

						p.Logger.Info("Scaling DOWN", "queue_size", queueLength, "stopping_workers", workersToStop)
						for _, ch := range channelsToClose {
							close(ch)
						}
					} else {
						p.mu.Unlock()
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
	if p.TestMode {
		time.Sleep(2 * time.Second)
	}
	defer func() {
		atomic.AddInt64(&p.ActiveWorkers, -1)
		p.removeWorkerCancelChan(workerCancelChan)
	}()

	baseCtx, baseCancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-workerCancelChan:
			baseCancel()
			p.Logger.Info("Worker received scale down signal... Exiting", "worker_id", id)
		case <-p.stopSignal:
			baseCancel()
			p.Logger.Warn("Worker received shutdown signal", "worker_id", id)
		}
	}()

	delay := 1 * time.Second

	for {
		select {
		case <-baseCtx.Done():
			return
		default:
		}

		ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
		res := p.Cache.Client.Do(ctx, p.Cache.Client.B().Brpop().Key(EmailQueueKey).Timeout(2).Build())
		cancel()

		if err := res.Error(); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			if strings.Contains(err.Error(), "nil message") {
				continue
			}
			p.Logger.Error("Valkey BRPOP error. Retrying...", "error", err, "worker_id", id)
			time.Sleep(delay)
			if delay < 30*time.Second {
				delay *= 2
			}
			continue
		}

		elements, err := res.AsStrSlice()
		if err != nil || len(elements) < 2 {
			p.Logger.Error("Failed to extract job data from BRPOP", "error", err, "worker_id", id)
			continue
		}

		var job MailJob
		if err := json.Unmarshal([]byte(elements[1]), &job); err != nil {
			p.Logger.Error("Failed to unmarshal job JSON", "job_skipped", elements[1], "worker_id", id)
			continue
		}

		if p.TestMode {
			time.Sleep(2 * time.Second)
			p.Logger.Info("Successfully sent mail", "recipient", job.Recipient, "worker_id", id)
			continue
		}

		if err := p.Mailer.Send(job.Recipient, job.TemplateFile, job.TemplateData); err != nil {
			p.Logger.Error("Final delivery mail failed", "recipient", job.Recipient, "template", job.TemplateFile, "error", err, "worker_id", id)
		} else {
			p.Logger.Info("Successfully sent mail", "recipient", job.Recipient, "worker_id", id)
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
