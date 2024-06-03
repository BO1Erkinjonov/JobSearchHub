package services

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	pb "jobs_service/genproto/jobs-service"
	"jobs_service/internal/entity"
	"jobs_service/internal/pkg/otlp"
	"jobs_service/internal/usecase"
)

type jobRPC struct {
	logger  *zap.Logger
	job     usecase.JobsService
	request usecase.RequestsService
}

const (
	serviceNameJobDelivery               = "JobDelivery"
	serviceNameJobDeliveryRepoPrefix     = "JobDelivery"
	serviceNameRequestDelivery           = "RequestDelivery"
	serviceNameRequestDeliveryRepoPrefix = "RequestDelivery"
)

func JobRPC(logget *zap.Logger, jobService usecase.JobsService, requestService usecase.RequestsService) pb.JobsServiceServer {
	return &jobRPC{
		logget,
		jobService,
		requestService,
	}
}

func (r jobRPC) CreateRequests(ctx context.Context, request *pb.Request) (*pb.Request, error) {
	ctx, span := otlp.Start(ctx, serviceNameRequestDelivery, serviceNameRequestDeliveryRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateRequests").String(request.JobId))
	defer span.End()
	resp, err := r.request.CreateRequests(ctx, &entity.Request{
		JobId:     request.JobId,
		ClientId:  request.ClientId,
		SummaryId: request.SummaryId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Request{
		JobId:           resp.JobId,
		ClientId:        resp.ClientId,
		SummaryId:       resp.SummaryId,
		StatusResp:      resp.StatusResp,
		DescriptionResp: resp.DescriptionResp,
	}, nil
}

func (r jobRPC) GetRequestByJobIdOrClientId(ctx context.Context, request *pb.GetRequest) (*pb.Request, error) {
	ctx, span := otlp.Start(ctx, serviceNameRequestDelivery, serviceNameRequestDeliveryRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetRequestByJobIdOrClientId").String(request.JobId))
	defer span.End()
	resp, err := r.request.GetRequestByJobIdOrClientId(ctx, &entity.GetRequest{request.JobId, request.ClientId})
	if err != nil {
		return nil, err
	}
	return &pb.Request{
		JobId:           resp.JobId,
		ClientId:        resp.ClientId,
		SummaryId:       resp.SummaryId,
		StatusResp:      resp.StatusResp,
		DescriptionResp: resp.DescriptionResp,
	}, nil
}

func (r jobRPC) GetAllRequest(ctx context.Context, req *pb.GetAllReq) (*pb.ListRequests, error) {
	ctx, span := otlp.Start(ctx, serviceNameRequestDelivery, serviceNameRequestDeliveryRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("GetAllRequest").String(req.Value))
	defer span.End()
	resp, err := r.request.GetAllRequest(ctx, &entity.GetAllReq{
		Page:  req.Page,
		Limit: req.Limit,
		Field: req.Field,
		Value: req.Value,
	})
	if err != nil {
		return nil, err
	}
	var requests []*pb.Request
	for _, request := range resp {
		requests = append(requests, &pb.Request{
			JobId:           request.JobId,
			ClientId:        request.ClientId,
			SummaryId:       request.SummaryId,
			StatusResp:      request.StatusResp,
			DescriptionResp: request.DescriptionResp,
		})
	}
	return &pb.ListRequests{
		Req: requests,
	}, nil
}

func (r jobRPC) UpdateRequest(ctx context.Context, request *pb.Request) (*pb.Request, error) {
	ctx, span := otlp.Start(ctx, serviceNameRequestDelivery, serviceNameRequestDeliveryRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateRequest").String(request.JobId))
	defer span.End()
	resp, err := r.request.UpdateRequest(ctx, &entity.Request{
		JobId:           request.JobId,
		ClientId:        request.ClientId,
		SummaryId:       request.SummaryId,
		StatusResp:      request.StatusResp,
		DescriptionResp: request.DescriptionResp,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Request{
		JobId:           resp.JobId,
		ClientId:        resp.ClientId,
		SummaryId:       resp.SummaryId,
		StatusResp:      resp.StatusResp,
		DescriptionResp: resp.DescriptionResp,
	}, nil
}

func (r jobRPC) DeleteRequest(ctx context.Context, request *pb.GetRequest) (*pb.StatusReq, error) {
	ctx, span := otlp.Start(ctx, serviceNameRequestDelivery, serviceNameRequestDeliveryRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteRequest").String(request.JobId))
	defer span.End()
	resp, err := r.request.DeleteRequest(ctx, &entity.GetRequest{request.JobId, request.ClientId})
	if err != nil {
		return nil, err
	}
	return &pb.StatusReq{Status: resp.Status}, nil
}

// Jobs
func (j jobRPC) CreateJob(ctx context.Context, in *pb.Job) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobDelivery, serviceNameJobDeliveryRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateJob").String(in.Id))
	defer span.End()
	resp, err := j.job.CreateJob(ctx, &entity.Job{
		Id:          in.Id,
		Owner_id:    in.OwnerId,
		Title:       in.Title,
		Description: in.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Job{
		Id:          resp.Id,
		OwnerId:     resp.Owner_id,
		Title:       resp.Title,
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt.String(),
	}, nil
}

func (j jobRPC) GetJobById(ctx context.Context, req *pb.GetReq) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobDelivery, serviceNameJobDeliveryRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetJobById").String(req.Id))
	defer span.End()
	resp, err := j.job.GetJobById(ctx, &entity.GetReq{
		Id:       req.Id,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Job{
		Id:          resp.Id,
		OwnerId:     resp.Owner_id,
		Title:       resp.Title,
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt.String(),
		UpdatedAt:   resp.UpdatedAt.String(),
		DeletedAt:   resp.DeletedAt.String(),
	}, nil
}

func (j jobRPC) GetAllJobs(ctx context.Context, all *pb.GetAll) (*pb.ListJobs, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobDelivery, serviceNameJobDeliveryRepoPrefix+"get all jobs")
	span.SetAttributes(attribute.Key("GetAllJobs").String(all.Value))
	defer span.End()
	resp, err := j.job.GetAllJobs(ctx, &entity.GetAll{
		Page:  all.Page,
		Limit: all.Limit,
		Field: all.Field,
		Value: all.Value,
	})
	if err != nil {
		return nil, err
	}
	var jobs pb.ListJobs
	for _, job := range resp {
		jobs.Jobs = append(jobs.Jobs, &pb.Job{
			Id:          job.Id,
			OwnerId:     job.Owner_id,
			Title:       job.Title,
			Description: job.Description,
			Responses:   job.Response,
			CreatedAt:   job.CreatedAt.String(),
			UpdatedAt:   job.UpdatedAt.String(),
			DeletedAt:   job.DeletedAt.String(),
		})
	}
	return &jobs, nil
}

func (j jobRPC) UpdateJob(ctx context.Context, job *pb.Job) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobDelivery, serviceNameJobDeliveryRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateJob").String(job.Id))
	defer span.End()
	resp, err := j.job.UpdateJob(ctx, &entity.Job{
		Id:          job.Id,
		Title:       job.Title,
		Description: job.Description,
		Response:    job.Responses,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Job{
		Id:          resp.Id,
		OwnerId:     resp.Owner_id,
		Title:       resp.Title,
		Description: resp.Description,
		Responses:   resp.Response,
		CreatedAt:   resp.CreatedAt.String(),
		UpdatedAt:   resp.UpdatedAt.String(),
		DeletedAt:   resp.DeletedAt.String(),
	}, nil
}

func (j jobRPC) DeleteJob(ctx context.Context, req *pb.DelReq) (*pb.StatusJob, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobDelivery, serviceNameJobDeliveryRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteJob").String(req.Id))
	defer span.End()
	status, err := j.job.DeleteJob(ctx, &entity.DelReq{
		Id:            req.Id,
		IsActive:      req.IsActive,
		IsHardDeleted: req.IsHardDeleted,
	})
	if err != nil {
		return nil, err
	}
	return &pb.StatusJob{
		Status: status.Status,
	}, nil
}
