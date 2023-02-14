package parser

import (
	"context"

	"github.com/tupyy/finance/internal/entity"
)

func Parse(ctx context.Context, transactions []*entity.Transaction, rules []entity.Rule) []*entity.Transaction {
	for i := 0; i < len(transactions); i++ {
		t := transactions[i]
		// try to match each rule. each matched rule adds its own fields and labels.
		for _, rule := range rules {
			if t.Match(rule) {
				t.AddLabels(rule.Labels)
				t.AddFields(rule.Fields)
			}
		}
	}
	return transactions
}
