package mongodb

import (
	"client_service/internal/entity"
	mongodb "client_service/internal/pkg/mongo"
	"client_service/internal/pkg/otlp"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
)

const (
	summaryCollectionName = "summary"
	summaryServiceName    = "summaryService"
	summarySpanRepoPrefix = "summaryRepo"
)

type SummaryRepo struct {
	collection *mongo.Collection
}

func NewSummaryRepo(client *mongodb.MongoDB, dbName string) *SummaryRepo {
	return &SummaryRepo{
		collection: client.Client.Database(dbName).Collection(summaryCollectionName),
	}
}

func (s *SummaryRepo) CreateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateSummary").Int(int(in.Id)))
	defer span.End()

	_, err := s.collection.InsertOne(ctx, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *SummaryRepo) GetSummaryById(ctx context.Context, in *entity.GetRequestSummary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetSummaryById").Int(int(in.Id)))
	defer span.End()

	var summary entity.Summary
	err := s.collection.FindOne(ctx, bson.M{"id": in.Id}).Decode(&summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (s *SummaryRepo) GetAllSummary(ctx context.Context, in *entity.GetAllRequestSummary) ([]*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllSummary").String(in.Value))
	defer span.End()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []*entity.Summary
	for cursor.Next(ctx) {
		var summary entity.Summary
		if err := cursor.Decode(&summary); err != nil {
			return nil, err
		}
		summaries = append(summaries, &summary)
	}

	return summaries, nil
}

func (s *SummaryRepo) UpdateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateSummary").Int(int(in.Id)))
	defer span.End()

	filter := bson.M{"id": in.Id}
	update := bson.M{"$set": in}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (s *SummaryRepo) DeleteSummary(ctx context.Context, in *entity.GetRequestSummary) (*entity.StatusSummary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSummary").Int(int(in.Id)))
	defer span.End()

	filter := bson.M{"id": in.Id}
	_, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return &entity.StatusSummary{Status: false}, err
	}

	return &entity.StatusSummary{Status: true}, nil
}
