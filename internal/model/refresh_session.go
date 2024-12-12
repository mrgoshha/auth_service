package model

import "time"

type Session struct {
	Id int

	SessionId string

	RefreshToken string

	Ip string

	UserId string

	ExpiresAt time.Time
}
