// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package sqlc

import (
	"context"
)

type Querier interface {
	CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (string, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (string, error)
	CreateUserWithRole(ctx context.Context, arg CreateUserWithRoleParams) (string, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	DeleteUser(ctx context.Context, id string) error
	GetRefreshToken(ctx context.Context, id string) (RefreshToken, error)
	GetUser(ctx context.Context, id string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	GetUserByProperty(ctx context.Context, email string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateUserStatus(ctx context.Context, arg UpdateUserStatusParams) error
	UpdateUserStatusByEmail(ctx context.Context, arg UpdateUserStatusByEmailParams) error
}

var _ Querier = (*Queries)(nil)
