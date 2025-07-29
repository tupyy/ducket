package inbound

import "git.tls.tupangiu.ro/cosmin/finante/internal/entity"

// LabelForm represents the HTTP request body structure for creating or updating labels.
// It contains the key-value pair that defines a label for categorizing transactions.
type LabelForm struct {
	Key   string `form:"key" json:"key" binding:"required"`
	Value string `form:"value" json:"value" binding:"required"`
}

// ToEntity converts a LabelForm to an entity.Label for business logic processing.
// This method transforms the HTTP request data into the internal domain model representation.
func (lf LabelForm) ToEntity() entity.Label {
	return entity.Label{
		Key:   lf.Key,
		Value: lf.Value,
	}
}
