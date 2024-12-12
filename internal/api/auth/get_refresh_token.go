package auth

import (
	"context"

	"github.com/sSmok/auth/internal/model"
	"github.com/sSmok/auth/internal/utils"
	descAuth "github.com/sSmok/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken - метод для обновления токена
func (api *API) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(api.tokenConfig.RefreshTokenSecretKey()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	user, err := api.userRepo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Email: claims.Email,
		Role:  user.Info.Role,
	},
		[]byte(api.tokenConfig.RefreshTokenSecretKey()),
		api.tokenConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return nil, err
	}

	return &descAuth.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
