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
