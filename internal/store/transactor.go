package store

import (
	"context"
	"database/sql"
)

type txKeyT int

var txKey txKeyT = 0

type DBTransactor struct {
	db *sql.DB
}

func newTransactor(db *sql.DB) *DBTransactor {
	return &DBTransactor{db: db}
}

func (t *DBTransactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	txContext := context.WithValue(ctx, txKey, tx)

	if err := fn(txContext); err != nil {
		return err
	}

	return tx.Commit()
}
