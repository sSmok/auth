package main

import (
	"context"
	"flag"
	"log"

	"github.com/sSmok/auth/internal/app"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	newApp, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("fail to init application: %v", err)
	}
	err = newApp.Run()
	if err != nil {
		log.Fatalf("fail to run application: %v", err)
	}
}
