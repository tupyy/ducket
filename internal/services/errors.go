package services

import "fmt"

type ErrResourceNotFound struct {
	error
}

func NewErrTransactionNotFound(id int) *ErrResourceNotFound {
	return &ErrResourceNotFound{fmt.Errorf("transaction %d not found", id)}
}

func NewErrRuleNotFound(id int) *ErrResourceNotFound {
	return &ErrResourceNotFound{fmt.Errorf("rule %d not found", id)}
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
