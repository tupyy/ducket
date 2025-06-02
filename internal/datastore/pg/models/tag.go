package models

import (
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Tags map[string][]Tag

func (tags Tags) ToEntity() []entity.Tag {
	return []entity.Tag{}
}

type Tag struct {
	Value  string
	RuleID string `db:"rule_id"`
}
