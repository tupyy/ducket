package entity

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

type TransactionKind string

const (
	DebitTransaction  TransactionKind = "debit"
	CreditTransaction TransactionKind = "credit"
)

type Transaction struct {
	ID         string
	Kind       TransactionKind
	Date       time.Time
	RawContent string
	Amount     float32
	Tags       Tags
}

func NewTransaction(kind TransactionKind, date time.Time, sum float32, rawContent string) *Transaction {
	idString := fmt.Sprintf("%s%s%f%s", kind, date, sum, rawContent)
	h := sha256.New()
	h.Write([]byte(idString))
	id := fmt.Sprintf("%x", h.Sum(nil))

	return &Transaction{
		ID:         string(id[:100]),
		RawContent: rawContent,
		Date:       date,
		Kind:       kind,
		Amount:     sum,
		Tags:       make(Tags),
	}
}

func (t *Transaction) Match(rule Rule) bool {
	pattern, err := regexp.Compile(rule.Pattern)
	if err != nil {
		zap.S().Errorf("rule pattern %q does not compile", rule.Pattern)
		return false
	}
	return pattern.MatchString(t.RawContent)
}

func (t *Transaction) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s | %s | %.2f€ | ", t.Date, t.RawContent, t.Amount)
	for k := range t.Tags {
		fmt.Fprintf(&sb, "%s ", k)
	}
	return sb.String()
}
