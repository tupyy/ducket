package models

import (
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type RuleRows map[string][]RuleRow

// ToEntity converts the database model rules to entity rules.
func (r RuleRows) ToEntity() []entity.Rule {
	rules := []entity.Rule{}
	for _, v := range r {
		rule := entity.Rule{Labels: []entity.Label{}}
		for _, vv := range v {
			rule.Name = vv.ID
			rule.Pattern = vv.Pattern
			rule.CreatedAt = vv.CreatedAt
			if vv.LabelKey != nil && vv.LabelValue != nil {
				label := entity.Label{
					ID:    *vv.LabelID,
					Key:   *vv.LabelKey,
					Value: *vv.LabelValue,
				}
				rule.Labels = append(rule.Labels, label)
			}
		}
		rules = append(rules, rule)
	}

	return rules
}

// Add appends a rule row to the collection, grouping by rule ID.
func (r RuleRows) Add(ruleRow RuleRow) {
	rows, ok := r[ruleRow.ID]
	if !ok {
		r[ruleRow.ID] = []RuleRow{ruleRow}
		return
	}
	r[ruleRow.ID] = append(rows, ruleRow)
}

type RuleRow struct {
	ID         string `db:"rule_id"`
	Pattern    string
	CreatedAt  time.Time
	LabelID    *int    `db:"label_id"`
	LabelKey   *string `db:"key"`
	LabelValue *string `db:"value"`
}
