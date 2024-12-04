package user

import (
	"context"

	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *Api) DeleteUser(ctx context.Context, req *descUser.DeleteUserRequest) (*emptypb.Empty, error) {
	err := api.service.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
