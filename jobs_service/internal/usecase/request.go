package usecase

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"jobs_service/internal/entity"
	"jobs_service/internal/infrastructure/repository"
	"jobs_service/internal/pkg/otlp"
	"time"
)

const (
	RequestUseCaseServiceName    = "RequestUseCaseService"
	RequestUseCaseSpanRepoPrefix = "RequestUseCaseRepo"
)

type newsRequestService struct {
	BaseUseCase
	repo       repository.RequestsService
	ctxTimeout time.Duration
}

func NewRequestsServiceService(ctxTimeout time.Duration, repo repository.RequestsService) newsRequestService {
	return newsRequestService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

type RequestsService interface {
	CreateRequests(ctx context.Context, in *entity.Request) (*entity.Request, error)
	GetRequestByJobIdOrClientId(ctx context.Context, in *entity.GetRequest) (*entity.Request, error)
	GetAllRequest(ctx context.Context, in *entity.GetAllReq) ([]*entity.Request, error)
	UpdateRequest(ctx context.Context, in *entity.Request) (*entity.Request, error)
	DeleteRequest(ctx context.Context, in *entity.GetRequest) (*entity.StatusReq, error)
}

func (n newsRequestService) CreateRequests(ctx context.Context, in *entity.Request) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestUseCaseServiceName, RequestUseCaseSpanRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateRequests").String(in.JobId))
	defer span.End()
	return n.repo.CreateRequests(ctx, in)
}

func (n newsRequestService) GetRequestByJobIdOrClientId(ctx context.Context, in *entity.GetRequest) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestUseCaseServiceName, RequestUseCaseSpanRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetRequestByJobIdOrClientId").String(in.JobId))
	defer span.End()
	return n.repo.GetRequestByJobIdOrClientId(ctx, in)
}

func (n newsRequestService) GetAllRequest(ctx context.Context, in *entity.GetAllReq) ([]*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestUseCaseServiceName, RequestUseCaseSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("GetAllRequest").String(in.Value))
	defer span.End()
	return n.repo.GetAllRequest(ctx, in)
}

func (n newsRequestService) UpdateRequest(ctx context.Context, in *entity.Request) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestUseCaseServiceName, RequestUseCaseSpanRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateRequest").String(in.JobId))
	defer span.End()
	return n.repo.UpdateRequest(ctx, in)
}

func (n newsRequestService) DeleteRequest(ctx context.Context, in *entity.GetRequest) (*entity.StatusReq, error) {
	ctx, span := otlp.Start(ctx, RequestUseCaseServiceName, RequestUseCaseSpanRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteRequest").String(in.JobId))
	defer span.End()
	return n.repo.DeleteRequest(ctx, in)
}
