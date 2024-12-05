package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	userAPI "github.com/sSmok/auth/internal/api/user"
	"github.com/sSmok/auth/internal/closer"
	"github.com/sSmok/auth/internal/config"
	"github.com/sSmok/auth/internal/repository"
	userRepository "github.com/sSmok/auth/internal/repository/user"
	"github.com/sSmok/auth/internal/service"
	"github.com/sSmok/auth/internal/service/user"
)

type container struct {
	grpcConfig config.GRPCConfigI
	pgConfig   config.PGConfigI

	pgPool *pgxpool.Pool

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

func (c *container) PGPool(ctx context.Context) *pgxpool.Pool {
	if c.pgPool == nil {
		pool, err := pgxpool.New(ctx, c.PGConfig().DSN())
		if err != nil {
			log.Fatalf("db pool connection error: %v", err)
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})

		c.pgPool = pool
	}
	return c.pgPool
}

func (c *container) UserRepository(ctx context.Context) repository.UserRepositoryI {
	if c.userRepository == nil {
		c.userRepository = userRepository.NewUserRepository(c.PGPool(ctx))
	}
	return c.userRepository
}

func (c *container) UserService(ctx context.Context) service.UserServiceI {
	if c.userService == nil {
		c.userService = user.NewService(c.UserRepository(ctx))
	}
	return c.userService
}

func (c *container) UserAPI(ctx context.Context) *userAPI.API {
	if c.userAPI == nil {
		c.userAPI = userAPI.NewAPI(c.UserService(ctx))
	}
	return c.userAPI
}
