package models

import (
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Tags map[string][]Tag

func (tags Tags) ToEntity() []entity.Tag {
	tt := []entity.Tag{}
	for _, v := range tags {
		tag := entity.Tag{RuleIDs: []string{}}
		for _, vv := range v {
			tag.Value = vv.Value
			tag.RuleIDs = append(tag.RuleIDs, vv.RuleID)
		}
		tt = append(tt, tag)
	}

	return tt
}

func (tags Tags) Add(t Tag) {
	rows, ok := tags[t.Value]
	if !ok {
		tags[t.Value] = []Tag{t}
		return
	}
	tags[t.Value] = append(rows, t)
}

type Tag struct {
	Value  string
	RuleID string
}
