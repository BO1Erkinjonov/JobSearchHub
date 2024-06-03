package services

import (
	pb "client_service/genproto/client-service"
	"client_service/internal/entity"
	"client_service/internal/pkg/otlp"
	"client_service/internal/usecase"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type clientRPC struct {
	logger  *zap.Logger
	client  usecase.Client
	summary usecase.Summary
}

const (
	serviceNameClientDelivery            = "ClientDelivery"
	serviceNameClientDeliveryRepoPrefix  = "ClientDelivery"
	serviceNameSummaryDelivery           = "summaryDelivery"
	serviceNameSummaryDeliveryRepoSuffix = "summaryDeliveryRepo"
)

func ClientRPC(logget *zap.Logger, clientSercie usecase.Client, summary usecase.Summary) pb.ClientServiceServer {
	return &clientRPC{
		logget,
		clientSercie,
		summary,
	}

}

func (s clientRPC) CreateSummary(ctx context.Context, summary *pb.Summary) (*pb.Summary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummaryDelivery, serviceNameSummaryDeliveryRepoSuffix+"Create")
	span.SetAttributes(attribute.Key("CreateSummary").Int(int(summary.Id)))
	defer span.End()
	resp, err := s.summary.CreateSummary(ctx, &entity.Summary{
		OwnerId:   summary.OwnerId,
		Skills:    summary.Skills,
		Bio:       summary.Bio,
		Languages: summary.Languages,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Summary{
		Id:        resp.Id,
		OwnerId:   resp.OwnerId,
		Skills:    resp.Skills,
		Bio:       resp.Bio,
		Languages: resp.Languages,
	}, nil
}

func (s clientRPC) GetSummaryById(ctx context.Context, request *pb.GetRequestSummary) (*pb.Summary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummaryDelivery, serviceNameSummaryDeliveryRepoSuffix+"Get")
	span.SetAttributes(attribute.Key("GetSummaryById").Int(int(request.Id)))
	defer span.End()
	resp, err := s.summary.GetSummaryById(ctx, &entity.GetRequestSummary{Id: request.Id})
	if err != nil {
		return nil, err
	}
	return &pb.Summary{
		Id:        resp.Id,
		OwnerId:   resp.OwnerId,
		Skills:    resp.Skills,
		Bio:       resp.Bio,
		Languages: resp.Languages,
	}, nil
}

func (s clientRPC) GetAllSummary(ctx context.Context, request *pb.GetAllRequestSummary) (*pb.GetAllResponseSummary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummaryDelivery, serviceNameSummaryDeliveryRepoSuffix+"Get all")
	span.SetAttributes(attribute.Key("GetAllSummary").String(request.Value))
	defer span.End()
	resp, err := s.summary.GetAllSummary(ctx, &entity.GetAllRequestSummary{
		Page:  request.Page,
		Limit: request.Limit,
		Field: request.Field,
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	var results pb.GetAllResponseSummary
	for _, i2 := range resp {
		results.Summary = append(results.Summary, &pb.Summary{
			Id:        i2.Id,
			OwnerId:   i2.OwnerId,
			Skills:    i2.Skills,
			Bio:       i2.Bio,
			Languages: i2.Languages,
		})
	}
	return &results, nil
}

func (s clientRPC) UpdateSummary(ctx context.Context, summary *pb.Summary) (*pb.Summary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummaryDelivery, serviceNameSummaryDeliveryRepoSuffix+"Update")
	span.SetAttributes(attribute.Key("UpdateSummary").Int(int(summary.Id)))
	defer span.End()
	resp, err := s.summary.UpdateSummary(ctx, &entity.Summary{
		Id:        summary.Id,
		OwnerId:   summary.OwnerId,
		Skills:    summary.Skills,
		Bio:       summary.Bio,
		Languages: summary.Languages,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Summary{
		Id:        resp.Id,
		OwnerId:   resp.OwnerId,
		Skills:    resp.Skills,
		Bio:       resp.Bio,
		Languages: resp.Languages,
	}, nil
}

func (s clientRPC) DeleteSummary(ctx context.Context, request *pb.GetRequestSummary) (*pb.StatusSummary, error) {
	ctx, span := otlp.Start(ctx, serviceNameSummaryDelivery, serviceNameSummaryDeliveryRepoSuffix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSummary").Int(int(request.Id)))
	defer span.End()
	resp, err := s.summary.DeleteSummary(ctx, &entity.GetRequestSummary{Id: request.Id})
	if err != nil {
		return nil, err
	}
	return &pb.StatusSummary{
		Status: resp.Status,
	}, nil
}

func (c clientRPC) CreateClient(ctx context.Context, client *pb.Client) (*pb.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientDelivery, serviceNameClientDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateClient").String(client.Id))
	defer span.End()
	resp, err := c.client.CreateClient(ctx, &entity.Client{
		Id:           client.Id,
		Role:         client.Role,
		FirstName:    client.FirstName,
		LastName:     client.LastName,
		Email:        client.Email,
		Password:     client.Password,
		RefreshToken: client.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Client{
		Id:           resp.Id,
		Role:         resp.Role,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		Email:        resp.Email,
		Password:     resp.Password,
		RefreshToken: resp.RefreshToken,
		CreatedAt:    resp.CreatedAt.String(),
	}, nil
}

func (c clientRPC) GetClientById(ctx context.Context, request *pb.GetRequest) (*pb.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientDelivery, serviceNameClientDeliveryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetClientById").String(request.ClientId))
	defer span.End()
	resp, err := c.client.GetClientById(ctx, &entity.GetRequest{
		ClientId: request.ClientId,
		IsActive: request.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Client{
		Id:           resp.Id,
		Role:         resp.Role,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		Email:        resp.Email,
		Password:     resp.Password,
		RefreshToken: resp.RefreshToken,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}, nil
}

func (c clientRPC) GetAllClients(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientDelivery, serviceNameClientDeliveryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllClients").String(request.Value))
	defer span.End()
	resp, err := c.client.GetAllClients(ctx, &entity.GetAllRequest{
		Page:  request.Page,
		Limit: request.Limit,
		Field: request.Field,
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	var clients pb.GetAllResponse
	for _, i := range resp {
		clients.AllClients = append(clients.AllClients, &pb.Client{
			Id:           i.Id,
			Role:         i.Role,
			FirstName:    i.FirstName,
			LastName:     i.LastName,
			Email:        i.Email,
			Password:     i.Password,
			RefreshToken: i.RefreshToken,
			CreatedAt:    i.CreatedAt.String(),
			UpdatedAt:    i.UpdatedAt.String(),
			DeletedAt:    i.DeletedAt.String(),
		})
		clients.Count += 1
	}
	return &clients, nil
}

func (c clientRPC) UpdateClient(ctx context.Context, client *pb.Client) (*pb.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientDelivery, serviceNameClientDeliveryRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateClient").String(client.Id))
	defer span.End()
	resp, err := c.client.UpdateClient(ctx, &entity.Client{
		Id:           client.Id,
		FirstName:    client.FirstName,
		LastName:     client.LastName,
		Email:        client.Email,
		Password:     client.Password,
		RefreshToken: client.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Client{
		Id:           resp.Id,
		Role:         resp.Role,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		Email:        resp.Email,
		Password:     resp.Password,
		RefreshToken: resp.RefreshToken,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}, nil
}

func (c clientRPC) DeleteClient(ctx context.Context, req *pb.DeleteReq) (*pb.Status, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientDelivery, serviceNameClientDeliveryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteClient").String(req.ClientId))
	defer span.End()
	status, err := c.client.DeleteClient(ctx, &entity.DeleteReq{
		ClientId:      req.ClientId,
		IsActive:      req.IsActive,
		IsHardDeleted: req.IsHardDeleted,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Status{
		Status: status.Status,
	}, nil
}

func (c clientRPC) CheckUniques(ctx context.Context, request *pb.CheckUniquesRequest) (*pb.CheckUniquesResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientDelivery, serviceNameClientDeliveryRepoPrefix+"Check")
	span.SetAttributes(attribute.Key("CheckUniques").String(request.Value))
	defer span.End()
	resp, err := c.client.CheckUniques(ctx, request.Field, request.Value)
	if err != nil {
		return nil, err
	}
	return &pb.CheckUniquesResponse{
		IsExist: resp,
	}, nil
}

func (c clientRPC) Exists(ctx context.Context, request *pb.EmailRequest) (*pb.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientDelivery, serviceNameClientDeliveryRepoPrefix+"Exists")
	span.SetAttributes(attribute.Key("Exists").String(request.Email))
	defer span.End()
	resp, err := c.client.Exists(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	return &pb.Client{
		Id:           resp.Id,
		Role:         resp.Role,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		Email:        resp.Email,
		Password:     resp.Password,
		RefreshToken: resp.RefreshToken,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}, nil
}
