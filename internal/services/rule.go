package services

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

//go:generate go run github.com/ecordell/optgen -output zz_generated.rule_filter.go . RuleFilter
type RuleFilter struct {
	Id string `debugmap:"visible"`
}

// QueriesFn returns a slice of query filters based on the rule filter criteria.
func (rf *RuleFilter) QueriesFn() []pg.QueryFilter {
	qf := []pg.QueryFilter{}

	qf = append(qf,
		pg.RuleIDQueryFilter(rf.Id),
	)

	return qf
}

type RuleService struct {
	dt *pg.Datastore
}

// NewRuleService creates a new instance of RuleService with the provided datastore.
func NewRuleService(dt *pg.Datastore) *RuleService {
	return &RuleService{dt: dt}
}

// GetRules retrieves all rules from the database.
func (r *RuleService) GetRules(ctx context.Context) ([]entity.Rule, error) {
	return r.dt.QueryRules(ctx)
}

// GetRule retrieves a single rule by its name.
// Returns nil if no rule is found with the given name.
func (r *RuleService) GetRule(ctx context.Context, name string) (*entity.Rule, error) {
	rules, err := r.dt.QueryRules(ctx, NewRuleFilterWithOptions(WithId(name)).QueriesFn()...)
	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return nil, nil
	}

	return &rules[0], nil
}

// Create creates a new rule in the database.
// Returns an error if a rule with the same name already exists.
// Also ensures that all tags referenced by the rule exist in the database.
func (r *RuleService) Create(ctx context.Context, rule entity.Rule) error {
	existingRule, err := r.GetRule(ctx, rule.Name)
	if err != nil {
		return err
	}

	if existingRule != nil {
		return fmt.Errorf("rule %s already exists", rule.Name)
	}

	labelSrv := NewLabelService(r.dt)
	existingLabels, err := labelSrv.GetLabels(ctx)
	if err != nil {
		return err
	}

	return r.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		// write rule
		if err := w.WriteRule(ctx, rule, false); err != nil {
			return err
		}

		relationships := []entity.Relationship{}
		for _, label := range rule.Labels {
			found := false
			labelID := 0
			for _, existingLabel := range existingLabels {
				if label.Equal(existingLabel) {
					found = true
					labelID = existingLabel.ID
					break
				}
			}
			if !found {
				id, err := w.WriteLabel(ctx, label)
				if err != nil {
					return err
				}
				labelID = id
			}

			relationships = append(relationships, entity.NewLabeRuleRelationship(labelID, rule.Name))
		}

		return w.WriteRelationships(ctx, relationships)
	})
}

// UpdateOrCreate creates a new rule or updates an existing one based on the name.
// Returns true if an existing rule was updated, false if a new rule was created.
// Also ensures that all tags referenced by the rule exist in the database.
func (r *RuleService) Update(ctx context.Context, rule entity.Rule) error {
	existingRule, err := r.GetRule(ctx, rule.Name)
	if err != nil {
		return err
	}

	labelSrv := NewLabelService(r.dt)
	existingLabels, err := labelSrv.GetLabels(ctx)
	if err != nil {
		return err
	}

	if existingRule == nil {
		return NewErrRuleNotFound(rule.Name)
	}

	err = r.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		existingRelationships := []entity.Relationship{}
		// remove old relationships
		for _, label := range existingRule.Labels {
			existingRelationships = append(existingRelationships, entity.NewLabeRuleRelationship(label.ID, rule.Name))
		}

		if err := w.DeleteRelationships(ctx, existingRelationships); err != nil {
			return err
		}

		// update the rule
		if err := w.WriteRule(ctx, rule, true); err != nil {
			return err
		}

		// add new relationships
		relationships := []entity.Relationship{}
		for _, label := range rule.Labels {
			found := false
			labelID := 0
			for _, existingLabel := range existingLabels {
				if label.Equal(existingLabel) {
					found = true
					labelID = existingLabel.ID
					break
				}
			}
			if !found {
				id, err := w.WriteLabel(ctx, label)
				if err != nil {
					return err
				}
				labelID = id
			}

			relationships = append(relationships, entity.NewLabeRuleRelationship(labelID, rule.Name))
		}

		return w.WriteRelationships(ctx, relationships)
	})

	return err
}

// DeleteRule removes a rule from the database by its name.
func (r *RuleService) DeleteRule(ctx context.Context, name string) error {
	return r.dt.WriteTx(ctx, func(ctx context.Context, w *pg.Writer) error {
		return w.DeleteRule(ctx, name)
	})
}
