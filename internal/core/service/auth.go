package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/redis/go-redis/v9"
	"gitlab.com/tour/internal/config"
	"gitlab.com/tour/internal/pkg/notification"
	"go.uber.org/zap"

	"gitlab.com/tour/internal/pkg/serializer"
	"gopkg.in/guregu/null.v4/zero"

	pbAuth "gitlab.com/tour/generated/auth_service"
	"gitlab.com/tour/internal/core/repository"
	"gitlab.com/tour/internal/core/repository/psql/sqlc"
	"gitlab.com/tour/internal/pkg/logger"
	"gitlab.com/tour/internal/pkg/security"

	models "gitlab.com/tour/internal/core/lib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	pbAuth.UnimplementedAuthServiceServer
	store                    repository.Store
	emailNotificationService *notification.EmailNotificationService
	issuer                   security.IssuerInterface
	redis                    *redis.Client
	logger                   *zap.Logger
	config                   *config.Config
	orgKeyGenerator          *security.OrganizationKeyGenerator
}

func NewAuthService(
	store repository.Store,
	issuer security.IssuerInterface,
	redisClient *redis.Client,
	_logger *zap.Logger,
	emailNotificationService *notification.EmailNotificationService,
	config *config.Config,
	orgKeyGenerator *security.OrganizationKeyGenerator) *AuthService {
	return &AuthService{
		store:                    store,
		emailNotificationService: emailNotificationService,
		issuer:                   issuer,
		redis:                    redisClient,
		logger:                   _logger,
		config:                   config,
		orgKeyGenerator:          orgKeyGenerator,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pbAuth.RegisterRequest) (*emptypb.Empty, error) {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	isValidEmail := emailRegex.MatchString(req.Email)
	if !isValidEmail {
		return nil, status.Error(codes.InvalidArgument, "Email is not valid")
	}

	_, err := s.store.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, status.Error(codes.InvalidArgument, "User already exists")
	}

	if len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 6 characters long")
	}

	hashPassword, err := security.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, status.Error(codes.Canceled, "could not create tokens")
	}

	dbReq := sqlc.CreateUserWithRoleParams{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   hashPassword,
		UserStatus: zero.StringFrom("PENDING"),
	}

	_, err = s.store.CreateUserWithRole(ctx, dbReq)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, "Internal error: could not create user")
	}

	go s.emailNotificationService.SendVerificationCode(req.Email) // todo: write to logs if occurs error

	return &emptypb.Empty{}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pbAuth.LoginRequest) (*pbAuth.LoginResponse, error) {
	resp, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "There is no such admin with this username or invalid password")
	}

	if match, err := security.ComparePassword(resp.Password, req.Password); !match {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	role, err := serializer.SerializeEnumToString(resp.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	claims := models.UserClaims{
		Id:   resp.ID,
		Role: role,
	}

	accessToken, err := s.issuer.NewAccessToken(claims)
	if err != nil {
		logger.Log.Error("Internal error: could not create tokens")
		return nil, status.Error(codes.PermissionDenied, "Invalid password")
	}

	refreshToken, err := s.issuer.NewRefreshToken(claims)
	if err != nil {
		return nil, status.Error(codes.Canceled, "Internal error: could not create tokens")
	}

	return &pbAuth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Id:           resp.ID,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		Status:       resp.UserStatus.String,
	}, nil
}

