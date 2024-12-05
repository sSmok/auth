package model

import "time"

// User описывает структуру полей пользователя для работы с БД
type User struct {
	ID        int64     `db:"id"`
	Info      UserInfo  `db:""`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserInfo описывает структуру полей информации о пользователе для работы с БД
type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}
