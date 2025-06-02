package models

import "git.tls.tupangiu.ro/cosmin/finante/internal/entity"

type RuleRows map[string][]RuleRow

func (r RuleRows) ToEntity() []entity.Rule {
	rules := []entity.Rule{}
	for _, v := range r {
		rule := entity.Rule{Tags: make(entity.Tags)}
		for _, vv := range v {
			rule.ID = vv.ID
			rule.Name = vv.Name
			rule.Pattern = vv.Pattern
			if vv.TagValue != nil {
				rule.Tags.Add(entity.Tag{Value: *vv.TagValue, RuleIDs: []string{vv.ID}})
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
	ID       string
	Name     string
	Pattern  string
	TagValue *string `db:"value"`
}
