package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/auth/internal/config"
	"github.com/sSmok/auth/internal/converter"
	"github.com/sSmok/auth/internal/repository"
	userRepository "github.com/sSmok/auth/internal/repository/user"
	"github.com/sSmok/auth/internal/service"
	"github.com/sSmok/auth/internal/service/user"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
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
	descUser.RegisterUserV1Server(serv, &server{repo: userRepo, service: userService})
	if err = serv.Serve(lis); err != nil {
		log.Printf("fail to serve: %v\n", err)
	}
}

type server struct {
	descUser.UnimplementedUserV1Server
	repo    repository.UserRepositoryI
	service service.UserServiceI
}

func (s *server) CreateUser(ctx context.Context, req *descUser.CreateUserRequest) (*descUser.CreateUserResponse, error) {
	userID, err := s.service.CreateUser(ctx, converter.ToUserInfoFromDesc(req.GetInfo()), converter.ToUserPasswordFromDesc(req.GetPass()))
	if err != nil {
		return nil, err
	}
	return &descUser.CreateUserResponse{Id: userID}, nil
}

func (s *server) GetUser(ctx context.Context, req *descUser.GetUserRequest) (*descUser.GetUserResponse, error) {
	userRepo, err := s.service.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &descUser.GetUserResponse{User: converter.ToDescFromUser(userRepo)}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *descUser.UpdateUserRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateUser(ctx, req.GetId(), converter.ToUserInfoFromDescUpdate(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *descUser.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
