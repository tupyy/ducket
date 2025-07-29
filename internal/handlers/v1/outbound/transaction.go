package outbound

import (
	"fmt"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

const (
	apiV1 = "/api/v1"
)

// TransactionLabelAssociation represents a label attached to a transaction in API responses.
// It includes the label's key-value pair, HREF for the label resource, and optional
// reference to the rule that applied this label.
type TransactionLabelAssociation struct {
	Href    string `json:"href"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	RuleRef string `json:"rule_href,omitempty"`
}

// Transaction represents a financial transaction in API responses.
// It contains all transaction details including labels, formatted for client consumption
// with HREF links for related resources and RFC3339 formatted dates.
type Transaction struct {
	ID          int                           `json:"id"`
	Href        string                        `json:"href"`
	Kind        string                        `json:"kind"`
	Date        string                        `json:"date"`
	Account     int64                         `json:"account"`
	Description string                        `json:"description"`
	Amount      float32                       `json:"amount"`
	Info        *string                       `json:"info,omitempty"`
	Labels      []TransactionLabelAssociation `json:"labels"`
}

// FromEntity converts an entity.Transaction to an outbound Transaction model
// suitable for API responses, including proper HREF generation and date formatting.
func FromEntity(t entity.Transaction) Transaction {
	transaction := Transaction{
		ID:          int(t.ID),
		Href:        fmt.Sprintf("%s/transactions/%d", apiV1, t.ID),
		Kind:        string(t.Kind),
		Date:        t.Date.Format(time.RFC3339),
		Description: t.RawContent,
		Account:     t.Account,
		Amount:      t.Amount,
		Info:        t.Info,
		Labels:      make([]TransactionLabelAssociation, 0, len(t.Labels)),
	}

	for _, a := range t.Labels {
		tLabelAssociation := TransactionLabelAssociation{
			Key:   a.Label.Key,
			Value: a.Label.Value,
			Href:  fmt.Sprintf("/api/v1/labels/%d", a.Label.ID),
		}
		if a.RuleID != nil {
			tLabelAssociation.RuleRef = fmt.Sprintf("/api/v1/rules/%s", *a.RuleID)
		}
		transaction.Labels = append(transaction.Labels, tLabelAssociation)
	}

	return transaction
}

// Transactions represents a paginated collection of transactions in API responses.
// It includes metadata about the total count, date range filter applied, and the
// actual transaction items in the response.
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
