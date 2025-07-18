package outbound

import (
	"fmt"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

// NewLabel creates a new Label response structure with the given key, value and associated rules.
func NewLabel(l entity.Label) v1.Label {
	rules := make([]v1.Rule, 0)
	label := v1.Label{
		Href:  fmt.Sprintf("/api/v1/labels/%d", l.ID),
		Key:   l.Key,
		Value: l.Value,
	}

	for _, rule := range l.Rules {
		rules = append(rules,
			v1.Rule{
				Href: fmt.Sprintf("/api/v1/rules/%s", rule),
				Name: rule,
			})
	}

	label.Rules = &rules

	return label
}

// NewLabels creates a new Labels response structure from a slice of entity.Label.
func NewLabels(labels []entity.Label) v1.Labels {
	mlabels := v1.Labels{
		Labels: make([]v1.Label, 0),
	}

	for _, label := range labels {
		mlabels.Labels = append(mlabels.Labels, NewLabel(label))
	}

	mlabels.Total = len(mlabels.Labels)

	return mlabels
}

func ptr(s string) *string {
	return &s
}
