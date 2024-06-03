package app

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	pb "jobs_service/genproto/jobs-service"
	grpc_server "jobs_service/internal/delivery/grpc/server"
	invest_grpc "jobs_service/internal/delivery/grpc/services"
	"jobs_service/internal/infrastructure/grpc_service_clients"
	"jobs_service/internal/infrastructure/repository/postgresql"
	"jobs_service/internal/pkg/config"
	"jobs_service/internal/pkg/logger"
	mongodb "jobs_service/internal/pkg/mongo"
	"jobs_service/internal/pkg/otlp"
	"jobs_service/internal/pkg/postgres"
	"jobs_service/internal/usecase"
	"jobs_service/internal/usecase/event"
	"time"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	MongoDB        *mongodb.MongoDB
	GrpcServer     *grpc.Server
	ShutdownOTLP   func() error
	ServiceClients grpc_service_clients.ServiceClients
	BrokerProducer event.BrokerProducer
}

func NewApp(cfg *config.Config) (*App, error) {
	// init logger
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	//kafkaProducer := kafka.NewProducer(cfg, logger)

	// otlp collector initialization
	shutdownOTLP, err := otlp.InitOTLPProvider(cfg)
	if err != nil {
		return nil, err
	}

	// init db
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	mongoDB, err := mongodb.New(cfg)
	if err != nil {
		return nil, err
	}

	// grpc server init
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_server.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_recovery.UnaryServerInterceptor(),
			),
			grpc_server.UnaryInterceptorData(logger),
		)),
	)

	return &App{
		Config:       cfg,
		Logger:       logger,
		DB:           db,
		MongoDB:      mongoDB,
		GrpcServer:   grpcServer,
		ShutdownOTLP: shutdownOTLP,
	}, nil
}

func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}
	// Initialize Service Clients
	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	// repositories initialization
	job := postgresql.NewJobRepo(a.DB)
	request := postgresql.NewRequestRepo(a.DB)

	//job := mongo.NewJobRepo(a.MongoDB, "ekzamen5_db")
	//request := mongo.NewRequestRepo(a.MongoDB, "ekzamen5_db")

	//// usecase initialization
	jobusecase := usecase.NewJobsServiceService(contextTimeout, job)
	requests := usecase.NewRequestsServiceService(contextTimeout, request)
	pb.RegisterJobsServiceServer(a.GrpcServer, invest_grpc.JobRPC(a.Logger, jobusecase, requests))

	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := grpc_server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}
	return nil
}

func (a *App) Stop() {
	// close broker producer
	a.BrokerProducer.Close()
	// closing client service connections
	a.ServiceClients.Close()
	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
