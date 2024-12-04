package converter

import (
	"github.com/sSmok/auth/internal/model"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserInfoFromDesc(protoUserInfo *descUser.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  protoUserInfo.GetName(),
		Email: protoUserInfo.GetEmail(),
		Role:  int(protoUserInfo.GetRole()),
	}
}

func ToUserInfoFromDescUpdate(protoUserInfo *descUser.UpdateUserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  protoUserInfo.GetName().GetValue(),
		Email: protoUserInfo.GetEmail().GetValue(),
		Role:  int(protoUserInfo.GetRole()),
	}
}

func ToDescFromUser(user *model.User) *descUser.User {
	return &descUser.User{
		Id:        user.ID,
		Info:      ToDescFromUserInfo(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToDescFromUserInfo(userInfo model.UserInfo) *descUser.UserInfo {
	return &descUser.UserInfo{
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Role:  descUser.Role(userInfo.Role),
	}
}
