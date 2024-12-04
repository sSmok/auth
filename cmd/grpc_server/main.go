package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	userApi "github.com/sSmok/auth/internal/api/user"
	"github.com/sSmok/auth/internal/config"
	userRepository "github.com/sSmok/auth/internal/repository/user"
	"github.com/sSmok/auth/internal/service/user"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Printf("fail to listen: %v\n", err)
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Fatalf("listener cannot be closed: %v", err)
		}
	}()

	serv := grpc.NewServer()
	reflection.Register(serv)
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("db pool connection error: %v", err)
	}
	defer pool.Close()

	userRepo := userRepository.NewUserRepository(pool)
	userService := user.NewService(userRepo)
	userApiServer := userApi.NewApi(userService)

	descUser.RegisterUserV1Server(serv, userApiServer)
	if err = serv.Serve(lis); err != nil {
		log.Printf("fail to serve: %v\n", err)
	}
}
