package pg

import (
	"context"
	"errors"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"github.com/jackc/pgx/v5"
)

type transactionWriterTx struct {
	tx pgx.Tx
}

func (t *transactionWriterTx) Write(ctx context.Context, transaction entity.Transaction) error {
	return errors.New("not implementated")
}

func (r *transactionWriterTx) Delete(ctx context.Context, id string) error {
	return errors.New("not implementated")
}
