package entity

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type TransactionKind string

const (
	DebitTransaction  TransactionKind = "debit"
	CreditTransaction TransactionKind = "credit"
)

type Transaction struct {
	ID         int
	Kind       TransactionKind
	Date       time.Time
	RawContent string
	Hash       string
	Amount     float32
	// Tags holds the tags applied on this transaction
	// the key is the tag and the value the rule_id which applied the tag
	Tags map[string]string
}

func NewTransaction(kind TransactionKind, date time.Time, sum float32, rawContent string) *Transaction {
	hash := fmt.Sprintf("%s%s%f%s", kind, date, sum, rawContent)
	h := sha256.New()
	h.Write([]byte(hash))

	return &Transaction{
		RawContent: rawContent,
		Date:       date,
		Hash:       fmt.Sprintf("%x", h.Sum(nil)),
		Kind:       kind,
		Amount:     sum,
		Tags:       make(map[string]string),
	}
}

func (t *Transaction) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s | %s | %.2f€ | ", t.Date, t.RawContent, t.Amount)
	for k := range t.Tags {
		fmt.Fprintf(&sb, "%s ", k)
	}
	return sb.String()
}
