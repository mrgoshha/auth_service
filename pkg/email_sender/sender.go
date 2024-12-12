package emailsender

import (
	"fmt"
	"net/smtp"
)

//go:generate mockgen -source=sender.go -destination=mocks/mock_sender.go

type EmailSender interface {
	SendEmail(toEmail string, message string) error
}

type Sender struct {
	fromEmail string

	auth smtp.Auth
	addr string
}

func NewSender(cfg *Config) *Sender {
	auth := smtp.PlainAuth("", cfg.FromEmail, cfg.FromPassword, cfg.SmtpHost)
	addr := fmt.Sprintf("%s:%s", cfg.SmtpHost, cfg.SmtpPort)
	return &Sender{
		auth: auth,
		addr: addr,
	}
}

func (s *Sender) SendEmail(toEmail string, message string) error {
	err := smtp.SendMail(s.addr, s.auth, s.fromEmail, []string{toEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("send email %w", err)
	}
	return nil
}
