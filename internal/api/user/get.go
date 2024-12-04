package user

import (
	"context"

	"github.com/sSmok/auth/internal/converter"
	descUser "github.com/sSmok/auth/pkg/user_v1"
)

func (api *Api) GetUser(ctx context.Context, req *descUser.GetUserRequest) (*descUser.GetUserResponse, error) {
	userRepo, err := api.service.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &descUser.GetUserResponse{User: converter.ToDescFromUser(userRepo)}, nil
}
