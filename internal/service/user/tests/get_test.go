package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/sSmok/auth/internal/model"
	"github.com/sSmok/auth/internal/repository"
	"github.com/sSmok/auth/internal/repository/mocks"
	userService "github.com/sSmok/auth/internal/service/user"
	"github.com/sSmok/platform_common/pkg/client/db"
	txMocks "github.com/sSmok/platform_common/pkg/client/db/mocks"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/stretchr/testify/require"
)

func Test_serv_GetUser(t *testing.T) {
	type userRepositoryIMockFunc func(mc *minimock.Controller) repository.UserRepositoryI
	type transactorIMockFunc func(mc *minimock.Controller) db.TransactorI

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		name          = gofakeit.Name()
		email         = gofakeit.Email()
		roleSlice     = []int32{0, 1}
		role          = roleSlice[gofakeit.Number(0, 1)]
		repoErr       = errors.New("error")

		info = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		user = &model.User{
			ID:   id,
			Info: *info,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *model.User
		err                 error
		userRepositoryIMock userRepositoryIMockFunc
		transactorIMock     transactorIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: user,
			err:  nil,
			userRepositoryIMock: func(mc *minimock.Controller) repository.UserRepositoryI {
				mock := mocks.NewUserRepositoryIMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				return mock
			},
		},
		{
			name: "fail case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErr,
			userRepositoryIMock: func(mc *minimock.Controller) repository.UserRepositoryI {
				mock := mocks.NewUserRepositoryIMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepositoryIMock := tt.userRepositoryIMock(minimockContr)
			txManager := transaction.NewManager(tt.transactorIMock(minimockContr))
			serv := userService.NewService(userRepositoryIMock, txManager)

			u, err := serv.GetUser(tt.args.ctx, tt.args.id)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, u)
		})
	}
}
