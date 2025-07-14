package entity

import "time"

type Label struct {
	ID                int
	Key               string
	Value             string
	CreatedAt         time.Time
	Rules             []string
	CountTransactions int
}

func (l Label) Equal(anotherLabel Label) bool {
	return l.Key == anotherLabel.Key && l.Value == anotherLabel.Value
}

type Rule struct {
	Name              string
	Pattern           string
	CreatedAt         time.Time
	CountTransactions int
	Labels            []Label
}

// NewRule creates a new Rule entity with the specified name, pattern, and associated tags.
func NewRule(name, pattern string, labels ...Label) Rule {
	_labels := make([]Label, 0, len(labels))
	for _, l := range labels {
		_labels = append(_labels, l)
	}

	return Rule{
		Name:    name,
		Pattern: pattern,
		Labels:  _labels,
	}
}
