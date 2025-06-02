package models

import (
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Transactions map[string][]Transaction

func (tt Transactions) Add(t Transaction) {
	rows, ok := tt[t.ID]
	if !ok {
		tt[t.ID] = []Transaction{t}
		return
	}
	tt[t.ID] = append(rows, t)
}

func (tt Transactions) ToEntity() []entity.Transaction {
	transactions := []entity.Transaction{}
	for _, v := range tt {
		transaction := entity.Transaction{Tags: make(entity.Tags)}
		for _, vv := range v {
			transaction.ID = vv.ID
			transaction.Date = vv.Date
			transaction.RawContent = vv.Content
			transaction.Amount = vv.Amount
			transaction.Kind = entity.TransactionKind(vv.Kind)
			if vv.Tag != nil {
				transaction.Tags.Add(entity.Tag{Value: *vv.Tag, RuleIDs: []string{*vv.RuleID}})
			}
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}

type Transaction struct {
	ID      string
	Kind    string
	Date    time.Time
	Content string
	Amount  float32
	Tag     *string `db:"tag_id"`
	RuleID  *string `db:"rule_id"`
}
