package outbound

import (
	"fmt"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

const (
	apiV1 = "/api/v1"
)

type TransactionLabelAssociation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Href  string `json:"href"`
	Rule  string `json:"rule"`
}

type Transaction struct {
	Href        string                        `json:"href"`
	Kind        string                        `json:"kind"`
	Date        string                        `json:"date"`
	Account     int64                         `json:"account"`
	Description string                        `json:"description"`
	Amount      float32                       `json:"amount"`
	Labels      []TransactionLabelAssociation `json:"labels"`
}

// FromEntity converts an entity.Transaction to an outbound Transaction model
// suitable for API responses.
func FromEntity(t entity.Transaction) Transaction {
	transaction := Transaction{
		Href:        fmt.Sprintf("%s/transactions/%d", apiV1, t.ID),
		Kind:        string(t.Kind),
		Date:        t.Date.Format(time.RFC3339),
		Description: t.RawContent,
		Account:     t.Account,
		Amount:      t.Amount,
		Labels:      make([]TransactionLabelAssociation, 0, len(t.Labels)),
	}

	// TODO: Need to lookup label information by ID to get key-value pairs
	// The t.Labels is now map[int]string where key is label_id and value is rule_id
	// This requires service layer access to resolve label IDs to key-value pairs
	// For now, we'll leave this empty until we can properly implement the lookup
	for labelID, ruleID := range t.Labels {
		transaction.Labels = append(transaction.Labels, TransactionLabelAssociation{
			Key:   fmt.Sprintf("label_%d", labelID), // Placeholder - needs proper lookup
			Value: fmt.Sprintf("value_%d", labelID), // Placeholder - needs proper lookup
			Href:  fmt.Sprintf("%s/labels/%d", apiV1, labelID),
			Rule:  ruleID,
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
