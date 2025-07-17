package entity

type RelationshipType int

const (
	RelationshipLabelTransaction RelationshipType = iota
	RelationshipLabelRuleTransaction
	RelationshipLabelRule
)

type Relationship struct {
	Kind          RelationshipType
	LabelID       int
	TransactionID int
	RuleID        string
}
