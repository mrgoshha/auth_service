package session

import (
	"AuthenticationService/internal/adapter/dbs"
	"AuthenticationService/internal/adapter/dbs/postgres"
	"AuthenticationService/internal/adapter/dbs/postgres/entity"
	"AuthenticationService/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateSession(s *model.Session) error {
	query := `	INSERT INTO sessions (session_id, refresh_token, ip, user_id, expires_at)
				VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, s.SessionId, s.RefreshToken, s.Ip, s.UserId, s.ExpiresAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf(`create session %w`, dbs.ErrorRecordAlreadyExists)
		}
		return fmt.Errorf(`create session %w`, err)
	}

	return nil
}

func (r *Repository) UpdateSession(u *model.Session) error {
	query := `	UPDATE sessions
				SET session_id = $1, refresh_token = $2, ip = $3, user_id = $4, expires_at = $5
				WHERE id = $6`

	_, err := r.db.Exec(query, u.SessionId, u.RefreshToken, u.Ip, u.UserId, u.ExpiresAt, u.Id)

	if err != nil {
		return fmt.Errorf(`update session %w`, err)
	}

	return nil
}

func (r *Repository) GetSessionBySessionId(sessionId string) (*model.Session, error) {
	query := `	SELECT *
				FROM sessions
				WHERE session_id = $1`

	session := &entity.Session{}
	err := r.db.Get(session, query, sessionId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(`get session by sessionId %w`, dbs.ErrorRecordNotFound)
		}
		return nil, fmt.Errorf(`get session by sessionId %w`, err)
	}

	return postgres.ToSessionServiceModel(session), nil
}
