package smtp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

//go:embed templates/auth.html
var authTemplate embed.FS

type AuthEmailData struct {
	Subject string `json:"subject"`
	Otp     string `json:"otp"`
}

type mailerProviderImpl struct {
	cfg  config.MailerConfig
	auth smtp.Auth
}

func NewMailerProvider(cfg config.MailerConfig) port.MailerProvider {
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	return &mailerProviderImpl{
		cfg,
		auth,
	}
}

func (p *mailerProviderImpl) send(to, subject, body string) error {
	msg := fmt.Appendf(nil, "Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", subject, body)
	addr := fmt.Sprintf("%s:%d", p.cfg.Host, p.cfg.Port)
	return smtp.SendMail(addr, p.auth, p.cfg.User, []string{to}, msg)
}

func (p *mailerProviderImpl) AuthEmail(to, subject, otp string) error {
	tmpl, err := template.ParseFS(authTemplate, "templates/auth.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := AuthEmailData{
		Subject: subject,
		Otp:     otp,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return p.send(to, subject, body.String())
}
