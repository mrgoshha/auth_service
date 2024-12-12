package entity

import "time"

type Session struct {
	Id int `db:"id"`

	SessionId string `db:"session_id"`

	Token string `db:"refresh_token"`

	Ip string `db:"ip"`

	UserId string `db:"user_id"`

	ExpiresAt time.Time `db:"expires_at"`
}
