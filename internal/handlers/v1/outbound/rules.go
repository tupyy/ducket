package outbound

import (
	"fmt"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

// NewRules creates a new Rules response structure from a slice of entity.Rule.
func NewRules(rules []entity.Rule) v1.Rules {
	r := v1.Rules{
		Rules: make([]v1.Rule, 0, len(rules)),
	}
	for _, rr := range rules {
		r.Rules = append(r.Rules, NewRule(rr))
	}

	r.Total = len(rules)
	return r
}

// NewRule creates a new Rule response structure from an entity.Rule.
func NewRule(rule entity.Rule) v1.Rule {
	labels := make([]v1.Label, 0, len(rule.Labels))
	r := v1.Rule{
		Href:    fmt.Sprintf("%s/rules/%s", apiV1, rule.Name),
		Name:    rule.Name,
		Pattern: ptr(rule.Pattern),
	}

	for _, label := range rule.Labels {
		labels = append(labels, v1.Label{
			Key:   label.Key,
			Value: label.Value,
			Href:  fmt.Sprintf("/api/v1/labels/%d", label.ID),
		})
	}

	r.Labels = &labels

	return r
}
