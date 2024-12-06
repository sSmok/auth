package user

import (
	"context"
	"errors"

	"github.com/sSmok/auth/internal/model"
)

func (s serv) CreateUser(ctx context.Context, info *model.UserInfo, pass *model.UserPassword) (int64, error) {
	if pass.Password != pass.PasswordConfirm {
		return 0, errors.New("passwords don't match")
	}

	var userID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		userID, errTx = s.userRepository.CreateUser(ctx, info, pass.Password)
		if errTx != nil {
			return errTx
		}
		_, errTx = s.userRepository.GetUser(ctx, userID)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}
