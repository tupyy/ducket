package models

import (
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type List struct {
	Total int `json:"total"`
}

type TagRef struct {
	HRef  string `json:"href"`
	Value string `json:"value"`
}

type RuleRef struct {
	HRef string `json:"href"`
	Name string `json:"name"`
}

func newRuleRef(rule entity.Rule) RuleRef {
	return RuleRef{
		HRef: fmt.Sprintf("/api/v1/rules/%s", rule.ID),
		Name: rule.Name,
	}
}

type Tags struct {
	List
	Tags []Tag `json:"tags"`
}

type Tag struct {
	HRef  string    `json:"href"`
	Value string    `json:"value"`
	Rules []RuleRef `json:"rules,omitempty"`
}

func NewTags(tags []string, rules []entity.Rule) Tags {
	mtags := Tags{
		Tags: make([]Tag, 0),
	}
	r := make(map[string]Tag)

	for _, rule := range rules {
		for _, t := range rule.Tags {
			if tag, ok := r[t]; ok {
				tag.Rules = append(tag.Rules, newRuleRef(rule))
				r[t] = tag
				continue
			}
			r[t] = Tag{
				HRef:  fmt.Sprintf("/api/v1/tags/%s", t),
				Value: t,
				Rules: []RuleRef{
					newRuleRef(rule),
				},
			}
		}
	}

	for _, eTag := range tags {
		if _, ok := r[eTag]; !ok {
			r[eTag] = Tag{
				HRef:  fmt.Sprintf("/api/v1/tags/%s", eTag),
				Value: eTag,
			}
		}
	}

	for _, v := range r {
		mtags.Tags = append(mtags.Tags, v)
	}
	mtags.List.Total = len(r)

	return mtags
}
