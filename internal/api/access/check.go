package access

import (
	"context"
	"errors"
	"strings"

	"github.com/sSmok/auth/internal/utils"
	descAccess "github.com/sSmok/auth/pkg/access_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check - метод для проверки доступа к эндпоинту
func (api *API) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(api.tokenConfig.AccessTokenSecretKey()))
	if err != nil {
		return nil, errors.New("access token is invalid")
	}

	accessibleRolesMap, err := api.accessibleRoles(ctx)
	if err != nil {
		return nil, errors.New("failed to get accessible roles")
	}

	for _, role := range accessibleRolesMap[req.GetEndpointAddress()] {
		if role == claims.Role {
			return &emptypb.Empty{}, nil
		}
	}

	return nil, errors.New("access denied")
}
