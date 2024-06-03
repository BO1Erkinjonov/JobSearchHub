package usecase

import (
	"client_service/internal/entity"
	"client_service/internal/infrastructure/repository"
	"client_service/internal/pkg/otlp"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

const (
	serviceNameSummary           = "summaryUseCase"
	serviceNameSummaryRepoPrefix = "summaryUseCase"
)

type Summary interface {
	CreateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error)
	GetSummaryById(ctx context.Context, in *entity.GetRequestSummary) (*entity.Summary, error)
	GetAllSummary(ctx context.Context, in *entity.GetAllRequestSummary) ([]*entity.Summary, error)
	UpdateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error)
	DeleteSummary(ctx context.Context, in *entity.GetRequestSummary) (*entity.StatusSummary, error)
}

type newsSummaryService struct {
	BaseUseCase
	repo       repository.Summary
	ctxTimeout time.Duration
}

func (n newsSummaryService) CreateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummary, serviceNameSummaryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateSummary").Int(int(in.Id)))
	defer span.End()
	return n.repo.CreateSummary(ctx, in)
}

func (n newsSummaryService) GetSummaryById(ctx context.Context, in *entity.GetRequestSummary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummary, serviceNameSummaryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetSummaryById").Int(int(in.Id)))
	defer span.End()
	return n.repo.GetSummaryById(ctx, in)
}

func (n newsSummaryService) GetAllSummary(ctx context.Context, in *entity.GetAllRequestSummary) ([]*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummary, serviceNameSummaryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllSummary").String(in.Value))
	defer span.End()
	return n.repo.GetAllSummary(ctx, in)
}

func (n newsSummaryService) UpdateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummary, serviceNameSummaryRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateSummary").Int(int(in.Id)))
	defer span.End()
	return n.repo.UpdateSummary(ctx, in)
}

func (n newsSummaryService) DeleteSummary(ctx context.Context, in *entity.GetRequestSummary) (*entity.StatusSummary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummary, serviceNameSummaryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSummary").Int(int(in.Id)))
	defer span.End()
	return n.repo.DeleteSummary(ctx, in)
}

func NewSummaryService(ctxTimeout time.Duration, repo repository.Summary) newsSummaryService {
	return newsSummaryService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}
