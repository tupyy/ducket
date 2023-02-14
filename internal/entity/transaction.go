package entity

import "time"

type TransactionKind int

const (
	DebitTransaction TransactionKind = iota
	CreditTransaction
)

type Transaction struct {
	// kind of transaction
	Kind TransactionKind
	// date of transaction
	Date time.Time
	// RawContent holds the raw content of the record
	RawContent string
	// sum of the transaction
	Sum float32
	// fields
	Fields map[string]string
	// Labels holds the labels to be applied.
	Labels []Label
}

type Label struct {
	Key   string
	Value string
}

func NewTransaction(kind TransactionKind, date time.Time, sum float32, rawContent string) *Transaction {
	return &Transaction{
		RawContent: rawContent,
		Date:       date,
		Kind:       kind,
		Sum:        sum,
		Fields:     make(map[string]string),
		Labels: []Label{
			{
				Key:   "category",
				Value: "unknown",
			},
		},
	}
}

func (t *Transaction) Match(rule Rule) bool {
	return rule.Pattern.MatchString(t.RawContent)
}

func (t *Transaction) AddFields(fields map[string]string) {
	for k, v := range fields {
		t.Fields[k] = v
	}
}

func (t *Transaction) AddLabels(labels []Label) {
	t.Labels = append(t.Labels, labels...)
}
