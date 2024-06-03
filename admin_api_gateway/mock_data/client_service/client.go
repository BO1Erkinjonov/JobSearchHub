package client_service

import (
	"admin_api_gateway/entity"
	"context"
	"time"
)

type NewClientMockClient interface {
	CreateClient(ctx context.Context, client *entity.Client) (*entity.Client, error)
	GetClientById(ctx context.Context, req *entity.GetRequest) (*entity.Client, error)
	GetAllClients(ctx context.Context, all *entity.GetAllRequest) ([]*entity.Client, error)
	UpdateClient(ctx context.Context, up *entity.Client) (*entity.Client, error)
	DeleteClient(ctx context.Context, req *entity.DeleteReq) (*entity.Status, error)
	CheckUniques(ctx context.Context, field, value string) (bool, error)
	Exists(ctx context.Context, email string) (*entity.Client, error)

	CreateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error)
	GetSummaryById(ctx context.Context, in *entity.GetRequestSummary) (*entity.Summary, error)
	GetAllSummary(ctx context.Context, in *entity.GetAllRequestSummary) ([]*entity.Summary, error)
	UpdateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error)
	DeleteSummary(ctx context.Context, in *entity.GetRequestSummary) (*entity.StatusSummary, error)
}

type mockClientServiceClient struct {
}

func (m mockClientServiceClient) CreateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	return &entity.Summary{
		Id:        1,
		OwnerId:   "Mock Owner Id",
		Skills:    "Mock skills",
		Bio:       "Mock bio",
		Languages: "Mock languages",
	}, nil
}

func (m mockClientServiceClient) GetSummaryById(ctx context.Context, in *entity.GetRequestSummary) (*entity.Summary, error) {
	return &entity.Summary{
		Id:        1,
		OwnerId:   "Mock Owner Id",
		Skills:    "Mock skills",
		Bio:       "Mock bio",
		Languages: "Mock languages",
	}, nil
}

func (m mockClientServiceClient) GetAllSummary(ctx context.Context, in *entity.GetAllRequestSummary) ([]*entity.Summary, error) {
	respSum := []*entity.Summary{
		&entity.Summary{
			Id:        1,
			OwnerId:   "Mock Owner Id",
			Skills:    "Mock skills",
			Bio:       "Mock bio",
			Languages: "Mock languages",
		},
		&entity.Summary{
			Id:        1,
			OwnerId:   "Mock Owner Id",
			Skills:    "Mock skills",
			Bio:       "Mock bio",
			Languages: "Mock languages",
		},
		&entity.Summary{
			Id:        1,
			OwnerId:   "Mock Owner Id",
			Skills:    "Mock skills",
			Bio:       "Mock bio",
			Languages: "Mock languages",
		},
	}
	return respSum, nil
}

func (m mockClientServiceClient) UpdateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	return &entity.Summary{
		Id:        1,
		OwnerId:   "Mock Owner Id",
		Skills:    "Mock skills",
		Bio:       "Mock bio",
		Languages: "Mock languages",
	}, nil
}

func (m mockClientServiceClient) DeleteSummary(ctx context.Context, in *entity.GetRequestSummary) (*entity.StatusSummary, error) {
	return &entity.StatusSummary{Status: true}, nil
}

func (m mockClientServiceClient) CreateClient(ctx context.Context, client *entity.Client) (*entity.Client, error) {
	return &entity.Client{
		Id:           "Mock Id",
		Role:         "Mock Role",
		FirstName:    "Mock First Name",
		LastName:     "Mock Last Name",
		Email:        "Mock Email",
		Password:     "Mock Password",
		RefreshToken: "Mock Token",
		CreatedAt:    time.Now(),
	}, nil
}

func (m mockClientServiceClient) GetClientById(ctx context.Context, req *entity.GetRequest) (*entity.Client, error) {
	return &entity.Client{
		Id:           "Mock Id",
		Role:         "Mock Role",
		FirstName:    "Mock First Name",
		LastName:     "Mock Last Name",
		Email:        "Mock Email",
		Password:     "Mock Password",
		RefreshToken: "Mock Token",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Time{},
		DeletedAt:    time.Time{},
	}, nil
}

func (m mockClientServiceClient) GetAllClients(ctx context.Context, all *entity.GetAllRequest) ([]*entity.Client, error) {
	clients := []*entity.Client{
		&entity.Client{
			Id:           "Mock Id",
			Role:         "Mock Role",
			FirstName:    "Mock First Name",
			LastName:     "Mock Last Name",
			Email:        "Mock Email",
			Password:     "Mock Password",
			RefreshToken: "Mock Token",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Time{},
			DeletedAt:    time.Time{},
		},
		&entity.Client{
			Id:           "Mock Id",
			Role:         "Mock Role",
			FirstName:    "Mock First Name",
			LastName:     "Mock Last Name",
			Email:        "Mock Email",
			Password:     "Mock Password",
			RefreshToken: "Mock Token",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Time{},
			DeletedAt:    time.Time{},
		},
		&entity.Client{
			Id:           "Mock Id",
			Role:         "Mock Role",
			FirstName:    "Mock First Name",
			LastName:     "Mock Last Name",
			Email:        "Mock Email",
			Password:     "Mock Password",
			RefreshToken: "Mock Token",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Time{},
			DeletedAt:    time.Time{},
		},
	}
	return clients, nil
}

func (m mockClientServiceClient) UpdateClient(ctx context.Context, up *entity.Client) (*entity.Client, error) {
	return &entity.Client{
		Id:           "Mock Id",
		Role:         "Mock Role",
		FirstName:    "Mock First Name",
		LastName:     "Mock Last Name",
		Email:        "Mock Email",
		Password:     "Mock Password",
		RefreshToken: "Mock Token",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    time.Time{},
	}, nil
}

func (m mockClientServiceClient) DeleteClient(ctx context.Context, req *entity.DeleteReq) (*entity.Status, error) {
	return &entity.Status{
		Status: true,
	}, nil
}

func (m mockClientServiceClient) CheckUniques(ctx context.Context, field, value string) (bool, error) {
	return true, nil
}

func (m mockClientServiceClient) Exists(ctx context.Context, email string) (*entity.Client, error) {
	return &entity.Client{
		Id:           "Mock Id",
		Role:         "Mock Role",
		FirstName:    "Mock First Name",
		LastName:     "Mock Last Name",
		Email:        "Mock Email",
		Password:     "Mock Password",
		RefreshToken: "Mock Token",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    time.Time{},
	}, nil
}

func NewMockClientServiceClient() NewClientMockClient {
	return &mockClientServiceClient{}
}
