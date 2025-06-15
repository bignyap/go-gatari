package dbutils

import (
	"context"
	"database/sql"
)

// https://docs.sqlc.dev/en/stable/howto/insert.html

type BulkInserter interface {
	InsertRows(ctx context.Context, tx *sql.Tx) (int64, error)
}

func InsertWithTransaction(ctx context.Context, db *sql.DB, inserter BulkInserter) (int64, error) {

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	affectedRows, err := inserter.InsertRows(ctx, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.ExecContext(ctx, "SHOW WARNINGS")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}
