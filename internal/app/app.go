package app

import (
	"context"
	"flag"
	"net"

	"github.com/sSmok/auth/internal/closer"
	"github.com/sSmok/auth/internal/config"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.Parse()
}

type App struct {
	container  *container
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}
	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return app.runGRPCSever(ctx)
}

func (app *App) initDeps(ctx context.Context) error {
	deps := []func(context.Context) error{
		app.initConfig,
		app.initContainer,
		app.initGRPCSever,
	}

	for _, f := range deps {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) initContainer(ctx context.Context) error {
	app.container = newContainer()
	return nil
}

func (app *App) initGRPCSever(ctx context.Context) error {
	app.grpcServer = grpc.NewServer()
	reflection.Register(app.grpcServer)
	descUser.RegisterUserV1Server(app.grpcServer, app.container.UserApi(ctx))

	return nil
}

func (app *App) runGRPCSever(ctx context.Context) error {
	lis, err := net.Listen("tcp", app.container.GRPCConfig().Address())
	if err != nil {
		return err
	}
	closer.Add(lis.Close)

	if err = app.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
