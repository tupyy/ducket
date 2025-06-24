package outbound

import (
	"fmt"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

const (
	apiV1 = "/api/v1"
)

type TransactionTagAssociation struct {
	Value string `json:"value"`
	Href  string `json:"href"`
	Rule  Rule   `json:"rule"`
}

type Transaction struct {
	Href        string                      `json:"href"`
	Kind        string                      `json:"kind"`
	Date        string                      `json:"date"`
	Description string                      `json:"description"`
	Amount      float32                     `json:"amount"`
	Tags        []TransactionTagAssociation `json:"tags"`
}

// FromEntity converts an entity.Transaction to an outbound Transaction model
// suitable for API responses.
func FromEntity(t entity.Transaction) Transaction {
	transaction := Transaction{
		Href:        fmt.Sprintf("%s/transactions/%d", apiV1, t.ID),
		Kind:        string(t.Kind),
		Date:        t.Date.Format(time.RFC3339),
		Description: t.RawContent,
		Amount:      t.Amount,
		Tags:        make([]TransactionTagAssociation, 0, len(t.Tags)),
	}

	for tag, ruleID := range t.Tags {
		transaction.Tags = append(transaction.Tags, TransactionTagAssociation{
			Value: tag,
			Href:  fmt.Sprintf("%s/tags/%s", apiV1, tag),
			Rule:  NewRule(entity.Rule{Name: ruleID}),
		})
	}

	return transaction
}

type Transactions struct {
	Total        int           `json:"total"`
	StartingDate string        `json:"start"`
	EndingDate   string        `json:"end"`
	Items        []Transaction `json:"items"`
}

// NewTransactions creates a new Transactions response structure with the given
// total count, time range, and transaction data.
func NewTransactions(total int, start, end time.Time) *Transactions {
	return &Transactions{
		Total:        total,
		StartingDate: start.Format("02/01/2006"),
		EndingDate:   end.Format("02/01/2006"),
		Items:        make([]Transaction, 0, total),
	}
}

type Summary struct {
	TotalAmmount float32            `json:"total"`
	StartingDate string             `json:"start"`
	EndingDate   string             `json:"end"`
	Items        map[string]float32 `json:"items"`
}
