package parser

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

func Parse(ctx context.Context, transactions []*entity.Transaction, rules []entity.Rule) []*entity.Transaction {
	for i := 0; i < len(transactions); i++ {
		t := transactions[i]
		// try to match each rule. each matched rule adds its own labels.
		for _, rule := range rules {
			if t.Match(rule) {
				t.AddLabels(rule.Labels)
			}
		}
	}
	return transactions
}
