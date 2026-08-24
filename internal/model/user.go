package model

import "time"

type User struct {
	Id int64
	Login string
	PasswordHash string
	Role string
	CreatedAt time.Time
}
