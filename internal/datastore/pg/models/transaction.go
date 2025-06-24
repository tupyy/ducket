package models

import (
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Transactions map[int64][]Transaction

// Add appends a transaction to the collection, grouping by transaction ID.
func (tt Transactions) Add(t Transaction) {
	rows, ok := tt[t.ID]
	if !ok {
		tt[t.ID] = []Transaction{t}
		return
	}
	tt[t.ID] = append(rows, t)
}

// ToEntity converts the database model transactions to entity transactions.
func (tt Transactions) ToEntity() []entity.Transaction {
	transactions := []entity.Transaction{}
	for _, v := range tt {
		transaction := entity.Transaction{Tags: make(map[string]string)}
		for _, vv := range v {
			transaction.ID = vv.ID
			transaction.Date = vv.Date
			transaction.RawContent = vv.Content
			transaction.Amount = vv.Amount
			transaction.Hash = vv.Hash
			transaction.Kind = entity.TransactionKind(vv.Kind)
			if vv.Tag != nil {
				transaction.Tags[*vv.Tag] = *vv.RuleID
			}
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}

type Transaction struct {
	ID      int64
	Kind    string
	Date    time.Time
	Content string
	Amount  float32
	Hash    string
	Tag     *string `db:"tag_id"`
	RuleID  *string `db:"rule_id"`
}
