package user

import (
	"context"

	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser удаляет пользователя получая данные из proto объекта
func (api *API) DeleteUser(ctx context.Context, req *descUser.DeleteUserRequest) (*emptypb.Empty, error) {
	err := api.service.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
