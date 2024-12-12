package app

import (
	"context"
	"log"

	accessAPI "github.com/sSmok/auth/internal/api/access"
	authAPI "github.com/sSmok/auth/internal/api/auth"
	userAPI "github.com/sSmok/auth/internal/api/user"
	internalConfig "github.com/sSmok/auth/internal/config"
	"github.com/sSmok/auth/internal/repository"
	accessRepository "github.com/sSmok/auth/internal/repository/access"
	userRepository "github.com/sSmok/auth/internal/repository/user"
	"github.com/sSmok/auth/internal/service"
	"github.com/sSmok/auth/internal/service/user"
	"github.com/sSmok/platform_common/pkg/client/db"
	"github.com/sSmok/platform_common/pkg/client/db/pg"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/sSmok/platform_common/pkg/closer"
	"github.com/sSmok/platform_common/pkg/config"
)

type container struct {
	grpcConfig       config.GRPCConfigI
	pgConfig         config.PGConfigI
	dbClient         db.ClientI
	tokenConfig      internalConfig.TokenConfigI
	txManager        db.TxManagerI
	userRepository   repository.UserRepositoryI
	accessRepository repository.AccessRepositoryI
	userService      service.UserServiceI
	userAPI          *userAPI.API
	accessAPI        *accessAPI.API
	authAPI          *authAPI.API
}

func newContainer() *container {
	return &container{}
}

func (c *container) GRPCConfig() config.GRPCConfigI {
	if c.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		c.grpcConfig = cfg
	}
	return c.grpcConfig
}

func (c *container) PGConfig() config.PGConfigI {
	if c.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		c.pgConfig = cfg
	}
	return c.pgConfig
}

func (c *container) TokenConfig() internalConfig.TokenConfigI {
	if c.tokenConfig == nil {
		cfg, err := internalConfig.NewTokenConfig()
		if err != nil {
			log.Fatalf("failed to get token config: %v", err)
		}
		c.tokenConfig = cfg
	}
	return c.tokenConfig
}

func (c *container) ClientDB(ctx context.Context) db.ClientI {
	if c.dbClient == nil {
		client, err := pg.NewPGClient(ctx, c.PGConfig().DSN())
		if err != nil {
			log.Fatalf("db client connection error: %v", err)
		}
		closer.Add(client.Close)
		c.dbClient = client
	}
	return c.dbClient
}

func (c *container) TxManager(ctx context.Context) db.TxManagerI {
	if c.txManager == nil {
		c.txManager = transaction.NewManager(c.ClientDB(ctx).DB())
	}

	return c.txManager
}

func (c *container) UserRepository(ctx context.Context) repository.UserRepositoryI {
	if c.userRepository == nil {
		c.userRepository = userRepository.NewUserRepository(c.ClientDB(ctx))
	}
	return c.userRepository
}

func (c *container) AccessRepository(ctx context.Context) repository.AccessRepositoryI {
	if c.accessRepository == nil {
		c.accessRepository = accessRepository.NewAccessRepository(c.ClientDB(ctx))
	}
	return c.accessRepository
}

func (c *container) UserService(ctx context.Context) service.UserServiceI {
	if c.userService == nil {
		c.userService = user.NewService(c.UserRepository(ctx), c.TxManager(ctx))
	}
	return c.userService
}

func (c *container) UserAPI(ctx context.Context) *userAPI.API {
	if c.userAPI == nil {
		c.userAPI = userAPI.NewAPI(c.UserService(ctx))
	}
	return c.userAPI
}

func (c *container) AccessAPI(ctx context.Context) *accessAPI.API {
	if c.accessAPI == nil {
		c.accessAPI = accessAPI.NewAPI(c.AccessRepository(ctx), c.TokenConfig())
	}
	return c.accessAPI
}

func (c *container) AuthAPI(ctx context.Context) *authAPI.API {
	if c.authAPI == nil {
		c.authAPI = authAPI.NewAPI(c.UserRepository(ctx), c.TokenConfig())
	}
	return c.authAPI
}
