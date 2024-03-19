package consts

type ContextKey string

const (
	GrpcGatewayHttpPath   = "grpc-gateway-http-path"
	GrpcGatewayHttpMethod = "grpc-gateway-http-method"
)

type AuthInfo string

const RoleKey ContextKey = "role"
const UserIdKey ContextKey = "userId"
