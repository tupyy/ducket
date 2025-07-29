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
		return nil, NewErrTransactionNotFoundByHash(hash)
	}
	return &tt[0], err
}

func (t *TransactionService) GetTransactionById(ctx context.Context, id int) (*entity.Transaction, error) {
	tt, err := t.dt.QueryTransactions(ctx, pg.TransactionIDQueryFilter(id))
	if err != nil {
		return nil, err
	}
	if len(tt) == 0 {
		return nil, NewErrTransactionNotFound(id)
	}
	return &tt[0], err
}

// CreateOrUpdate creates a new transaction or updates an existing one based on the hash.
// If a transaction with the same hash exists, it updates the existing record.
// Otherwise, it creates a new transaction record.
func (t *TransactionService) Create(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error) {
	tt, err := t.dt.QueryTransactions(ctx, pg.TransactionHashQueryFilter(transaction.Hash))
	if err != nil {
		return transaction, err
	}

	if len(tt) == 1 {
		return tt[0], NewErrTransactionExistsAlready(int(tt[0].ID))
	}

	if err := t.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		id, err := w.WriteTransaction(ctx, transaction)
		if err != nil {
			return err
		}
		transaction.ID = id

		labelSrv := NewLabelService(t.dt)
		relationships := []entity.Relationship{}
		for _, a := range transaction.Labels {
			labelID := 0
			l, err := labelSrv.Get(ctx, a.Label.Key, a.Label.Value)
			if err != nil {
				return err
			}

			if l != nil {
				labelID = l.ID
			}

			if l == nil {
				label, err := labelSrv.Create(ctx, a.Label.Key, a.Label.Value)
				if err != nil {
					return err
				}
				labelID = label.ID
			}

			if a.RuleID != nil {
				relationships = append(relationships, entity.NewLabeRuleTransactionRelationship(labelID, *a.RuleID, int(transaction.ID)))
				continue
			}
			relationships = append(relationships, entity.NewLabelTransaction(labelID, int(transaction.ID)))
		}

		if len(relationships) == 0 {
			return nil
		}

		return w.WriteRelationships(ctx, relationships)
	}); err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (t *TransactionService) Update(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error) {
	tt, err := t.dt.QueryTransactions(ctx, pg.TransactionIDQueryFilter(int(transaction.ID)))
	if err != nil {
		return transaction, err
	}

	if len(tt) == 0 {
		return entity.Transaction{}, NewErrTransactionNotFound(int(transaction.ID))
	}

	oldTransaction := tt[0]

	if err := t.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		// remove old relationships
		existingRelationships := []entity.Relationship{}
		for _, label := range oldTransaction.Labels {
			if label.RuleID != nil {
				existingRelationships = append(existingRelationships, entity.NewLabeRuleTransactionRelationship(label.Label.ID, *label.RuleID, int(oldTransaction.ID)))
				continue
			}
			existingRelationships = append(existingRelationships, entity.NewLabelTransaction(label.Label.ID, int(oldTransaction.ID)))
		}

		if len(existingRelationships) > 0 {
			if err := w.DeleteRelationships(ctx, existingRelationships); err != nil {
				return err
			}
		}

		_, err := w.WriteTransaction(ctx, transaction)
		if err != nil {
			return err
		}

		labelSrv := NewLabelService(t.dt)
		relationships := []entity.Relationship{}
		for _, a := range transaction.Labels {
			labelID := 0
			l, err := labelSrv.Get(ctx, a.Label.Key, a.Label.Value)
			if err != nil {
				return err
			}

			if l != nil {
				labelID = l.ID
			}

			if l == nil {
				label, err := labelSrv.Create(ctx, a.Label.Key, a.Label.Value)
				if err != nil {
					return err
				}
				labelID = label.ID
			}

			if a.RuleID != nil {
				relationships = append(relationships, entity.NewLabeRuleTransactionRelationship(labelID, *a.RuleID, int(transaction.ID)))
				continue
			}
			relationships = append(relationships, entity.NewLabelTransaction(labelID, int(transaction.ID)))
		}

		if len(relationships) == 0 {
			return nil
		}

		return w.WriteRelationships(ctx, relationships)
	}); err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (t *TransactionService) Labels(ctx context.Context, transactionID int) ([]entity.Label, error) {
	tt, err := t.dt.QueryTransactions(ctx, pg.TransactionIDQueryFilter(transactionID))
	if err != nil {
		return []entity.Label{}, err
	}

	if len(tt) == 0 {
		return []entity.Label{}, NewErrTransactionNotFound(transactionID)
	}

	labels := make([]entity.Label, 0, len(tt[0].Labels))
	for _, a := range tt[0].Labels {
		labels = append(labels, a.Label)
	}

	return labels, nil
}

// Delete removes a transaction from the database by its ID.
func (t *TransactionService) Delete(ctx context.Context, id int64) error {
	return t.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		return w.DeleteTransaction(ctx, id)
	})
}

// UpdateInfo updates only the info field of an existing transaction.
func (t *TransactionService) UpdateInfo(ctx context.Context, id int64, info string) (entity.Transaction, error) {
	// First check if transaction exists
	transaction, err := t.GetTransactionById(ctx, int(id))
	if err != nil {
		return entity.Transaction{}, err
	}

	// Update the info field
	transaction.Info = &info

	// Write the updated transaction
	if err := t.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		_, err := w.WriteTransaction(ctx, *transaction)
		return err
	}); err != nil {
		return entity.Transaction{}, err
	}

	return *transaction, nil
}
