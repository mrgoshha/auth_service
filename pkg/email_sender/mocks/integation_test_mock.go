package mock_emailsender

import "github.com/stretchr/testify/mock"

type EmailSender struct {
	mock.Mock
}

func (m *EmailSender) SendEmail(toEmail, message string) error {
	return nil
}
