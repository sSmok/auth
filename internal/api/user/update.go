package user

import (
	"context"

	"github.com/sSmok/auth/internal/converter"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *Api) UpdateUser(ctx context.Context, req *descUser.UpdateUserRequest) (*emptypb.Empty, error) {
	err := api.service.UpdateUser(ctx, req.GetId(), converter.ToUserInfoFromDescUpdate(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}