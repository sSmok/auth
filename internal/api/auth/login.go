package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sSmok/auth/internal/model"
	"github.com/sSmok/auth/internal/utils"
	descAuth "github.com/sSmok/auth/pkg/auth_v1"
)

// Login - метод для аутентификации пользователя
func (api *API) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	user, err := api.userRepo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Email: req.GetEmail(),
		Role:  user.Info.Role,
	},
		[]byte(api.tokenConfig.RefreshTokenSecretKey()),
		api.tokenConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate token")
	}

	return &descAuth.LoginResponse{RefreshToken: refreshToken}, nil
}
