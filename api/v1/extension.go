package v1

import "git.tls.tupangiu.ro/cosmin/finante/internal/entity"

func NewTransactionFromEntity(t entity.Transaction) Transaction {
	txn := Transaction{
		Id:      t.ID,
		Hash:    t.Hash,
		Date:    t.Date,
		Account: t.Account,
		Kind:    TransactionKind(t.Kind),
		Amount:  t.Amount,
		Content: t.Content,
		Tags:    t.Tags,
	}
	if t.Info != nil {
		txn.Info = t.Info
	}
	if t.Recipient != nil {
		txn.Recipient = t.Recipient
	}
	if txn.Tags == nil {
		txn.Tags = []string{}
	}
	return txn
}

func NewRuleFromEntity(r entity.Rule) Rule {
	rule := Rule{
		Id:     r.ID,
		Name:   r.Name,
		Filter: r.Filter,
		Tags:   r.Tags,
	}
	if !r.CreatedAt.IsZero() {
		rule.CreatedAt = &r.CreatedAt
	}
	if rule.Tags == nil {
		rule.Tags = []string{}
	}
	return rule
}
