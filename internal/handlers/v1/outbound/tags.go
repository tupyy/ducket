package outbound

import (
	"fmt"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type Labels struct {
	Total  int             `json:"total"`
	Labels []LabelWithMeta `json:"labels"`
}

type LabelWithMeta struct {
	HRef         string `json:"href"`
	Key          string `json:"key"`
	Value        string `json:"value"`
	Rules        []Rule `json:"rules,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Transactions int    `json:"transactions,omitempty"`
}

// NewLabel creates a new Label response structure with the given key, value and associated rules.
func NewLabel(key, value string, rules ...entity.Rule) LabelWithMeta {
	label := LabelWithMeta{
		HRef:  fmt.Sprintf("/api/v1/labels/%s/%s", key, value),
		Key:   key,
		Value: value,
		Rules: make([]Rule, 0),
	}

	for _, rule := range rules {
		label.Rules = append(label.Rules, NewRule(rule))
	}

	return label
}

// NewLabels creates a new Labels response structure from a slice of entity.Label.
func NewLabels(labels []entity.Label) Labels {
	mlabels := Labels{
		Labels: make([]LabelWithMeta, 0),
	}

	for _, eLabel := range labels {
		l := LabelWithMeta{
			HRef:         fmt.Sprintf("/api/v1/labels/%d", eLabel.ID),
			Key:          eLabel.Key,
			Value:        eLabel.Value,
			Transactions: eLabel.CountTransactions,
			CreatedAt:    eLabel.CreatedAt.Format(time.RFC3339),
			Rules:        make([]Rule, 0),
		}
		for _, rule := range eLabel.Rules {
			l.Rules = append(l.Rules, Rule{Name: rule, HRef: fmt.Sprintf("/api/v1/rules/%s", rule)})
		}
		mlabels.Labels = append(mlabels.Labels, l)
	}

	mlabels.Total = len(mlabels.Labels)

	return mlabels
}
