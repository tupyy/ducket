package store

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/kubev2v/assisted-migration-agent/pkg/filter"
)

var transactionColumns = map[string]string{
	"id":        "t.id",
	"hash":      "t.hash",
	"date":      "t.date",
	"account":   "t.account",
	"kind":      "t.kind",
	"amount":    "t.amount",
	"content":   "t.content",
	"info":      "t.info",
	"recipient": "t.recipient",
}

var transactionMapFn = filter.MapFunc(func(name string) (string, error) {
	col, ok := transactionColumns[name]
	if !ok {
		return "", fmt.Errorf("unknown transaction field: %s", name)
	}
	return col, nil
})

func ParseTransactionFilter(expr string) (sq.Sqlizer, error) {
	return filter.Parse([]byte(expr), transactionMapFn)
}

var ruleColumns = map[string]string{
	"id":     "id",
	"name":   "name",
	"filter": "filter",
}

var ruleMapFn = filter.MapFunc(func(name string) (string, error) {
	col, ok := ruleColumns[name]
	if !ok {
		return "", fmt.Errorf("unknown rule field: %s", name)
	}
	return col, nil
})

func ParseRuleFilter(expr string) (sq.Sqlizer, error) {
	return filter.Parse([]byte(expr), ruleMapFn)
}
