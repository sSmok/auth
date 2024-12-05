package user

import (
	"github.com/sSmok/auth/internal/service"
	descUser "github.com/sSmok/auth/pkg/user_v1"
)

// API - апи слой для работы с пользователем, взаимодействует с сервисным слоем
type API struct {
	descUser.UnimplementedUserV1Server
	service service.UserServiceI
}

// NewAPI - конструктор апи слоя
func NewAPI(userService service.UserServiceI) *API {
	return &API{
		service: userService,
	}
}
