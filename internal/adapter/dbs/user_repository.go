package dbs

import "AuthenticationService/internal/model"

type UserRepository interface {
	GetUserById(id string) (*model.User, error)
}
