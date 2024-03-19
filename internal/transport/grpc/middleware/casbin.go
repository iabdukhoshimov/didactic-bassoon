package middleware

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/spf13/cast"
	"gitlab.com/tour/internal/config/consts"
	"gitlab.com/tour/internal/pkg/logger"
	"gitlab.com/tour/internal/pkg/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	notProtectedMethods = []string{
		"/auth_service.AuthService/Register",
		"/auth_service.AuthService/Login",
		"/auth_service.AuthService/RefreshToken",
		"/auth_service.AuthService/EmailVerification",
		"/auth_service.AuthService/ParseOrganizationToken",
		"/auth_service.AuthService/ParseUserToken",
		"/auth_service.AuthService/SendCode",
		"/organization_service.OrganizationService/PhoneCallCreate",
		"/organization_service.OrganizationService/NewsGetAll",
		"/organization_service.OrganizationService/NewsGet",
		"/organization_service.OrganizationService/GetAllFAQ",
		"/organization_service.OrganizationService/GetFAQ",
	}
)

type rbacMiddleware struct {
	enforcer       *casbin.Enforcer
	publicEnpoints map[string]string
}

func NewRBACMiddleware(enforcer *casbin.Enforcer) *rbacMiddleware {
	return &rbacMiddleware{enforcer: enforcer, publicEnpoints: map[string]string{
		"/v1/auth/register":            "POST",
		"/v1/auth/login":               "POST",
		"/v1/auth/email-verification":  "POST",
		"/v1/auth/logout":              "POST",
		"/v1/auth/refresh":             "POST",
		"/v1/auth/send-code":           "POST",
		"/v1/organizations/phone-call": "POST",
		"/v1/organizations/news":       "GET",
		"/v1/organizations/news/:id":   "GET",
		"/v1/organizations/faqs":       "GET",
		"/v1/organizations/faqs/:id":   "GET",
	}}
}

func JWTMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, issuer *security.Issuer) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata failed")
	}

	if slices.Contains(notProtectedMethods, info.FullMethod) {
		return handler(ctx, req)
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization header not provided")
	}

	tokenParts := strings.Split(authHeader[0], "Bearer ")
	if len(tokenParts) != 2 {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization header format")
	}

	tokenString := tokenParts[1]

	claims, err := issuer.ParseAccessToken(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	newCtx := context.WithValue(ctx, consts.RoleKey, claims.Role)
	newCtx = context.WithValue(newCtx, consts.UserIdKey, claims.Id)

	return handler(newCtx, req)
}

func (r *rbacMiddleware) ValidatePermissions(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing auth token")
	}

	if slices.Contains(notProtectedMethods, info.FullMethod) {
		return handler(ctx, req)
	}

	httpPath := getStringOrEmptyItem(md, consts.GrpcGatewayHttpPath)
	httpMethod := getStringOrEmptyItem(md, consts.GrpcGatewayHttpMethod)

	if httpPath == "" {
		return handler(ctx, req)
	}

	if method, ok := r.publicEnpoints[httpPath]; ok && method == httpMethod {
		return handler(ctx, req)
	}

	role := cast.ToString(ctx.Value(consts.RoleKey))

	accessGranted, err := r.enforcer.Enforce(role, httpPath, httpMethod)
	if err != nil {
		return resp, fmt.Errorf("failed to enforce: %w", err)
	}

	if !accessGranted {
		return resp, status.Errorf(codes.PermissionDenied, "access denied")
	}

	resp, err = handler(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
		return resp, err
	}

	return resp, nil
}

func getStringOrEmptyItem(md metadata.MD, key string) string {
	if len(md[key]) > 0 {
		return md[key][0]
	}

	return ""
}
