package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authService "gitlab.com/tour/generated/auth_service"
	userService "gitlab.com/tour/generated/user_service"
	"gitlab.com/tour/internal/config"
	"gitlab.com/tour/internal/config/consts"
	"gitlab.com/tour/pkg/wrapper"
	mainGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func New(ctx context.Context, cfg *config.Config) *runtime.ServeMux {
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.New(map[string]string{
				consts.GrpcGatewayHttpPath:   req.URL.Path,
				consts.GrpcGatewayHttpMethod: req.Method,
			})
		}),
		runtime.WithIncomingHeaderMatcher(wrapper.CustomMatcher),
	)
	connPingService, err := mainGrpc.Dial(
		fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port),
		mainGrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil
	}

	if err := authService.RegisterAuthServiceHandler(ctx, gwMux, connPingService); err != nil {
		return nil
	}

	if err := userService.RegisterUserServiceHandler(ctx, gwMux, connPingService); err != nil {
		return nil
	}

	return gwMux
}
