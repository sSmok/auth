package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/sSmok/auth/internal/api/user"
	"github.com/sSmok/auth/internal/service"
	"github.com/sSmok/auth/internal/service/mocks"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAPI_DeleteUser(t *testing.T) {
	type userServiceIMockFunc func(mc *minimock.Controller) service.UserServiceI

	type args struct {
		ctx context.Context
		req *descUser.DeleteUserRequest
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		serviceErr    = fmt.Errorf("service error")
		req           = &descUser.DeleteUserRequest{Id: id}
	)

	tests := []struct {
		name             string
		args             args
		want             *emptypb.Empty
		err              error
		userServiceIMock userServiceIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceIMock: func(mc *minimock.Controller) service.UserServiceI {
				mock := mocks.NewUserServiceIMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "fail case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceIMock: func(mc *minimock.Controller) service.UserServiceI {
				mock := mocks.NewUserServiceIMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceIMock := tt.userServiceIMock(minimockContr)
			api := user.NewAPI(userServiceIMock)

			newResp, err := api.DeleteUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newResp)
		})
	}
}
