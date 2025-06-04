package models

type Tags map[string][]Tag

func (tags Tags) ToEntity() []string {
	tt := []string{}
	for tag := range tags {
		tt = append(tt, tag)
	}

	return tt
}

func (tags Tags) Add(t Tag) {
	rows, ok := tags[t.Value]
	if !ok {
		tags[t.Value] = []Tag{t}
		return
	}
	tags[t.Value] = append(rows, t)
}

type Tag struct {
	Value  string
	RuleID *string
}
