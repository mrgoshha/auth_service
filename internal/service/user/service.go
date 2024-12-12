package user

import (
	"AuthenticationService/internal/adapter/dbs"
	"AuthenticationService/internal/model"
)

type Service struct {
	repository dbs.UserRepository
}

func NewService(r dbs.UserRepository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) GetUserById(id string) (*model.User, error) {
	user, err := s.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
