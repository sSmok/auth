package converter

import (
	"github.com/sSmok/auth/internal/model"
	modelRepo "github.com/sSmok/auth/internal/repository/user/model"
	descUser "github.com/sSmok/auth/pkg/user_v1"
)

// ToUserFromRepo конвертирует данные пользователя из БД для сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo конвертирует данные пользователя из БД для сервисного слоя
func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  descUser.Role_value[info.Role],
	}
}

// ToRepoFromUserInfo конвертирует данные пользователя из сервисного слоя в БД
func ToRepoFromUserInfo(info *model.UserInfo) *modelRepo.UserInfo {
	return &modelRepo.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  descUser.Role_name[info.Role],
	}
}
