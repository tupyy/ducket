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

// NewTag creates a new Tag response structure with the given value and associated rules.
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

// NewTags creates a new Tags response structure from a slice of entity.Tag.
func NewTags(tags []entity.Tag) Tags {
	mtags := Tags{
		Tags: make([]Tag, 0),
	}

	for _, eTag := range tags {
		t := Tag{
			HRef:         fmt.Sprintf("/api/v1/tags/%s", eTag.Value),
			Value:        eTag.Value,
			Transactions: eTag.CountTransactions,
			CreatedAt:    eTag.CreatedAt.Format(time.RFC3339),
			Rules:        make([]Rule, 0),
		}
		for _, rule := range eTag.Rules {
			t.Rules = append(t.Rules, Rule{Name: rule, HRef: fmt.Sprintf("/api/v1/rules/%s", rule)})
		}
		mtags.Tags = append(mtags.Tags, t)
	}

	mtags.Total = len(mtags.Tags)

	return mtags
}
