package dbs

import "AuthenticationService/internal/model"

//go:generate mockgen -source=session_repository.go -destination=mocks/mock_session_repository.go

type RefreshSessionRepository interface {
	CreateSession(session *model.Session) error
	UpdateSession(update *model.Session) error
	GetSessionBySessionId(session string) (*model.Session, error)
}
