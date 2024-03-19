package grpc_client

import (
	pb "gitlab.com/tour/generated/auth_service"
	"gitlab.com/tour/internal/config"
)

type GrpcClient struct {
	AuthService pb.AuthServiceClient
}

func NewGrpcClient(cfg *config.Config) *GrpcClient {

	return &GrpcClient{}
}
