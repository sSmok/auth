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

	userID, err := s.userRepository.CreateUser(ctx, info, pass.Password)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
