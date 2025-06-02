package entity

const (
	idLength = 12
)

type Tag struct {
	Value   string
	RuleIDs []string
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

func NewRule(id, name, pattern string, tags ...Tag) Rule {
	_tags := make(Tags)
	for _, tag := range tags {
		_tags.Add(tag)
	}

	return Rule{
		ID:      id,
		Name:    name,
		Pattern: pattern,
		Tags:    _tags,
	}
}
