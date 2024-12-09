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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAPI_GetUser(t *testing.T) {
	type userServiceIMockFunc func(mc *minimock.Controller) service.UserServiceI

	type args struct {
		ctx context.Context
		req *descUser.GetUserRequest
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		name          = gofakeit.Name()
		email         = gofakeit.Email()
		roleSlice     = []int32{0, 1}
		role          = roleSlice[gofakeit.Number(0, 1)]
		createdAt     = gofakeit.Date()
		updatedAt     = createdAt
		serviceErr    = fmt.Errorf("service error")

		req = &descUser.GetUserRequest{Id: id}

		info = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		userModel = &model.User{
			ID:        id,
			Info:      *info,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		resp = &descUser.GetUserResponse{User: &descUser.User{
			Id: id,
			Info: &descUser.UserInfo{
				Name:  name,
				Email: email,
				Role:  descUser.Role(role),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}}
	)

	tests := []struct {
		name             string
		args             args
		want             *descUser.GetUserResponse
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
				mock.GetUserMock.Expect(ctx, id).Return(userModel, nil)
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
				mock.GetUserMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceIMock := tt.userServiceIMock(minimockContr)
			api := user.NewAPI(userServiceIMock)

			newResp, err := api.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newResp)
		})
	}
}
