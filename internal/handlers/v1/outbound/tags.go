package outbound

import (
	"fmt"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Tags struct {
	Total int   `json:"total"`
	Tags  []Tag `json:"tags"`
}

type Tag struct {
	HRef         string `json:"href"`
	Value        string `json:"value"`
	Rules        []Rule `json:"rules,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Transactions int    `json:"transactions,omitempty"`
}

func NewTag(value string, rules ...entity.Rule) Tag {
	tag := Tag{
		HRef:  fmt.Sprintf("/api/v1/tags/%s", value),
		Value: value,
		Rules: make([]Rule, 0),
	}

	for _, rule := range rules {
		tag.Rules = append(tag.Rules, NewRule(rule))
	}

	return tag
}

func NewTags(tags []entity.Tag, rules []entity.Rule) Tags {
	mtags := Tags{
		Tags: make([]Tag, 0),
	}
	r := make(map[string]Tag)

	for _, rule := range rules {
		for _, t := range rule.Tags {
			if tag, ok := r[t]; ok {
				tag.Rules = append(tag.Rules, Rule{Name: rule.Name, HRef: fmt.Sprintf("/api/v1/rules/%s", rule.Name)})
				r[t] = tag
				continue
			}
			r[t] = Tag{
				HRef:  fmt.Sprintf("/api/v1/tags/%s", t),
				Value: t,
				Rules: []Rule{
					{Name: rule.Name, HRef: fmt.Sprintf("/api/v1/rules/%s", rule.Name)},
				},
			}
		}
	}

	for _, eTag := range tags {
		if _, ok := r[eTag.Value]; !ok {
			r[eTag.Value] = Tag{
				HRef:         fmt.Sprintf("/api/v1/tags/%s", eTag.Value),
				Value:        eTag.Value,
				Transactions: eTag.CountTransactions,
				CreatedAt:    eTag.CreatedAt.Format(time.RFC3339),
			}
		} else {
			tt := r[eTag.Value]
			tt.CreatedAt = eTag.CreatedAt.Format(time.RFC3339)
			tt.Transactions = eTag.CountTransactions
			r[eTag.Value] = tt
		}
	}

	for _, v := range r {
		mtags.Tags = append(mtags.Tags, v)
	}
	mtags.Total = len(r)

	return mtags
}
