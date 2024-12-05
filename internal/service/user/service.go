package user

import (
	"github.com/sSmok/auth/internal/repository"
	"github.com/sSmok/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepositoryI
}

// NewService создает объект сервиса для работы на уровне бизнес логики
func NewService(userRepo repository.UserRepositoryI) service.UserServiceI {
	return &serv{userRepository: userRepo}
}
