package user

import (
	"context"

	"github.com/sSmok/auth/internal/converter"
	descUser "github.com/sSmok/auth/pkg/user_v1"
)

// CreateUser создает нового пользователя получая данные из proto объекта
func (api *API) CreateUser(ctx context.Context, req *descUser.CreateUserRequest) (*descUser.CreateUserResponse, error) {
	userID, err := api.service.CreateUser(ctx, converter.ToUserInfoFromDesc(req.GetInfo()), converter.ToUserPasswordFromDesc(req.GetPass()))
	if err != nil {
		return nil, err
	}
	return &descUser.CreateUserResponse{Id: userID}, nil
}
