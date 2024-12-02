package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/auth/internal/config"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	userIDCol        = "id"
	userNameCol      = "name"
	userEmailCol     = "email"
	userRoleCol      = "role"
	userCreatedAtCol = "created_at"
	userUpdatedAtCol = "updated_at"
	userPasswordCol  = "password"
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

	descUser.RegisterUserV1Server(serv, &server{pool: pool})
	if err = serv.Serve(lis); err != nil {
		log.Printf("fail to serve: %v\n", err)
	}
}

type server struct {
	descUser.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) CreateUser(ctx context.Context, req *descUser.CreateUserRequest) (*descUser.CreateUserResponse, error) {
	pass := req.GetPass().GetPassword()
	passConfirm := req.GetPass().GetPasswordConfirm()
	if pass != passConfirm {
		return nil, errors.New("passwords don't match")
	}

	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns(userNameCol, userEmailCol, userRoleCol, userCreatedAtCol, userUpdatedAtCol, userPasswordCol).
		Values(req.GetInfo().GetName(), req.GetInfo().GetEmail(), req.GetInfo().GetRole(), time.Now().UTC(), time.Now().UTC(), pass).
		Suffix("RETURNING id")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &descUser.CreateUserResponse{Id: userID}, nil
}

func (s *server) GetUser(ctx context.Context, req *descUser.GetUserRequest) (*descUser.GetUserResponse, error) {
	var id int64
	var name, email, role string
	var createdAt time.Time
	var updatedAt time.Time
	builder := sq.Select(userIDCol, userNameCol, userEmailCol, userRoleCol, userCreatedAtCol, userUpdatedAtCol).
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIDCol: req.GetId()}).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	resp := &descUser.GetUserResponse{
		User: &descUser.User{
			Id: id,
			Info: &descUser.UserInfo{
				Name:  name,
				Email: email,
				Role:  descUser.Role(descUser.Role_value[role]),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		},
	}

	return resp, nil
}

func (s *server) UpdateUser(ctx context.Context, req *descUser.UpdateUserRequest) (*emptypb.Empty, error) {
	builder := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set(userNameCol, req.GetInfo().GetName()).
		Set(userEmailCol, req.GetInfo().GetEmail()).
		Set(userRoleCol, req.GetInfo().GetRole()).
		Set(userUpdatedAtCol, time.Now().UTC()).
		Where(sq.Eq{userIDCol: req.GetId()})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	exec, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	log.Printf("updated %d rows", exec.RowsAffected())

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *descUser.DeleteUserRequest) (*emptypb.Empty, error) {
	builder := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIDCol: req.GetId()})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	exec, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	log.Printf("deleted %d rows", exec.RowsAffected())
	return &emptypb.Empty{}, nil
}
