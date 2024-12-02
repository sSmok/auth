package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"github.com/sSmok/auth/internal/config"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to load grpc config: %v", err)
	}

	conn, err := grpc.NewClient(grpcConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v\n", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("listener cannot be closed: %v", err)
		}
	}()

	client := descUser.NewUserV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var userID int64
	//=============
	pass := gofakeit.Name()
	reqCreate := &descUser.CreateUserRequest{
		Info: &descUser.UserInfo{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
			Role:  descUser.Role_admin,
		},
		Pass: &descUser.UserPassword{
			Password:        pass,
			PasswordConfirm: pass,
		},
	}
	respCreate, err := client.CreateUser(ctx, reqCreate)
	if err != nil {
		log.Fatalf("failed to create user: %v\n", err)
	}
	userID = respCreate.GetId()
	log.Printf(color.RedString("User id after create: %v\n"), color.GreenString("%+v", userID))

	//=============
	resp, err := client.GetUser(ctx, &descUser.GetUserRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed on get request: %v\n", err)
	}
	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", resp.GetUser()))

	//=============
	reqUpd := &descUser.UpdateUserRequest{
		Id: userID,
		Info: &descUser.UpdateUserInfo{
			Name:  wrapperspb.String(gofakeit.Name()),
			Email: wrapperspb.String(gofakeit.Email()),
			Role:  descUser.Role_user,
		},
	}
	respUpd, err := client.UpdateUser(ctx, reqUpd)
	if err != nil {
		log.Fatalf("error at update user: %+v\n", err)
	}
	log.Printf("update successfully: %+v\n", respUpd)

	//=============
	reqDel := &descUser.DeleteUserRequest{Id: userID}
	respDel, err := client.DeleteUser(ctx, reqDel)
	if err != nil {
		log.Fatalf("error at delete user: %+v\n", err)
	}
	log.Printf("delete successfully: %+v\n", respDel)
}
