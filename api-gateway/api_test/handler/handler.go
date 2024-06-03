package handler

import (
	"api-gateway/internal/infastructure/grpc_service_client"
	"api-gateway/internal/pkg/config"
	token "api-gateway/internal/pkg/tokens"
	"api-gateway/internal/usecase/event"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"time"
)

type handlerV1 struct {
	ContextTimeout time.Duration
	jwthandler     token.JWTHandler
	log            *zap.Logger
	serviceManager grpc_service_clients.ServiceClient
	cfg            *config.Config
	BrokerProducer event.BrokerProducer
	//kafka          *kafka.Produce
}

// HandlerV1Config ...
type HandlerV1Config struct {
	ContextTimeout time.Duration
	Jwthandler     token.JWTHandler
	Logger         *zap.Logger
	Service        grpc_service_clients.ServiceClient
	Config         *config.Config
	Enforcer       casbin.Enforcer
	BrokerProducer event.BrokerProducer
	//Kafka          *kafka.Produce
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		jwthandler:     c.Jwthandler,
		log:            c.Logger,
		serviceManager: c.Service,
		cfg:            c.Config,
		BrokerProducer: c.BrokerProducer,
	}
}
