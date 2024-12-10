package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/sSmok/auth/internal/repository"
	"github.com/sSmok/auth/internal/repository/mocks"
	userService "github.com/sSmok/auth/internal/service/user"
	"github.com/sSmok/platform_common/pkg/client/db"
	txMocks "github.com/sSmok/platform_common/pkg/client/db/mocks"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/stretchr/testify/require"
)

func Test_serv_DeleteUser(t *testing.T) {
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
		repoErr       = errors.New("error")
	)

	tests := []struct {
		name                string
		args                args
		want                error
		userRepositoryIMock userRepositoryIMockFunc
		transactorIMock     transactorIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			userRepositoryIMock: func(mc *minimock.Controller) repository.UserRepositoryI {
				mock := mocks.NewUserRepositoryIMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
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
			want: repoErr,
			userRepositoryIMock: func(mc *minimock.Controller) repository.UserRepositoryI {
				mock := mocks.NewUserRepositoryIMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(repoErr)
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
			err := serv.DeleteUser(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.want, err)
		})
	}
}
