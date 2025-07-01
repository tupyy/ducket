package inbound

import (
	"fmt"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"github.com/go-playground/validator/v10"
)

const (
	queryDateFormat = "02/01/2006"
)

type TransactionForm struct {
	Kind    string            `form:"kind" json:"kind" binding:"required"`
	Date    string            `form:"date" json:"date" binding:"required"`
	Content string            `form:"content" json:"content" binding:"required"`
	Amount  float32           `form:"amount" json:"amount" binding:"required"`
	Account int64             `form:"account" json:"account" binding:"required"`
	Tags    map[string]string `form:"tags" json:"tags" binding:"required"`
}

// Entity converts a TransactionForm to an entity.Transaction for business logic processing.
func (t TransactionForm) Entity() (entity.Transaction, error) {
	date, err := time.Parse(queryDateFormat, t.Date)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("unable to parse transaction date: %w", err)
	}
	te := entity.NewTransaction(entity.TransactionKind(t.Kind), t.Account, date, t.Amount, t.Content)
	te.Tags = t.Tags

	return *te, nil
}

// TransactionFormValidation provides custom validation logic for TransactionForm structures.
// It implements the validator.StructLevel interface for complex validation rules.
func TransactionFormValidation(sl validator.StructLevel) {
	form := sl.Current().Interface().(TransactionForm)

	if form.Kind != "credit" && form.Kind != "debit" {
		sl.ReportError(form.Kind, "kind", "kind", "shoudl be credit or debit", "")
	}

	_, err := time.Parse(queryDateFormat, form.Date)
	if err != nil {
		sl.ReportError(form.Date, "date", "date", "format invalid. shoudl be 02/01/2006", "")
	}

	if len(form.Tags) == 0 {
		sl.ReportError(form.Tags, "tags", "tags", "ge 0", "")
	}

	for tag, rule := range form.Tags {
		if len(tag) > 20 {
			sl.ReportError(form.Tags, "tag", "tag", "lt 20", "")
		}
		if len(rule) > 20 {
			sl.ReportError(form.Tags, "tag", "tag", "rule_id should be less 20 chars", "")
		}
	}
}
