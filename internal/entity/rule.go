package entity

type Rule struct {
	ID      string
	Name    string
	Pattern string
	Tags    []string
}

func NewRule(id, name, pattern string, tags ...string) Rule {
	_tags := make([]string, 0, len(tags))
	for _, tag := range tags {
		_tags = append(_tags, tag)
	}

	return Rule{
		ID:      id,
		Name:    name,
		Pattern: pattern,
		Tags:    _tags,
	}
}
