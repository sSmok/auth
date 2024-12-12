package repository

import (
	"context"

	"github.com/sSmok/auth/internal/model"
)

// UserRepositoryI создает контракт для работы с репозиторием пользователя
type UserRepositoryI interface {
	CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, info *model.UserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}

// AccessRepositoryI создает контракт для работы с репозиторием доступа
type AccessRepositoryI interface {
	GetAccessibleRoles(ctx context.Context) (map[string][]int32, error)
}
