package worker

import (
// "time"
)

type MailJob struct {
	Recipient string `json:"recipient"`

	TemplateFile string `json:"template_file"`

	TemplateData any `json:"template_data"`
}
