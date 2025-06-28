package dbutils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BulkInserter interface {
	InsertRows(ctx context.Context, tx pgx.Tx) (int64, error)
}

func InsertWithTransaction(ctx context.Context, pool *pgxpool.Pool, inserter BulkInserter) (int64, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) // safe to call always

	affectedRows, err := inserter.InsertRows(ctx, tx)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}
