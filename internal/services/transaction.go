package services

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
)

type TransactionService struct {
	st *store.Store
}

func NewTransactionService(st *store.Store) *TransactionService {
	return &TransactionService{st: st}
}

func (t *TransactionService) List(ctx context.Context, filter string, limit, offset int) ([]entity.Transaction, error) {
	return t.st.ListTransactions(ctx, filter, limit, offset)
}

func (t *TransactionService) Get(ctx context.Context, id int64) (*entity.Transaction, error) {
	txn, err := t.st.GetTransaction(ctx, id)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return nil, NewErrTransactionNotFound(int(id))
	}
	return txn, nil
}

func (t *TransactionService) Create(ctx context.Context, txn entity.Transaction) (entity.Transaction, error) {
	existing, _ := t.st.GetTransactionByHash(ctx, txn.Hash)
	if existing != nil {
		return *existing, NewErrTransactionExistsAlready(int(existing.ID))
	}

	var id int64
	if err := t.st.WithTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = t.st.CreateTransaction(ctx, txn)
		return err
	}); err != nil {
		return txn, err
	}

	txn.ID = id
	return txn, nil
}

func (t *TransactionService) Update(ctx context.Context, txn entity.Transaction) error {
	existing, err := t.st.GetTransaction(ctx, txn.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return NewErrTransactionNotFound(int(txn.ID))
	}

	return t.st.WithTx(ctx, func(ctx context.Context) error {
		return t.st.UpdateTransaction(ctx, txn)
	})
}

func (t *TransactionService) Delete(ctx context.Context, id int64) error {
	return t.st.WithTx(ctx, func(ctx context.Context) error {
		return t.st.DeleteTransaction(ctx, id)
	})
}
