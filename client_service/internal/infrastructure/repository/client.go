package repository

import (
	"client_service/internal/entity"
	"context"
)

type Client interface {
	CreateClient(ctx context.Context, client *entity.Client) (*entity.Client, error)
	GetClientById(ctx context.Context, req *entity.GetRequest) (*entity.Client, error)
	GetAllClients(ctx context.Context, all *entity.GetAllRequest) ([]*entity.Client, error)
	UpdateClient(ctx context.Context, up *entity.Client) (*entity.Client, error)
	DeleteClient(ctx context.Context, req *entity.DeleteReq) (*entity.Status, error)
	CheckUniques(ctx context.Context, field, value string) (bool, error)
	Exists(ctx context.Context, email string) (*entity.Client, error)
}
