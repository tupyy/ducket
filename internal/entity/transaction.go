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

// LabelAssociation keeps association between a label applied on a transaction and the ruleID which matched the transaction.
// If RuleID is nil, it means the label has been manually applied.
type LabelAssociation struct {
	Label  Label
	RuleID *string
}

type Transaction struct {
	ID         int64
	Kind       TransactionKind
	Date       time.Time
	Account    int64
	RawContent string
	Hash       string
	Amount     float32
	Info       *string
	Labels     []LabelAssociation
}

// NewTransaction creates a new Transaction entity with the specified kind, date, amount, and raw content.
func NewTransaction(kind TransactionKind, account int64, date time.Time, sum float32, rawContent string) *Transaction {
	hash := fmt.Sprintf("%s%s%f%s", kind, date, sum, rawContent)
	h := sha256.New()
	h.Write([]byte(hash))

	return &Transaction{
		RawContent: rawContent,
		Date:       date,
		Account:    account,
		Hash:       fmt.Sprintf("%x", h.Sum(nil)),
		Kind:       kind,
		Amount:     sum,
		Labels:     []LabelAssociation{},
	}
}
