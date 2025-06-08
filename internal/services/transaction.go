package services

import (
	"context"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

const (
	defaultPageSize = 100
)

//go:generate go run github.com/ecordell/optgen -output zz_generated.transaction_filter.go . TransactionFilter
type TransactionFilter struct {
	Start  *time.Time `debugmap:"visible"`
	End    *time.Time `debugmap:"visible"`
	Tags   []string   `debugmap:"visible"`
	Limit  int        `debugmap:"visible"`
	Offset int        `debugmap:"visible"`
}

func (tf *TransactionFilter) QueriesFn() []pg.QueryFilter {
	qf := []pg.QueryFilter{}

	qf = append(qf,
		pg.IntervalDateQueryFilter(tf.Start, tf.End),
		pg.LimitQueryFilter(tf.Limit),
		pg.OffsetQueryFilter(tf.Offset),
	)

	return qf
}

type TransactionService struct {
	dt *pg.Datastore
}

func NewTransactionService(dt *pg.Datastore) *TransactionService {
	return &TransactionService{dt: dt}
}

func (t *TransactionService) GetTransactions(ctx context.Context, filter *TransactionFilter) ([]entity.Transaction, error) {
	return t.dt.QueryTransactions(ctx, filter.QueriesFn()...)
}

func (t *TransactionService) GetTransaction(ctx context.Context, hash string) (*entity.Transaction, error) {
	tt, err := t.dt.QueryTransactions(ctx, pg.TransactionHashQueryFilter(hash))
	if err != nil {
		return nil, err
	}
	if len(tt) == 0 {
		return nil, nil
	}
	return &tt[0], err
}

func (t *TransactionService) CreateOrUpdate(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error) {
	tt, err := t.dt.QueryTransactions(ctx, pg.TransactionHashQueryFilter(transaction.Hash))
	if err != nil {
		return transaction, err
	}

	if len(tt) == 1 {
		transaction.ID = tt[0].ID
	}

	if err := t.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		id, err := w.WriteTransaction(ctx, transaction)
		if err != nil {
			return err
		}
		transaction.ID = id
		return nil
	}); err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (t *TransactionService) Delete(ctx context.Context, id int64) error {
	return t.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		return w.DeleteTransaction(ctx, id)
	})
}
