package entity

type User struct {
	Id string `db:"user_id"`

	Email string `db:"email"`
}
