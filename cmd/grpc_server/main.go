package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50500

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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
	descUser.RegisterUserV1Server(serv, &server{})
	if err = serv.Serve(lis); err != nil {
		log.Printf("fail to serve: %v\n", err)
	}
}

type server struct {
	descUser.UnimplementedUserV1Server
}

func (s *server) Create(_ context.Context, req *descUser.CreateRequest) (*descUser.CreateResponse, error) {
	pass := req.GetPass()
	user := &descUser.User{
		Id:        gofakeit.Int64(),
		Info:      req.GetInfo(),
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}

	log.Printf("user pass: %+v", pass)
	log.Printf("user created: %+v", user)

	return &descUser.CreateResponse{Id: user.Id}, nil
}
func (s *server) Get(_ context.Context, req *descUser.GetRequest) (*descUser.GetResponse, error) {
	return &descUser.GetResponse{
		User: &descUser.User{
			Id: req.GetId(),
			Info: &descUser.UserInfo{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Role:  0,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Update(_ context.Context, req *descUser.UpdateRequest) (*emptypb.Empty, error) {
	user := &descUser.User{
		Id: 5,
		Info: &descUser.UserInfo{
			Name:  gofakeit.Name(),
			Email: gofakeit.Name(),
			Role:  descUser.Role_user,
		},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	log.Printf("user before update: %+v\n", user)

	user.Info.Name = req.Info.Name.GetValue()
	user.Info.Email = req.Info.Email.GetValue()
	user.Info.Role = req.Info.Role
	user.UpdatedAt = timestamppb.Now()

	log.Printf("user after update: %+v\n", user)

	return &emptypb.Empty{}, nil
}
func (s *server) Delete(_ context.Context, req *descUser.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() == 12 {
		log.Printf("user successfully deleted")
		return &emptypb.Empty{}, nil
	}

	return nil, fmt.Errorf("user cannot be deleted")
}
