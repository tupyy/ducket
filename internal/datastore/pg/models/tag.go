package models

import (
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

// store tags by value because we can have the same tag associated with multiple rules
type Tags map[string][]Tag

// ToEntity converts the database model tags to entity tags.
func (tags Tags) ToEntity() []entity.Tag {
	tt := []entity.Tag{}
	for k, tagRows := range tags {
		tag := entity.Tag{Value: k}
		for _, r := range tagRows {
			tag.CreatedAt = r.CreatedAt
			if r.RuleID != nil {
				tag.Rules = append(tag.Rules, *r.RuleID)
			}
		}
		tt = append(tt, tag)
	}

	return tt
}

// Add appends a tag to the collection, grouping by tag value.
func (tags Tags) Add(t Tag) {
	rows, ok := tags[t.Value]
	if !ok {
		tags[t.Value] = []Tag{t}
		return
	}
	tags[t.Value] = append(rows, t)
}

type Tag struct {
	Value     string
	CreatedAt time.Time `db:"created_at"`
	RuleID    *string
}

type TransactionCountRow struct {
	Value  string
	RuleID string
	Count  int
}
