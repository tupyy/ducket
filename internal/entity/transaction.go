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

func (k TransactionKind) Valid() bool {
	return k == DebitTransaction || k == CreditTransaction
}

type Transaction struct {
	ID        int64           `json:"id"`
	Hash      string          `json:"hash"`
	Date      time.Time       `json:"date"`
	Account   int64           `json:"account"`
	Kind      TransactionKind `json:"kind"`
	Amount    float64         `json:"amount"`
	Content   string          `json:"content"`
	Info      *string         `json:"info,omitempty"`
	Recipient *string         `json:"recipient,omitempty"`
	Tags      []string        `json:"tags"`
}

func NewTransaction(kind TransactionKind, account int64, date time.Time, amount float64, content string) *Transaction {
	hash := fmt.Sprintf("%s|%d|%s|%f|%s", kind, account, date.UTC().Format(time.RFC3339Nano), amount, content)
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

func (t *Transaction) RecomputeHash() {
	hash := fmt.Sprintf("%s|%d|%s|%f|%s", t.Kind, t.Account, t.Date.UTC().Format(time.RFC3339Nano), t.Amount, t.Content)
	h := sha256.New()
	h.Write([]byte(hash))
	t.Hash = fmt.Sprintf("%x", h.Sum(nil))
}
