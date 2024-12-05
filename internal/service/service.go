package service

import (
	"context"

	"github.com/sSmok/auth/internal/model"
)

// UserServiceI создает контракт для работы с пользователем на уровне бизнес логики
type UserServiceI interface {
	CreateUser(ctx context.Context, info *model.UserInfo, pass *model.UserPassword) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, info *model.UserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}
