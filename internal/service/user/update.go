package user

import (
	"context"

	"github.com/sSmok/auth/internal/model"
)

func (s serv) UpdateUser(ctx context.Context, id int64, info *model.UserInfo) error {
	err := s.userRepository.UpdateUser(ctx, id, info)
	if err != nil {
		return err
	}

	return nil
}
