package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	address := fmt.Sprintf("localhost:%v", os.Getenv("GRPC_PORT"))
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	reqCreate := &descUser.CreateRequest{
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
	respCreate, err := client.Create(ctx, reqCreate)
	if err != nil {
		log.Fatalf("failed to create user: %v\n", err)
	}
	userID = respCreate.GetId()
	log.Printf(color.RedString("User id after create: %v\n"), color.GreenString("%+v", userID))

	//=============
	resp, err := client.Get(ctx, &descUser.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed on get request: %v\n", err)
	}
	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", resp.GetUser()))

	//=============
	reqUpd := &descUser.UpdateRequest{
		Id: userID,
		Info: &descUser.UpdateUserInfo{
			Name:  wrapperspb.String(gofakeit.Name()),
			Email: wrapperspb.String(gofakeit.Email()),
			Role:  descUser.Role_user,
		},
	}
	respUpd, err := client.Update(ctx, reqUpd)
	if err != nil {
		log.Fatalf("error at update user: %+v\n", err)
	}
	log.Printf("update successfully: %+v\n", respUpd)

	//=============
	reqDel := &descUser.DeleteRequest{Id: userID}
	respDel, err := client.Delete(ctx, reqDel)
	if err != nil {
		log.Fatalf("error at delete user: %+v\n", err)
	}
	log.Printf("delete successfully: %+v\n", respDel)
}
