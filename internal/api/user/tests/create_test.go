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
)

func TestAPI_CreateUser(t *testing.T) {
	type userServiceIMockFunc func(mc *minimock.Controller) service.UserServiceI

	type args struct {
		ctx context.Context
		req *descUser.CreateUserRequest
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		name          = gofakeit.Name()
		email         = gofakeit.Email()
		roleSlice     = []int32{0, 1}
		role          = roleSlice[gofakeit.Number(0, 1)]
		pswd          = gofakeit.Password(true, true, true, true, false, 10)
		serviceErr    = fmt.Errorf("service error")

		req = &descUser.CreateUserRequest{
			Info: &descUser.UserInfo{
				Name:  name,
				Email: email,
				Role:  descUser.Role(role),
			},
			Pass: &descUser.UserPassword{
				Password:        pswd,
				PasswordConfirm: pswd,
			},
		}

		info = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		pass = &model.UserPassword{
			Password:        pswd,
			PasswordConfirm: pswd,
		}

		resp = &descUser.CreateUserResponse{
			Id: id,
		}
	)

	tests := []struct {
		name             string
		args             args
		want             *descUser.CreateUserResponse
		err              error
		userServiceIMock userServiceIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: resp,
			err:  nil,
			userServiceIMock: func(mc *minimock.Controller) service.UserServiceI {
				mock := mocks.NewUserServiceIMock(mc)
				mock.CreateUserMock.Expect(ctx, info, pass).Return(id, nil)
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
				mock.CreateUserMock.Expect(ctx, info, pass).Return(0, serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceIMock := tt.userServiceIMock(minimockContr)
			api := user.NewAPI(userServiceIMock)

			newResp, err := api.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newResp)
		})
	}
}
