package entity

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type TransactionKind string

const (
	DebitTransaction  TransactionKind = "debit"
	CreditTransaction TransactionKind = "credit"
)

type Transaction struct {
	ID        int64
	Hash      string
	Date      time.Time
	Account   int64
	Kind      TransactionKind
	Amount    float64
	Content   string
	Info      *string
	Recipient *string
	Tags      []string
}

func NewTransaction(kind TransactionKind, account int64, date time.Time, amount float64, content string) *Transaction {
	hash := fmt.Sprintf("%s%s%f%s", kind, date, amount, content)
	h := sha256.New()
	h.Write([]byte(hash))

	return &Transaction{
		Kind:    kind,
		Date:    date,
		Account: account,
		Hash:    fmt.Sprintf("%x", h.Sum(nil)),
		Amount:  amount,
		Content: content,
	}
}
