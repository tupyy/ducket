package entity

import (
	"crypto/sha256"
	"fmt"
)

type Tag struct {
	Value  string
	RuleID string
}

type Tags map[string]Tag

func (t Tags) Add(tag Tag) {
	t[tag.Value] = tag
}

func (t Tags) Has(tag string) bool {
	_, ok := t[tag]
	return ok
}

func (t Tags) Remove(tag string) {
	delete(t, tag)
}

type Rule struct {
	ID      string
	Name    string
	Pattern string
	Tags    Tags
}

func NewRule(name, pattern string, tags ...Tag) Rule {
	idString := fmt.Sprintf("%s%s", name, pattern)
	h := sha256.New()
	h.Write([]byte(idString))
	id := fmt.Sprintf("%x", h.Sum(nil))

	_tags := make(Tags)
	for _, tag := range tags {
		_tags.Add(tag)
	}

	return Rule{
		ID:      id[:100],
		Name:    name,
		Pattern: pattern,
		Tags:    _tags,
	}
}
