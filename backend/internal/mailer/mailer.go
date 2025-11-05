package mailer

import (
	"bytes"
	"embed"
	"errors"
	"time"

	"github.com/wneessen/go-mail"

	ht "html/template"
	tt "text/template"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (*Mailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		client: client,
		sender: sender,
	}

	return mailer, nil
}

func (m *Mailer) Send(recipient string, templateFile string, data any) error {

	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return errors.New("Error parsing file: " + err.Error())
	}

	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return errors.New("Error executing subject templ: " + err.Error())
	}

	plainBody := new(bytes.Buffer)
	if err := textTmpl.ExecuteTemplate(plainBody, "plainBody", data); err != nil {
		return errors.New("Error executing plainBody templ: " + err.Error())
	}

	htmlTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return errors.New("Error executing htmlTmpl template: " + err.Error())
	}

	htmlBody := new(bytes.Buffer)
	if err := htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data); err != nil {
		return errors.New("Error executing htmlBody template: " + err.Error())
	}

	msg := mail.NewMsg()
	if err := msg.To(recipient); err != nil {
		return errors.New("Recipient not found: " + err.Error())
	}

	if err := msg.From(m.sender); err != nil {
		return errors.New("Sender not found: " + err.Error())
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	for i := range 3 {
		err = m.client.DialAndSend(msg)
		if err == nil {
			break
		}
		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return err
}
