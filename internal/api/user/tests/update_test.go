package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/sSmok/auth/internal/api/user"
	"github.com/sSmok/auth/internal/model"
	"github.com/sSmok/auth/internal/service"
	"github.com/sSmok/auth/internal/service/mocks"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestAPI_UpdateUser(t *testing.T) {
	type userServiceIMockFunc func(mc *minimock.Controller) service.UserServiceI

	type args struct {
		ctx context.Context
		req *descUser.UpdateUserRequest
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		name          = gofakeit.Name()
		email         = gofakeit.Email()
		roleSlice     = []int32{0, 1}
		role          = roleSlice[gofakeit.Number(0, 1)]
		serviceErr    = fmt.Errorf("service error")

		updInfo = &descUser.UpdateUserInfo{
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(email),
			Role:  descUser.Role(role),
		}

		req = &descUser.UpdateUserRequest{
			Id:   id,
			Info: updInfo,
		}

		info = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}
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
				mock.UpdateUserMock.Expect(ctx, id, info).Return(nil)
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
				mock.UpdateUserMock.Expect(ctx, id, info).Return(serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceIMock := tt.userServiceIMock(minimockContr)
			api := user.NewAPI(userServiceIMock)

			newResp, err := api.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newResp)
		})
	}
}
