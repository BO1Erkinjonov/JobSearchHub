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
	JobUseCaseServiceName    = "JobUseCaseService"
	JobUseCaseSpanRepoPrefix = "JobUseCaseRepo"
)

type JobsService interface {
	CreateJob(ctx context.Context, in *entity.Job) (*entity.Job, error)
	GetJobById(ctx context.Context, in *entity.GetReq) (*entity.Job, error)
	GetAllJobs(ctx context.Context, in *entity.GetAll) ([]*entity.Job, error)
	UpdateJob(ctx context.Context, in *entity.Job) (*entity.Job, error)
	DeleteJob(ctx context.Context, in *entity.DelReq) (*entity.StatusJob, error)
}

type newsJobService struct {
	BaseUseCase
	repo       repository.JobsService
	ctxTimeout time.Duration
}

func NewJobsServiceService(ctxTimeout time.Duration, repo repository.JobsService) newsJobService {
	return newsJobService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (n newsJobService) CreateJob(ctx context.Context, in *entity.Job) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, JobUseCaseServiceName, JobUseCaseSpanRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateJob").String(in.Id))
	defer span.End()
	return n.repo.CreateJob(ctx, in)
}

func (n newsJobService) GetJobById(ctx context.Context, in *entity.GetReq) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, JobUseCaseServiceName, JobUseCaseSpanRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetJobById").String(in.Id))
	defer span.End()
	return n.repo.GetJobById(ctx, in)
}

func (n newsJobService) GetAllJobs(ctx context.Context, in *entity.GetAll) ([]*entity.Job, error) {
	ctx, span := otlp.Start(ctx, JobUseCaseServiceName, JobUseCaseSpanRepoPrefix+"get all jobs")
	span.SetAttributes(attribute.Key("GetAllJobs").String(in.Value))
	defer span.End()
	return n.repo.GetAllJobs(ctx, in)
}

func (n newsJobService) UpdateJob(ctx context.Context, in *entity.Job) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, JobUseCaseServiceName, JobUseCaseSpanRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateJob").String(in.Id))
	defer span.End()
	return n.repo.UpdateJob(ctx, in)
}

func (n newsJobService) DeleteJob(ctx context.Context, in *entity.DelReq) (*entity.StatusJob, error) {
	ctx, span := otlp.Start(ctx, JobUseCaseServiceName, JobUseCaseSpanRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteJob").String(in.Id))
	defer span.End()
	return n.repo.DeleteJob(ctx, in)
}
