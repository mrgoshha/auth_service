package service

import "AuthenticationService/internal/model"

//go:generate mockgen -source=user_service.go -destination=mocks/user_service.go

type User interface {
	GetUserById(id string) (*model.User, error)
}
