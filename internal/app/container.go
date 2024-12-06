package app

import (
	"context"
	"log"

	userAPI "github.com/sSmok/auth/internal/api/user"
	"github.com/sSmok/auth/internal/client/db"
	"github.com/sSmok/auth/internal/client/db/pg"
	"github.com/sSmok/auth/internal/client/db/transaction"
	"github.com/sSmok/auth/internal/closer"
	"github.com/sSmok/auth/internal/config"
	"github.com/sSmok/auth/internal/repository"
	userRepository "github.com/sSmok/auth/internal/repository/user"
	"github.com/sSmok/auth/internal/service"
	"github.com/sSmok/auth/internal/service/user"
)

type container struct {
	grpcConfig     config.GRPCConfigI
	pgConfig       config.PGConfigI
	dbClient       db.ClientI
	txManager      db.TxManagerI
	userRepository repository.UserRepositoryI
	userService    service.UserServiceI
	userAPI        *userAPI.API
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
