package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// "github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	"api-gateway/api"

	grpcService "api-gateway/internal/infastructure/grpc_service_client"
	"api-gateway/internal/infastructure/kafka"
	"api-gateway/internal/usecase/event"

	// "api-gateway/internal/infrastructure/kafka"
	// "api-gateway/internal/infrastructure/repository/postgresql"
	// redisrepo "api-gateway/internal/infrastructure/repository/redis"
	"api-gateway/internal/pkg/config"
	"api-gateway/internal/pkg/logger"

	"api-gateway/internal/pkg/otlp"
	// "api-gateway/internal/pkg/policy"
	"api-gateway/internal/pkg/postgres"
	// "api-gateway/internal/pkg/redis"
	// "api-gateway/internal/usecase/app_version"
	// "api-gateway/internal/usecase/event"
	// "evrone_api_gateway/internal/usecase/refresh_token"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	server         *http.Server
	ShutdownOTLP   func() error
	Clients        grpcService.ServiceClient
	BrokerProducer event.BrokerProducer
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// kafka producer init
	kafkaProducer := kafka.NewProducer(&cfg, logger)

	// postgres init
	db, err := postgres.New(&cfg)
	if err != nil {
		return nil, err
	}

	// otlp collector init
	shutdownOTLP, err := otlp.InitOTLPProvider(&cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:         &cfg,
		Logger:         logger,
		DB:             db,
		ShutdownOTLP:   shutdownOTLP,
		BrokerProducer: kafkaProducer,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	clients, err := grpcService.New(a.Config)
	if err != nil {
		return err
	}
	a.Clients = clients

	// initialize cache
	// cache := redisrepo.NewCache(a.RedisDB)

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Service:        clients,
		BrokerProducer: a.BrokerProducer,
	})
	// if err = a.Enforcer.LoadPolicy(); err != nil {
	// 	return fmt.Errorf("error during enforcer load policy: %w", err)
	// }

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close database
	a.DB.Close()

	// close grpc connections
	a.Clients.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
