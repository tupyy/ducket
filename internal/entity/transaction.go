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
	ID         int64
	Kind       TransactionKind
	Date       time.Time
	Account    int64
	RawContent string
	Hash       string
	Amount     float32
	// Labels holds the tags applied on this transaction
	// the key is the tag and the value the rule_id which applied the tag
	Labels map[int]string
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
		Labels:     make(map[int]string),
	}
}

// Holds how many transactions are associated with this tag and rule
type TransactionStat struct {
	LabelID int
	RuleID  string
	Count   int
}
