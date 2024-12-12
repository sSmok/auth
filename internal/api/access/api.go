package access

import (
	"context"

	"github.com/sSmok/auth/internal/config"
	"github.com/sSmok/auth/internal/repository"
	descAccess "github.com/sSmok/auth/pkg/access_v1"
)

const authPrefix = "Bearer "

// API - апи слой для авторизации пользователя
type API struct {
	descAccess.UnimplementedAccessV1Server
	accessRepo  repository.AccessRepositoryI
	tokenConfig config.TokenConfigI
	roles       map[string][]int32
}

// NewAPI - конструктор апи слоя
func NewAPI(accessRepo repository.AccessRepositoryI, tokenConfig config.TokenConfigI) *API {
	return &API{
		accessRepo:  accessRepo,
		tokenConfig: tokenConfig,
	}
}

// Возвращает мапу с адресом эндпоинта и ролью, которая имеет доступ к нему
func (api *API) accessibleRoles(ctx context.Context) (map[string][]int32, error) {
	if api.roles == nil {
		roles, err := api.accessRepo.GetAccessibleRoles(ctx)
		if err != nil {
			return nil, err
		}
		api.roles = roles
	}

	return api.roles, nil
}
