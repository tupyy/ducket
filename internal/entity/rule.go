package entity

import "time"

type Tag struct {
	Value             string
	CreatedAt         time.Time
	Rules             []string
	CountTransactions int
}

type Rule struct {
	Name              string
	Pattern           string
	CreatedAt         time.Time
	CountTransactions int
	Tags              []string
}

// NewRule creates a new Rule entity with the specified name, pattern, and associated tags.
func NewRule(name, pattern string, tags ...string) Rule {
	_tags := make([]string, 0, len(tags))
	for _, tag := range tags {
		_tags = append(_tags, tag)
	}

	return Rule{
		Name:    name,
		Pattern: pattern,
		Tags:    _tags,
	}
}
