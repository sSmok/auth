package converter

import (
	"github.com/sSmok/auth/internal/model"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserInfoFromDesc конвертирует proto данные пользователя в модель сервисного слоя
func ToUserInfoFromDesc(protoUserInfo *descUser.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  protoUserInfo.GetName(),
		Email: protoUserInfo.GetEmail(),
		Role:  int32(protoUserInfo.GetRole()),
	}
}

// ToUserPasswordFromDesc конвертирует proto данные паролей пользователя в модель сервисного слоя
func ToUserPasswordFromDesc(protoUserInfo *descUser.UserPassword) *model.UserPassword {
	return &model.UserPassword{
		Password:        protoUserInfo.GetPassword(),
		PasswordConfirm: protoUserInfo.GetPasswordConfirm(),
	}
}

// ToUserInfoFromDescUpdate конвертирует proto данные пользователя в модель сервисного слоя
func ToUserInfoFromDescUpdate(protoUserInfo *descUser.UpdateUserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  protoUserInfo.GetName().GetValue(),
		Email: protoUserInfo.GetEmail().GetValue(),
		Role:  int32(protoUserInfo.GetRole()),
	}
}

// ToDescFromUser конвертирует данные пользователя из сервисного слоя в proto данные
func ToDescFromUser(user *model.User) *descUser.User {
	return &descUser.User{
		Id:        user.ID,
		Info:      ToDescFromUserInfo(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToDescFromUserInfo конвертирует данные пользователя из сервисного слоя в proto данные
func ToDescFromUserInfo(userInfo model.UserInfo) *descUser.UserInfo {
	return &descUser.UserInfo{
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Role:  descUser.Role(userInfo.Role),
	}
}