func (s *AuthService) EmailVerification(ctx context.Context, req *pbAuth.EmailVerificationRequest) (*emptypb.Empty, error) {
	var errorCodeTriedCount int
	errorCodeTriedCountValue, err := s.redis.Get(ctx, fmt.Sprintf(req.Email+"_count")).Result()
	if err == nil {
		errorCodeTriedCount, err = strconv.Atoi(errorCodeTriedCountValue)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
		if errorCodeTriedCount == s.config.EmailNotificationConfig.RandomCodeErrorTriesCount {
			return nil, status.Error(codes.InvalidArgument, "Tried to many attempts, try again later")
		}
	}

	code, err := s.redis.Get(ctx, req.Email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, status.Error(codes.InvalidArgument, "Code expired")
		} else {
			s.logger.Error(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if code != req.Code {
		codeLiveDuration, err := time.ParseDuration(s.config.EmailNotificationConfig.RandomCodeErrorBlockTime)
		if err != nil {
			s.logger.Error("Failed to parse duration", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}
		s.redis.Set(ctx, req.Email+"_count", errorCodeTriedCount+1, codeLiveDuration)

		return nil, status.Error(codes.InvalidArgument, "Invalid code")
	}

	// organizationID, err := s.store.GetOrganizationByOwnersEmail(ctx, req.Email)
	// if err != nil {
	// 	s.logger.Error(err.Error())
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	// organizationClaims := security.OrganizationClaims{
	// 	Id: organizationID,
	// }

	// organizationSecretKey, err := s.orgKeyGenerator.GenerateOrganizationKey(organizationClaims)
	// if err != nil {
	// 	s.logger.Error(err.Error())
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	// err = s.store.UpdateOrganizationSecretKey(ctx, sqlc.UpdateOrganizationSecretKeyParams{
	// 	ID:        organizationID,
	// 	SecretKey: zero.StringFrom(organizationSecretKey),
	// })
	// if err != nil {
	// 	s.logger.Error(err.Error())
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	err = s.store.UpdateUserStatusByEmail(ctx, sqlc.UpdateUserStatusByEmailParams{
		Email:      req.Email,
		UserStatus: zero.StringFrom("ACTIVE"),
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *pbAuth.LogoutRequest) (*emptypb.Empty, error) {
	err := s.store.DeleteRefreshToken(ctx, req.Token)
	if err != nil {
		logger.Log.Error("There is an error while deleting refresh token at service level ==> " + err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pbAuth.RefreshTokenRequest) (*pbAuth.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		logger.Log.Error("Refresh token not provided")
		return nil, fmt.Errorf("refresh token not provided")
	}

	claims, err := s.issuer.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	accessToken, err := s.issuer.NewAccessToken(*claims)
	if err != nil {
		logger.Log.Error("Internal error: could not create tokens")
		return nil, status.Error(codes.Canceled, "Internal error: could not create tokens")
	}

	newRefreshToken, err := s.issuer.NewRefreshToken(*claims)
	if err != nil {
		logger.Log.Error("Internal error: could not create tokens")
		return nil, status.Error(codes.Canceled, "Internal error: could not create tokens")
	}

	return &pbAuth.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) SendCode(ctx context.Context, req *pbAuth.SendCodeRequest) (*empty.Empty, error) {
	errorCodeTriedCountValue, err := s.redis.Get(ctx, fmt.Sprintf(req.Email+"_count")).Result()
	if err == nil {
		errorCodeTriedCount, err := strconv.Atoi(errorCodeTriedCountValue)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
		if errorCodeTriedCount == s.config.EmailNotificationConfig.RandomCodeErrorTriesCount {
			return nil, status.Error(codes.InvalidArgument, "Tried too many attempts, try again later")
		}
	}

	_, err = s.redis.Get(ctx, req.Email).Result()
	if err == nil {
		return nil, status.Error(codes.InvalidArgument, "Code already sent")
	}

	go s.emailNotificationService.SendVerificationCode(req.Email)
	return &emptypb.Empty{}, nil
}

func (s *AuthService) ParseOrganizationToken(ctx context.Context, req *pbAuth.ParseOrganizationTokenRequest) (*pbAuth.ParseOrganizationTokenResponse, error) {
	claims, err := s.orgKeyGenerator.ParseOrganizationKey(req.Token)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, "Invalid token")
	}

	return &pbAuth.ParseOrganizationTokenResponse{
		OrgId: claims.Id,
	}, nil
}

func (s *AuthService) ParseUserToken(ctx context.Context, req *pbAuth.ParseUserTokenRequest) (*pbAuth.ParseUserTokenResponse, error) {
	claims, err := s.issuer.ParseAccessToken(req.Token)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, "Invalid token")
	}

	return &pbAuth.ParseUserTokenResponse{
		Id:   claims.Id,
		Role: claims.Role,
	}, nil
}
