package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	configpkg "jobs_service/internal/pkg/config"
)

type MongoDB struct {
	Client *mongo.Client
}

func New(config *configpkg.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(buildMongoURI(config))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to MongoDB: %s", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to ping MongoDB: %s", err.Error())
	}

	return &MongoDB{
		Client: client,
	}, nil
}

func buildMongoURI(config *configpkg.Config) string {
	uri := "mongodb://"
	uri += config.Mongodb.Host
	if len(config.Mongodb.Port) > 0 {
		uri += ":" + config.Mongodb.Port
	}
	if len(config.Mongodb.Database) > 0 {
		uri += "/" + config.Mongodb.Database
	}
	return uri
}

func (m *MongoDB) Close() {
	if m.Client != nil {
		_ = m.Client.Disconnect(context.Background())
	}
}
