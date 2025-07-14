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
	Labels []string   `debugmap:"visible"`
	Limit  int        `debugmap:"visible"`
	Offset int        `debugmap:"visible"`
}

// QueriesFn returns a slice of query filters based on the transaction filter criteria.
// It converts the filter parameters into database query filters for date range,
// pagination, and other constraints.
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

// NewTransactionService creates a new instance of TransactionService with the provided datastore.
func NewTransactionService(dt *pg.Datastore) *TransactionService {
	return &TransactionService{dt: dt}
}

// GetTransactions retrieves a list of transactions based on the provided filter criteria.
func (t *TransactionService) GetTransactions(ctx context.Context, filter *TransactionFilter) ([]entity.Transaction, error) {
	return t.dt.QueryTransactions(ctx, filter.QueriesFn()...)
}

// GetTransaction retrieves a single transaction by its hash identifier.
// Returns nil if no transaction is found with the given hash.
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

// CreateOrUpdate creates a new transaction or updates an existing one based on the hash.
// If a transaction with the same hash exists, it updates the existing record.
// Otherwise, it creates a new transaction record.
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

// Delete removes a transaction from the database by its ID.
func (t *TransactionService) Delete(ctx context.Context, id int64) error {
	return t.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		return w.DeleteTransaction(ctx, id)
	})
}
