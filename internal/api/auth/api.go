package auth

import (
	"github.com/sSmok/auth/internal/config"
	"github.com/sSmok/auth/internal/repository"
	descAuth "github.com/sSmok/auth/pkg/auth_v1"
)

// API - апи слой для работы с аутентификацией пользователя
type API struct {
	descAuth.UnimplementedAuthV1Server
	userRepo    repository.UserRepositoryI
	tokenConfig config.TokenConfigI
}

// NewAPI - конструктор апи слоя
func NewAPI(userRepo repository.UserRepositoryI, tokenConfig config.TokenConfigI) *API {
	return &API{
		userRepo:    userRepo,
		tokenConfig: tokenConfig,
	}
}
