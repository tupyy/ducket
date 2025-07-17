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

func NewLabeRuleRelationship(labelID int, ruleID string) Relationship {
	return Relationship{Kind: RelationshipLabelRule, LabelID: labelID, RuleID: ruleID}
}

func NewLabeRuleTransactionRelationship(labelID int, ruleID string, transactionID int) Relationship {
	return Relationship{Kind: RelationshipLabelRuleTransaction, LabelID: labelID, RuleID: ruleID, TransactionID: transactionID}
}

func NewLabelTransaction(labelID int, transactionID int) Relationship {
	return Relationship{Kind: RelationshipLabelTransaction, LabelID: labelID, TransactionID: transactionID}
}
