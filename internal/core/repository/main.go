package repository

import (
	"context"

	"gitlab.com/tour/internal/core/repository/psql"
	"gitlab.com/tour/internal/core/repository/psql/sqlc"
)

type Store interface {
	sqlc.Querier
}

// this is the main function that will be used by the service layer
func New(ctx context.Context, dsn string) Store {
	return psql.NewStore(ctx, dsn)
}
