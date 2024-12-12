package emailsender

import (
	"fmt"
	"os"
)

const (
	smtpHost     = "SMTP_HOST"
	smtpPort     = "SMTP_PORT"
	fromEmail    = "FROM_EMAIL"
	fromPassword = "FROM_PASSWORD"
)

type Config struct {
	SmtpHost     string
	SmtpPort     string
	FromEmail    string
	FromPassword string
}

func NewConfig() (*Config, error) {
	host := os.Getenv(smtpHost)
	if len(host) == 0 {
		return nil, fmt.Errorf("config error")
	}
	port := os.Getenv(smtpPort)
	if len(port) == 0 {
		return nil, fmt.Errorf("config error")
	}
	email := os.Getenv(fromEmail)
	if len(email) == 0 {
		return nil, fmt.Errorf("config error")
	}
	password := os.Getenv(fromPassword)
	if len(password) == 0 {
		return nil, fmt.Errorf("config error")
	}

	return &Config{
		SmtpHost:     host,
		SmtpPort:     port,
		FromEmail:    email,
		FromPassword: password,
	}, nil
}
