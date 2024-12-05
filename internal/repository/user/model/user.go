package model

import "time"

// User описывает структуру полей пользователя для работы с БД
type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserInfo описывает структуру полей информации о пользователе для работы с БД
type UserInfo struct {
	Name  string
	Email string
	Role  string
}
