package mongodb

import (
	mongodb "client_service/internal/pkg/mongo"
	"client_service/internal/pkg/otlp"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"time"

	"client_service/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	clientCollectionName = "clients"
	clientServiceName    = "clientService"
	clientSpanRepoPrefix = "clientRepo"
)

type clientRepo struct {
	clientCollection *mongo.Collection
}

func NewClientRepo(client *mongodb.MongoDB, dbName string) *clientRepo {
	return &clientRepo{
		clientCollection: client.Client.Database(dbName).Collection(clientCollectionName),
	}
}

func (c *clientRepo) CreateClient(ctx context.Context, client *entity.Client) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateClient").String(client.Id))
	defer span.End()

	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()

	_, err := c.clientCollection.InsertOne(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *clientRepo) GetClientById(ctx context.Context, req *entity.GetRequest) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetClientById").String(req.ClientId))
	defer span.End()

	var client entity.Client
	err := c.clientCollection.FindOne(ctx, bson.M{"id": req.ClientId}).Decode(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *clientRepo) GetAllClients(ctx context.Context, all *entity.GetAllRequest) ([]*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("GetAllClients").String(all.Value))
	defer span.End()

	var clients []*entity.Client

	cursor, err := c.clientCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var client entity.Client
		if err := cursor.Decode(&client); err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}

	return clients, nil
}

func (c *clientRepo) UpdateClient(ctx context.Context, up *entity.Client) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateClient").String(up.Id))
	defer span.End()

	up.UpdatedAt = time.Now()

	filter := bson.M{"id": up.Id}
	update := bson.M{"$set": up}

	_, err := c.clientCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return up, nil
}

func (c *clientRepo) DeleteClient(ctx context.Context, req *entity.DeleteReq) (*entity.Status, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteClient").String(req.ClientId))
	defer span.End()

	filter := bson.M{"id": req.ClientId}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	_, err := c.clientCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &entity.Status{Status: true}, nil
}

func (c *clientRepo) CheckUniques(ctx context.Context, field, value string) (bool, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"check uniques")
	span.SetAttributes(attribute.Key("CheckUniques").String(value))
	defer span.End()

	var count int64
	filter := bson.M{field: value}
	count, err := c.clientCollection.CountDocuments(ctx, filter)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}

func (c *clientRepo) Exists(ctx context.Context, email string) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"exists")
	span.SetAttributes(attribute.Key("Exists").String(email))
	defer span.End()

	var client entity.Client
	err := c.clientCollection.FindOne(ctx, bson.M{"email": email}).Decode(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
