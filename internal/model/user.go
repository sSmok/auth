package model

import "time"

// User описывает структуру полей пользователя для работы в приложении
type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserInfo описывает структуру полей информации о пользователе для работы в приложении
type UserInfo struct {
	Name  string
	Email string
	Role  int32
}

// UserPassword описывает структуру полей паролей пользователя для работы в приложении
type UserPassword struct {
	Password, PasswordConfirm string
}
