package outbound

import (
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Rules struct {
	Rules []Rule `json:"rules"`
	Total int    `json:"total"`
}

// NewRules creates a new Rules response structure from a slice of entity.Rule.
func NewRules(rules []entity.Rule) Rules {
	r := Rules{
		Rules: make([]Rule, 0),
	}
	for _, rr := range rules {
		r.Rules = append(r.Rules, NewRule(rr))
	}

	r.Total = len(rules)
	return r
}

type Rule struct {
	HRef    string  `json:"href"`
	Name    string  `json:"name"`
	Pattern string  `json:"pattern,omitempty"`
	Labels  []Label `json:"labels,omitempty"`
}

// NewRule creates a new Rule response structure from an entity.Rule.
func NewRule(rule entity.Rule) Rule {
	r := Rule{
		HRef:    fmt.Sprintf("%s/rules/%s", apiV1, rule.Name),
		Name:    rule.Name,
		Pattern: rule.Pattern,
		Labels:  make([]Label, 0, len(rule.Labels)),
	}

	for _, label := range rule.Labels {
		r.Labels = append(r.Labels, Label{
			Key:   label.Key,
			Value: label.Value,
			HRef:  fmt.Sprintf("/api/v1/labels/%d", label.ID),
		})
	}

	return r
}
