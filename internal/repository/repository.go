package repository

import (
	"context"

	"github.com/sSmok/auth/internal/model"
)

type UserRepositoryI interface {
	CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, info *model.UserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}
