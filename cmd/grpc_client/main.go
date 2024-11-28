package main

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"time"
)

const address = "localhost:50500"

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v\n", err)
	}
	defer conn.Close()

	client := descUser.NewUserV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	//=============
	resp, err := client.Get(ctx, &descUser.GetRequest{Id: 5})
	if err != nil {
		log.Fatalf("failed to get request: %v\n", err)
	}
	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", resp.GetUser()))

	//=============
	reqCreate := &descUser.CreateRequest{
		Info: &descUser.UserInfo{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
			Role:  descUser.Role_admin,
		},
		Pass: &descUser.UserPassword{
			Password:        gofakeit.Name(),
			PasswordConfirm: gofakeit.Name(),
		},
	}
	respCreate, err := client.Create(ctx, reqCreate)
	if err != nil {
		log.Fatalf("failed to crate user: %v\n", err)
	}
	log.Printf(color.RedString("User id after create: %v\n"), color.GreenString("%+v", respCreate.GetId()))

	//=============
	reqUpd := &descUser.UpdateRequest{
		Id: 5,
		Info: &descUser.UpdateUserInfo{
			Name:  wrapperspb.String(gofakeit.Name()),
			Email: wrapperspb.String(gofakeit.Email()),
			Role:  descUser.Role_admin,
		},
	}
	respUpd, err := client.Update(ctx, reqUpd)
	if err != nil {
		log.Fatalf("error at update user: %+v\n", err)
	}
	log.Printf("update successfully: %+v\n", respUpd)

	//=============
	reqDel := &descUser.DeleteRequest{Id: 10}
	respDel, err := client.Delete(ctx, reqDel)
	if err != nil {
		log.Fatalf("error at delete user: %+v\n", err)
	}
	log.Printf("delete successfully: %+v\n", respDel)
}
