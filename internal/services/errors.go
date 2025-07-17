package services

import "fmt"

type ErrResourceNotFound struct {
	error
}

func NewErrTransactionNotFound(id int) *ErrResourceNotFound {
	return &ErrResourceNotFound{fmt.Errorf("transaction %d not found", id)}
}

func NewErrTransactionNotFoundByHash(hash string) *ErrResourceNotFound {
	return &ErrResourceNotFound{fmt.Errorf("transaction with hash %s not found", hash)}
}

func NewErrRuleNotFound(id string) *ErrResourceNotFound {
	return &ErrResourceNotFound{fmt.Errorf("rule %s not found", id)}
}

type ErrResourceExistsAlready struct {
	error
}

func NewErrTransactionExistsAlready(id int) *ErrResourceExistsAlready {
	return &ErrResourceExistsAlready{fmt.Errorf("transaction %d already exists", id)}
}

func IsErrResourceNotFound(err error) bool {
	_, ok := err.(*ErrResourceNotFound)
	return ok
}
