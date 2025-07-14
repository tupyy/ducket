package inbound

import "git.tls.tupangiu.ro/cosmin/finante/internal/entity"

type TagForm struct {
	Value string `form:"value" json:"value" binding:"required"`
}

type LabelForm struct {
	Key   string `form:"key" json:"key" binding:"required"`
	Value string `form:"value" json:"value" binding:"required"`
}

// ToEntity converts a LabelForm to an entity.Label
func (lf LabelForm) ToEntity() entity.Label {
	return entity.Label{
		Key:   lf.Key,
		Value: lf.Value,
	}
}
