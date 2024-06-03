package event

import (
	"api-gateway/api/models"
	"context"
)

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupID() string
	GetHandler() ConsumerHandler
}

type ConsumerHandler interface {
	Handle(ctx context.Context, key, value []byte) error
}

type BrokerConsumer interface {
	Run() error
	RegisterConsumer(cfg ConsumerConfig)
	Close()
}

type BrokerProducer interface {
	ProduceUserToCreate(ctx context.Context, key string, value *models.RespClient) error
	Close()
}
