package repository

import (
	"context"
	"jobs_service/internal/entity"
)

type JobsService interface {
	CreateJob(ctx context.Context, in *entity.Job) (*entity.Job, error)
	GetJobById(ctx context.Context, in *entity.GetReq) (*entity.Job, error)
	GetAllJobs(ctx context.Context, in *entity.GetAll) ([]*entity.Job, error)
	UpdateJob(ctx context.Context, in *entity.Job) (*entity.Job, error)
	DeleteJob(ctx context.Context, in *entity.DelReq) (*entity.StatusJob, error)
}
