package repository

import (
	"client_service/internal/entity"
	"context"
)

type Summary interface {
	CreateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error)
	GetSummaryById(ctx context.Context, in *entity.GetRequestSummary) (*entity.Summary, error)
	GetAllSummary(ctx context.Context, in *entity.GetAllRequestSummary) ([]*entity.Summary, error)
	UpdateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error)
	DeleteSummary(ctx context.Context, in *entity.GetRequestSummary) (*entity.StatusSummary, error)
}
