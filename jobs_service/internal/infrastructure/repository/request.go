package repository

import (
	"context"
	"jobs_service/internal/entity"
)

type RequestsService interface {
	CreateRequests(ctx context.Context, in *entity.Request) (*entity.Request, error)
	GetRequestByJobIdOrClientId(ctx context.Context, in *entity.GetRequest) (*entity.Request, error)
	GetAllRequest(ctx context.Context, in *entity.GetAllReq) ([]*entity.Request, error)
	UpdateRequest(ctx context.Context, in *entity.Request) (*entity.Request, error)
	DeleteRequest(ctx context.Context, in *entity.GetRequest) (*entity.StatusReq, error)
}
