package converter

import (
	"github.com/sSmok/auth/internal/model"
	modelRepo "github.com/sSmok/auth/internal/repository/user/model"
	descUser "github.com/sSmok/auth/pkg/user_v1"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  int(descUser.Role_value[info.Role]),
	}
}

func ToRepoFromUserInfo(info *model.UserInfo) *modelRepo.UserInfo {
	return &modelRepo.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  descUser.Role_name[int32(info.Role)],
	}
}
