package user

import (
	"AuthenticationService/internal/adapter/dbs"
	"AuthenticationService/internal/adapter/dbs/postgres"
	"AuthenticationService/internal/adapter/dbs/postgres/entity"
	"AuthenticationService/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserById(id string) (*model.User, error) {
	query := `	SELECT *
				FROM users
				WHERE user_id = $1`

	user := &entity.User{}
	err := r.db.Get(user, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(`get user by id %w`, dbs.ErrorRecordNotFound)
		}
		return nil, fmt.Errorf(`get user by id %w`, err)
	}

	return postgres.ToUserServiceModel(user), nil
}
