package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"ecommerce/internal/mailer"
	"ecommerce/internal/worker"

	"github.com/valkey-io/valkey-go"
)

type MailJob struct {
	Recipient    string `json:"recipient"`
	TemplateFile string `json:"template_file"`
	TemplateData any    `json:"template_data"`
}

const EmailQueueKey = "queue:emails"

func RunEmailWorker(m *mailer.Mailer, valkeyClient valkey.Client) {
	for {
		ctx := context.Background()

		res := valkeyClient.Do(ctx,
			valkeyClient.B().Bzpopmax().Key(EmailQueueKey).Timeout(0).Build())

		if err := res.Error(); err != nil {
			log.Printf("Valkey BZPOP error: %v. Retrying in 1s...", err)
			time.Sleep(1 * time.Second)
			continue
		}

		elements, err := res.AsStrSlice()
		if err != nil || len(elements) < 2 {
			log.Printf("Failed to extract job data from BRPOP: %v", err)
			continue
		}

		jobJSON := elements[1]

		var job worker.MailJob
		if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
			log.Printf("ERROR: Failed to unmarshal job JSON. Job skipped: %s", jobJSON)
			continue
		}

		if err := m.Send(job.Recipient, job.TemplateFile, job.TemplateData); err != nil {
			log.Printf("FINAL MAIL DELIVERY FAILED for %s (%s): %v",
				job.Recipient, job.TemplateFile, err)
		} else {
			log.Printf("Successfully sent mail to %s", job.Recipient)
		}
	}
}
