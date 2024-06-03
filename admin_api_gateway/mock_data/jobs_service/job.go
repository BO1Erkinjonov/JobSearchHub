package jobs_service

import (
	"admin_api_gateway/entity"
	"context"
	"time"
)

type JobsServiceClient interface {
	CreateJob(ctx context.Context, in *entity.Job) (*entity.Job, error)
	GetJobById(ctx context.Context, in *entity.GetReq) (*entity.Job, error)
	GetAllJobs(ctx context.Context, in *entity.GetAll) ([]*entity.Job, error)
	UpdateJob(ctx context.Context, in *entity.Job) (*entity.Job, error)
	DeleteJob(ctx context.Context, in *entity.DelReq) (*entity.StatusJob, error)

	CreateRequests(ctx context.Context, in *entity.Request) (*entity.Request, error)
	GetRequestByJobIdOrClientId(ctx context.Context, in *entity.GetRequestReq) (*entity.Request, error)
	GetAllRequest(ctx context.Context, in *entity.GetAllReq) ([]*entity.Request, error)
	UpdateRequest(ctx context.Context, in *entity.Request) (*entity.Request, error)
	DeleteRequest(ctx context.Context, in *entity.GetRequestReq) (*entity.StatusReq, error)
}

type mockJobServiceClient struct {
}

func (m mockJobServiceClient) CreateRequests(ctx context.Context, in *entity.Request) (*entity.Request, error) {
	return &entity.Request{
		JobId:           "Mock Job Id",
		ClientId:        "Mock Client Id",
		SummaryId:       1,
		StatusResp:      "Mock Status Resp",
		DescriptionResp: "Mock Description Resp",
	}, nil
}

func (m mockJobServiceClient) GetRequestByJobIdOrClientId(ctx context.Context, in *entity.GetRequestReq) (*entity.Request, error) {
	return &entity.Request{
		JobId:           "Mock Job Id",
		ClientId:        "Mock Client Id",
		SummaryId:       1,
		StatusResp:      "Mock Status Resp",
		DescriptionResp: "Mock Description Resp",
	}, nil
}

func (m mockJobServiceClient) GetAllRequest(ctx context.Context, in *entity.GetAllReq) ([]*entity.Request, error) {
	lsummary := []*entity.Request{
		&entity.Request{
			JobId:           "Mock Job Id",
			ClientId:        "Mock Client Id",
			SummaryId:       1,
			StatusResp:      "Mock Status Resp",
			DescriptionResp: "Mock Description Resp",
		},
		&entity.Request{
			JobId:           "Mock Job Id",
			ClientId:        "Mock Client Id",
			SummaryId:       1,
			StatusResp:      "Mock Status Resp",
			DescriptionResp: "Mock Description Resp",
		},
		&entity.Request{
			JobId:           "Mock Job Id",
			ClientId:        "Mock Client Id",
			SummaryId:       1,
			StatusResp:      "Mock Status Resp",
			DescriptionResp: "Mock Description Resp",
		},
	}
	return lsummary, nil
}

func (m mockJobServiceClient) UpdateRequest(ctx context.Context, in *entity.Request) (*entity.Request, error) {
	return &entity.Request{
		JobId:           "Mock Job Id",
		ClientId:        "Mock Client Id",
		SummaryId:       1,
		StatusResp:      "Mock Status Resp",
		DescriptionResp: "Mock Description Resp",
	}, nil
}

func (m mockJobServiceClient) DeleteRequest(ctx context.Context, in *entity.GetRequestReq) (*entity.StatusReq, error) {
	return &entity.StatusReq{Status: true}, nil
}

func (m mockJobServiceClient) CreateJob(ctx context.Context, in *entity.Job) (*entity.Job, error) {
	return &entity.Job{
		Id:          "Mock Job Id",
		Owner_id:    "Mock Job Owner",
		Title:       "Mock title",
		Description: "Mock description",
		Response:    4,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Time{},
		DeletedAt:   time.Time{},
	}, nil
}

func (m mockJobServiceClient) GetJobById(ctx context.Context, in *entity.GetReq) (*entity.Job, error) {
	return &entity.Job{
		Id:          "Mock Job Id",
		Owner_id:    "Mock Job Owner",
		Title:       "Mock title",
		Description: "Mock description",
		Response:    4,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Time{},
		DeletedAt:   time.Time{},
	}, nil
}

func (m mockJobServiceClient) GetAllJobs(ctx context.Context, in *entity.GetAll) ([]*entity.Job, error) {
	jobs := []*entity.Job{
		&entity.Job{
			Id:          "Mock Job Id",
			Owner_id:    "Mock Job Owner",
			Title:       "Mock title",
			Description: "Mock description",
			Response:    4,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Time{},
			DeletedAt:   time.Time{},
		},
		&entity.Job{
			Id:          "Mock Job Id",
			Owner_id:    "Mock Job Owner",
			Title:       "Mock title",
			Description: "Mock description",
			Response:    2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Time{},
			DeletedAt:   time.Time{},
		},
		&entity.Job{
			Id:          "Mock Job Id",
			Owner_id:    "Mock Job Owner",
			Title:       "Mock title",
			Description: "Mock description",
			Response:    13,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Time{},
			DeletedAt:   time.Time{},
		},
	}
	return jobs, nil
}

func (m mockJobServiceClient) UpdateJob(ctx context.Context, in *entity.Job) (*entity.Job, error) {
	return &entity.Job{
		Id:          "Mock Job Id",
		Owner_id:    "Mock Job Owner",
		Title:       "Mock title",
		Description: "Mock description",
		Response:    4,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Time{},
		DeletedAt:   time.Time{},
	}, nil
}

func (m mockJobServiceClient) DeleteJob(ctx context.Context, in *entity.DelReq) (*entity.StatusJob, error) {
	return &entity.StatusJob{
		Status: true,
	}, nil
}

func NewMockClientServiceClient() JobsServiceClient {
	return &mockJobServiceClient{}
}
