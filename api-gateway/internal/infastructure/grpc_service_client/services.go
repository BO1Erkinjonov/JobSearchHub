package grpc_service_clients

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"api-gateway/genproto/client-service"
	"api-gateway/genproto/jobs-service"
	"google.golang.org/grpc"

	"api-gateway/internal/pkg/config"
)

type ServiceClient interface {
	ClientService() client.ClientServiceClient
	JobService() job.JobsServiceClient
	Close()
}

type serviceClient struct {
	connections   []*grpc.ClientConn
	clientService client.ClientServiceClient
	jobService    job.JobsServiceClient
}

func New(cfg *config.Config) (ServiceClient, error) {
	connClientService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.ClientService.Host, cfg.ClientService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)

	connJobService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.JobsService.Host, cfg.JobsService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return &serviceClient{
		clientService: client.NewClientServiceClient(connClientService),
		connections: []*grpc.ClientConn{
			connClientService,
			connJobService,
		},
		jobService: job.NewJobsServiceClient(connJobService),
	}, nil
}

func (s *serviceClient) JobService() job.JobsServiceClient {
	return s.jobService
}

func (s *serviceClient) ClientService() client.ClientServiceClient {
	return s.clientService
}

func (s *serviceClient) Close() {
	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			// should be replaced by logger soon
			fmt.Printf("error while closing grpc connection: %v", err)
		}
	}
}
