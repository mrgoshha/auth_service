package postgres

import (
	"AuthenticationService/internal/adapter/dbs/postgres/entity"
	"AuthenticationService/internal/model"
)

func ToSessionServiceModel(s *entity.Session) *model.Session {
	return &model.Session{
		Id:           s.Id,
		SessionId:    s.SessionId,
		RefreshToken: s.Token,
		Ip:           s.Ip,
		UserId:       s.UserId,
		ExpiresAt:    s.ExpiresAt,
	}
}

func ToUserServiceModel(u *entity.User) *model.User {
	return &model.User{
		Id:    u.Id,
		Email: u.Email,
	}
}
