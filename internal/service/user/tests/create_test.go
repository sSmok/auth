package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
	errorsUtil "github.com/pkg/errors"
	"github.com/sSmok/auth/internal/model"
	"github.com/sSmok/auth/internal/repository"
	"github.com/sSmok/auth/internal/repository/mocks"
	userService "github.com/sSmok/auth/internal/service/user"
	"github.com/sSmok/platform_common/pkg/client/db"
	txMocks "github.com/sSmok/platform_common/pkg/client/db/mocks"
	"github.com/sSmok/platform_common/pkg/client/db/pg"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/stretchr/testify/require"
)

func Test_serv_CreateUser(t *testing.T) {
	type userRepositoryIMockFunc func(mc *minimock.Controller) repository.UserRepositoryI
	type transactorIMockFunc func(mc *minimock.Controller) db.TransactorI

	type args struct {
		ctx  context.Context
		info *model.UserInfo
		pass *model.UserPassword
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
		repoErr       = errors.New("error")
		txErr         = errorsUtil.Wrap(repoErr, "failed executing code inside transaction")

		info = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		user = &model.User{
			ID:   id,
			Info: *info,
		}

		pass = &model.UserPassword{
			Password:        pswd,
			PasswordConfirm: pswd,
		}

		tx txMocks.TxMock
	)

	tests := []struct {
		name                string
		args                args
		want                int64
		err                 error
		userRepositoryIMock userRepositoryIMockFunc
		transactorIMock     transactorIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				info: info,
				pass: pass,
			},
			want: id,
			err:  nil,
			userRepositoryIMock: func(mc *minimock.Controller) repository.UserRepositoryI {
				ctxTx := pg.MakeContextTransaction(ctx, &tx)
				mock := mocks.NewUserRepositoryIMock(mc)
				mock.CreateUserMock.Expect(ctxTx, info, pswd).Return(id, nil)
				mock.GetUserMock.Expect(ctxTx, id).Return(user, nil)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				mock.BeginTxMock.Expect(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Return(&tx, nil)
				return mock
			},
		},
		{
			name: "fail case",
			args: args{
				ctx:  ctx,
				info: info,
				pass: pass,
			},
			want: 0,
			err:  txErr,
			userRepositoryIMock: func(mc *minimock.Controller) repository.UserRepositoryI {
				ctxTx := pg.MakeContextTransaction(ctx, &tx)
				mock := mocks.NewUserRepositoryIMock(mc)
				mock.CreateUserMock.Expect(ctxTx, info, pswd).Return(0, repoErr)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				mock.BeginTxMock.Expect(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Return(&tx, nil)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepositoryIMock := tt.userRepositoryIMock(minimockContr)
			txManager := transaction.NewManager(tt.transactorIMock(minimockContr))
			serv := userService.NewService(userRepositoryIMock, txManager)

			userID, err := serv.CreateUser(tt.args.ctx, tt.args.info, tt.args.pass)

			if tt.name == "success case" {
				require.Equal(t, tt.err, err)
			} else {
				require.Contains(t, err.Error(), tt.err.Error())
			}
			require.Equal(t, tt.want, userID)
		})
	}
}
