package models

import "git.tls.tupangiu.ro/cosmin/finante/internal/entity"

type RuleRows map[string][]RuleRow

func (r RuleRows) ToEntity() []entity.Rule {
	rules := []entity.Rule{}
	for _, v := range r {
		rule := entity.Rule{Tags: []string{}}
		for _, vv := range v {
			rule.Name = vv.ID
			rule.Pattern = vv.Pattern
			if vv.Tag != nil {
				rule.Tags = append(rule.Tags, *vv.Tag)
			}
		}
		rules = append(rules, rule)
	}

	return rules
}

func (r RuleRows) Add(ruleRow RuleRow) {
	rows, ok := r[ruleRow.ID]
	if !ok {
		r[ruleRow.ID] = []RuleRow{ruleRow}
		return
	}
	r[ruleRow.ID] = append(rows, ruleRow)
}

type RuleRow struct {
	ID      string
	Pattern string
	Tag     *string `db:"value"`
}
