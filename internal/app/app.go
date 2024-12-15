package app

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sSmok/auth/internal/interceptor"
	"github.com/sSmok/auth/internal/metric"
	descAccess "github.com/sSmok/auth/pkg/access_v1"
	descAuth "github.com/sSmok/auth/pkg/auth_v1"
	descUser "github.com/sSmok/auth/pkg/user_v1"
	"github.com/sSmok/platform_common/pkg/closer"
	"github.com/sSmok/platform_common/pkg/config"
	platformInterceptor "github.com/sSmok/platform_common/pkg/interceptor"
	"github.com/sSmok/platform_common/pkg/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "auth-service"

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.Parse()
}

// App представляет основную логику приложения, содержит DI контейнер и сервер gRPC
type App struct {
	container  *container
	grpcServer *grpc.Server
}

// NewApp создает новое приложение
func NewApp(ctx context.Context) (*App, error) {
	app := &App{}
	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// Run запускает grpc сервер и контролирует закрытие ресурсов
func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	go func() {
		err := runPrometheus()
		if err != nil {
			log.Fatalf("failed to run prometheus: %v", err)
		}
	}()

	return app.runGRPCServer()
}

func (app *App) initDeps(ctx context.Context) error {
	deps := []func(context.Context) error{
		app.initConfig,
		app.initContainer,
		app.initTracing,
		app.initMetric,
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

func (app *App) initContainer(_ context.Context) error {
	app.container = newContainer()
	return nil
}

func (app *App) initTracing(_ context.Context) error {
	return tracing.Init(serviceName)
}

func (app *App) initGRPCSever(ctx context.Context) error {
	chain := grpcMiddleware.ChainUnaryServer(
		interceptor.MetricsInterceptor,
		platformInterceptor.Log,
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
	)
	app.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(chain))
	reflection.Register(app.grpcServer)
	descUser.RegisterUserV1Server(app.grpcServer, app.container.UserAPI(ctx))
	descAccess.RegisterAccessV1Server(app.grpcServer, app.container.AccessAPI(ctx))
	descAuth.RegisterAuthV1Server(app.grpcServer, app.container.AuthAPI(ctx))

	return nil
}

func (app *App) initMetric(ctx context.Context) error {
	return metric.Init(ctx)
}

func (app *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", app.container.GRPCConfig().Address())
	if err != nil {
		return err
	}
	closer.Add(lis.Close)

	log.Printf("GRPC server is running on %s", app.container.GRPCConfig().Address())

	if err = app.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:              "localhost:2112",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("Prometheus server is running on %s", "localhost:2112")

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
