package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/attribute"
	mongodb "jobs_service/internal/pkg/mongo"
	"jobs_service/internal/pkg/otlp"

	"jobs_service/internal/entity"
)

const (
	RequestCollectionName = "requests"
	RequestServiceName    = "requestService"
	RequestSpanRepoPrefix = "requestRepo"
)

type RequestRepo struct {
	requestCollection *mongo.Collection
}

func NewRequestRepo(client *mongodb.MongoDB, dbName string) *RequestRepo {
	return &RequestRepo{
		requestCollection: client.Client.Database(dbName).Collection(RequestCollectionName),
	}
}

func (r *RequestRepo) CreateRequest(ctx context.Context, req *entity.Request) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateRequest").String(req.JobId))
	defer span.End()

	_, err := r.requestCollection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *RequestRepo) GetRequestByJobIdOrClientId(ctx context.Context, in *entity.GetRequest) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetRequestByJobIdOrClientId").String(in.JobId))
	defer span.End()

	filter := bson.M{}
	if in.JobId != "" {
		filter["job_id"] = in.JobId
	} else if in.ClientId != "" {
		filter["client_id"] = in.ClientId
	}

	var req entity.Request
	err := r.requestCollection.FindOne(ctx, filter).Decode(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *RequestRepo) GetAllRequests(ctx context.Context, all *entity.GetAllReq) ([]*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("GetAllRequests").String(all.Value))
	defer span.End()

	filter := bson.M{}
	if all.Field == "client_id" {
		filter["client_id"] = all.Value
	} else if all.Field != "" {
		filter[all.Field] = bson.M{"$regex": fmt.Sprintf("^%s", all.Value)}
	}

	options := options.Find()
	if all.Limit != 0 {
		options.SetLimit(int64(all.Limit))
		options.SetSkip(int64(all.Limit * (all.Page - 1)))
	}

	cursor, err := r.requestCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []*entity.Request
	for cursor.Next(ctx) {
		var req entity.Request
		if err := cursor.Decode(&req); err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}

	return requests, nil
}

func (r *RequestRepo) UpdateRequest(ctx context.Context, req *entity.Request) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateRequest").String(req.JobId))
	defer span.End()

	filter := bson.M{"job_id": req.JobId, "client_id": req.ClientId}
	update := bson.M{"$set": bson.M{"status_resp": req.StatusResp, "description_resp": req.DescriptionResp}}

	_, err := r.requestCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *RequestRepo) DeleteRequest(ctx context.Context, in *entity.GetRequest) (*entity.StatusReq, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteRequest").String(in.JobId))
	defer span.End()

	filter := bson.M{}
	if in.JobId != "" {
		filter["job_id"] = in.JobId
	} else if in.ClientId != "" {
		filter["client_id"] = in.ClientId
	}

	_, err := r.requestCollection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &entity.StatusReq{Status: true}, nil
}
