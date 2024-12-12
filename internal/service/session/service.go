package session

import (
	"AuthenticationService/internal/adapter/dbs"
	"AuthenticationService/internal/service"
	"AuthenticationService/pkg/auth"
	emailsender "AuthenticationService/pkg/email_sender"
	"time"
)

type Service struct {
	repository      dbs.RefreshSessionRepository
	tokenManager    auth.TokenManager
	refreshTokenTTL time.Duration
	userService     service.User
	sender          emailsender.EmailSender
}

func NewService(
	r dbs.RefreshSessionRepository,
	t auth.TokenManager,
	refreshTokenTTL time.Duration,
	userService service.User,
	sender emailsender.EmailSender,
) *Service {
	return &Service{
		repository:      r,
		tokenManager:    t,
		refreshTokenTTL: refreshTokenTTL,
		userService:     userService,
		sender:          sender,
	}
}
