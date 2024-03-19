# Build stage
FROM golang:1.21-alpine3.19 AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .

RUN go build -o auth_service cmd/auth_service/main.go


# Final stage
FROM alpine:3.16
WORKDIR /app
RUN mkdir -p config 
COPY --from=builder /app/auth_service .
COPY --from=builder /app/doc ./doc
COPY ./internal/transport/grpc/middleware/rbac.conf /app/config/rbac.conf
COPY ./internal/transport/grpc/middleware/policy_effect.csv /app/config/policy_effect.csv

ENTRYPOINT [ "/app/auth_service" ]
