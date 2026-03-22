package services

import (
	"context"
	"fmt"

	"github.com/tupyy/ducket/internal/entity"
	"github.com/tupyy/ducket/internal/store"
	srvErrors "github.com/tupyy/ducket/pkg/errors"
)

type TransactionService struct {
	st *store.Store
}

func NewTransactionService(st *store.Store) *TransactionService {
	return &TransactionService{st: st}
}

func (t *TransactionService) List(ctx context.Context, filter string, tags []string, sort []store.SortParam, limit, offset int) ([]entity.Transaction, int, error) {
	return t.st.ListTransactions(ctx, filter, tags, sort, limit, offset)
}

func (t *TransactionService) Get(ctx context.Context, id int64) (*entity.Transaction, error) {
	return t.st.GetTransaction(ctx, id)
}

func (t *TransactionService) Create(ctx context.Context, txn entity.Transaction) (entity.Transaction, error) {
	if !txn.Kind.Valid() {
		return txn, srvErrors.NewValidationError(fmt.Sprintf("invalid transaction kind: %q", txn.Kind))
	}

	var id int64
	if err := t.st.WithTx(ctx, func(ctx context.Context) error {
		existing, err := t.st.GetTransactionByHash(ctx, txn.Hash)
		if err != nil && !srvErrors.IsResourceNotFoundError(err) {
			return err
		}
		if existing != nil {
			return srvErrors.NewDuplicateResourceError("transaction", "hash", txn.Hash)
		}

		id, err = t.st.CreateTransaction(ctx, txn)
		return err
	}); err != nil {
		return txn, err
	}

	txn.ID = id
	return txn, nil
}

func (t *TransactionService) Update(ctx context.Context, txn entity.Transaction) error {
	if !txn.Kind.Valid() {
		return srvErrors.NewValidationError(fmt.Sprintf("invalid transaction kind: %q", txn.Kind))
	}

	txn.RecomputeHash()

	return t.st.WithTx(ctx, func(ctx context.Context) error {
		return t.st.UpdateTransaction(ctx, txn)
	})
}

func (t *TransactionService) Delete(ctx context.Context, id int64) error {
	return t.st.WithTx(ctx, func(ctx context.Context) error {
		return t.st.DeleteTransaction(ctx, id)
	})
}
