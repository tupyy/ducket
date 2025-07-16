package services

import "fmt"

type ErrTransactionNotFound struct {
	error
}

func NewErrTransactionNotFound(id int) *ErrTransactionNotFound {
	return &ErrTransactionNotFound{fmt.Errorf("transaction %d not found", id)}
}

func NewErrTransactionNotFoundByHash(hash string) *ErrTransactionNotFound {
	return &ErrTransactionNotFound{fmt.Errorf("transaction with hash %s not found", hash)}
}
