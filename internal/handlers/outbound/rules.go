package outbound

import (
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Rules struct {
	Rules []Rule `json:"rules"`
	Total int    `json:"total"`
}

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
	HRef    string `json:"href"`
	Name    string `json:"name"`
	Pattern string `json:"pattern,omitempty"`
	Tags    []Tag  `json:"tags,omitempty"`
}

func NewRule(rule entity.Rule) Rule {
	r := Rule{
		HRef:    fmt.Sprintf("/api/v1/rules/%s", rule.Name),
		Name:    rule.Name,
		Pattern: rule.Pattern,
	}

	for _, t := range rule.Tags {
		r.Tags = append(r.Tags, NewTag(t))
	}

	return r
}
