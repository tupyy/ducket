package outbound

import (
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Tags struct {
	Total int   `json:"total"`
	Tags  []Tag `json:"tags"`
}

type Tag struct {
	HRef  string `json:"href"`
	Value string `json:"value"`
	Rules []Rule `json:"rules,omitempty"`
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

func NewTags(tags []string, rules []entity.Rule) Tags {
	mtags := Tags{
		Tags: make([]Tag, 0),
	}
	r := make(map[string]Tag)

	for _, rule := range rules {
		for _, t := range rule.Tags {
			if tag, ok := r[t]; ok {
				tag.Rules = append(tag.Rules, NewRule(rule))
				r[t] = tag
				continue
			}
			r[t] = Tag{
				HRef:  fmt.Sprintf("/api/v1/tags/%s", t),
				Value: t,
				Rules: []Rule{
					NewRule(rule),
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
	mtags.Total = len(r)

	return mtags
}
