package sqlc

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

type QueriesImpl struct {
	*Queries
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit(ctx)
}

func (store *Store) BeginTx(ctx context.Context, opts *pgx.TxOptions) (*QueriesImpl, error) {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &QueriesImpl{New(tx)}, nil
}
