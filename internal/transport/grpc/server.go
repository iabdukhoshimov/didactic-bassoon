package grpc

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/redis/go-redis/v9"
	auth "gitlab.com/tour/generated/auth_service"
	user "gitlab.com/tour/generated/user_service"
	"gitlab.com/tour/internal/config"
	"gitlab.com/tour/internal/core/repository"
	"gitlab.com/tour/internal/core/service"
	"gitlab.com/tour/internal/pkg/notification"
	"gitlab.com/tour/internal/pkg/security"
	"gitlab.com/tour/internal/transport/grpc/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func New(repo repository.Store, redisClient *redis.Client, cfg *config.Config, logger *zap.Logger) (grpcServer *grpc.Server) {
	issuer := security.NewIssuer(cfg)

	// enforcer, err := casbin.NewEnforcer(cfg.Casbin.ConfigPath, cfg.Casbin.PolicyPath)
	// if err != nil {
	// 	logger.Fatal("Error while initializing casbin enforcer", zap.Error(err), zap.String("config_path", cfg.Casbin.ConfigPath), zap.String("policy_path", cfg.Casbin.PolicyPath))
	// 	return nil
	// }

	// rbacMiddleware := middleware.NewRBACMiddleware()

	grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.GrpcLoggerMiddleware,
				func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

					return middleware.JWTMiddleware(ctx, req, info, handler, issuer)
				},
				// rbacMiddleware.ValidatePermissions,
			),
		),
	)

	// grpcClient := grpc_client.NewGrpcClient(cfg)

	reflection.Register(grpcServer)

	emailNotificationService := notification.NewEmailNotificationService(cfg, redisClient, logger)

	orgKeyGenerator := security.NewOrganizationKeyGenerator(
		cfg.Organization.OrganizationSecretKey,
		cfg.Organization.OrganizationKeyLifeTime)

	auth.RegisterAuthServiceServer(grpcServer, service.NewAuthService(
		repo,
		issuer,
		redisClient,
		logger,
		emailNotificationService,
		cfg,
		orgKeyGenerator))

	user.RegisterUserServiceServer(grpcServer, service.NewUserService(repo, issuer))

	// org.RegisterOrganizationServiceServer(grpcServer, service.NewOrganizationService(repo, issuer, grpcClient, orgKeyGenerator))

	return grpcServer

}
