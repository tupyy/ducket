package entity

import (
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
	// kind of transaction
	Kind TransactionKind `json:"kind"`
	// date of transaction
	Date time.Time `json:"date"`
	// RawContent holds the raw content of the record
	RawContent string `json:"content"`
	// sum of the transaction
	Sum float32 `json:"sum"`
	// Labels holds the labels to be applied.
	Labels map[string]string `json:"labels,omitempty"`
}

func NewTransaction(kind TransactionKind, date time.Time, sum float32, rawContent string) *Transaction {
	return &Transaction{
		RawContent: rawContent,
		Date:       date,
		Kind:       kind,
		Sum:        sum,
		Labels:     make(map[string]string),
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

func (t *Transaction) AddLabels(labels map[string]string) {
	for k, v := range labels {
		t.Labels[k] = v
	}
}

func (t *Transaction) Matched() bool {
	return len(t.Labels) != 0
}

func (t *Transaction) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s | %s | %.2f€ | ", t.Date, t.RawContent, t.Sum)
	for k, v := range t.Labels {
		fmt.Fprintf(&sb, "%s:%s ", k, v)
	}
	return sb.String()
}
