package models

import (
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

// store labels by value because we can have the same label associated with multiple rules
type Labels map[int][]Label

// ToEntity converts the database model labels to entity labels.
func (labels Labels) ToEntity() []entity.Label {
	tt := []entity.Label{}
	for id, rows := range labels {
		label := entity.Label{ID: id}
		for _, r := range rows {
			label.CreatedAt = r.CreatedAt
			label.Key = r.Key
			label.Value = r.Value
			if r.RuleID != nil {
				label.Rules = append(label.Rules, *r.RuleID)
			}
		}
		tt = append(tt, label)
	}

	return tt
}

// Add appends a label to the collection, grouping by label value.
func (labels Labels) Add(t Label) {
	rows, ok := labels[t.ID]
	if !ok {
		labels[t.ID] = []Label{t}
		return
	}
	labels[t.ID] = append(rows, t)
}

type Label struct {
	ID        int
	Key       string
	Value     string
	CreatedAt time.Time `db:"created_at"`
	RuleID    *string
}

type TransactionCountRow struct {
	LabelID int
	RuleID  string
	Count   int
}
