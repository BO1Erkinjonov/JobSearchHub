package usecase

import (
	"client_service/internal/entity"
	"client_service/internal/infrastructure/repository"
	"client_service/internal/pkg/otlp"
	"go.opentelemetry.io/otel/attribute"

	//"client_service/internal/pkg/otlp"
	"context"
	"time"
)

const (
	serviceNameClient           = "clientUseCase"
	serviceNameClientRepoPrefix = "clientUseCase"
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

type newsService struct {
	BaseUseCase
	repo       repository.Client
	ctxTimeout time.Duration
}

func NewClientService(ctxTimeout time.Duration, repo repository.Client) newsService {
	return newsService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u newsService) CreateClient(ctx context.Context, client *entity.Client) (*entity.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameClient, serviceNameClientRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateClient").String(client.Id))

	defer span.End()

	return u.repo.CreateClient(ctx, client)
}

func (u newsService) GetClientById(ctx context.Context, req *entity.GetRequest) (*entity.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameClient, serviceNameClientRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetClientById").String(req.ClientId))
	defer span.End()

	return u.repo.GetClientById(ctx, req)
}

func (u newsService) GetAllClients(ctx context.Context, all *entity.GetAllRequest) ([]*entity.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameClient, serviceNameClientRepoPrefix+"List")
	span.SetAttributes(attribute.Key("GetAllClients").String(all.Value))

	defer span.End()

	return u.repo.GetAllClients(ctx, all)
}

func (u newsService) UpdateClient(ctx context.Context, up *entity.Client) (*entity.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameClient, serviceNameClientRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateClient").String(up.Id))

	defer span.End()

	return u.repo.UpdateClient(ctx, up)
}

func (u newsService) DeleteClient(ctx context.Context, req *entity.DeleteReq) (*entity.Status, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameClient, serviceNameClientRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteClient").String(req.ClientId))

	defer span.End()

	return u.repo.DeleteClient(ctx, req)
}

func (u newsService) CheckUniques(ctx context.Context, field, value string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameClient, serviceNameClientRepoPrefix+"CheckUniques")
	span.SetAttributes(attribute.Key("CheckUniques").String(value))

	defer span.End()

	return u.repo.CheckUniques(ctx, field, value)
}

func (u newsService) Exists(ctx context.Context, email string) (*entity.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameClient, serviceNameClientRepoPrefix+"Exists")
	span.SetAttributes(attribute.Key("Exists").String(email))
	defer span.End()

	return u.repo.Exists(ctx, email)
}
