package user

import (
	"github.com/sSmok/auth/internal/service"
	descUser "github.com/sSmok/auth/pkg/user_v1"
)

type Api struct {
	descUser.UnimplementedUserV1Server
	service service.UserServiceI
}

func NewApi(userService service.UserServiceI) *Api {
	return &Api{
		service: userService,
	}
}
