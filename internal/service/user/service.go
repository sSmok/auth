package user

import (
	"github.com/sSmok/auth/internal/repository"
	"github.com/sSmok/auth/internal/service"
	"github.com/sSmok/platform_common/pkg/client/db"
)

type serv struct {
	userRepository repository.UserRepositoryI
	txManager      db.TxManagerI
}

// NewService создает объект сервиса для работы на уровне бизнес логики
func NewService(userRepo repository.UserRepositoryI, txManager db.TxManagerI) service.UserServiceI {
	return &serv{
		userRepository: userRepo,
		txManager:      txManager,
	}
}
