package model

import "time"

type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserInfo struct {
	Name  string
	Email string
	Role  int
}

type UserPassword struct {
	Password, PasswordConfirm string
}
