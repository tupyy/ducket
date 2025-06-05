package services

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

type RuleService struct {
	dt *pg.Datastore
}

func NewRuleService(dt *pg.Datastore) *RuleService {
	return &RuleService{dt: dt}
}

func (r *RuleService) GetRules(ctx context.Context) ([]entity.Rule, error) {
	return r.dt.QueryRules(ctx, pg.RuleFilter{}, &pg.QueryRuleOptions{})
}

func (r *RuleService) GetRule(ctx context.Context, name string) (*entity.Rule, error) {
	rules, err := r.dt.QueryRules(ctx, pg.RuleFilter{Name: &name}, &pg.QueryRuleOptions{})
	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return nil, nil
	}

	return &rules[0], nil
}

func (r *RuleService) Create(ctx context.Context, rule entity.Rule) error {
	existingRule, err := r.GetRule(ctx, rule.Name)
	if err != nil {
		return err
	}

	if existingRule != nil {
		return fmt.Errorf("rule %s already exists", rule.Name)
	}

	tagSrv := NewTagService(r.dt)
	existingTags, err := tagSrv.GetTags(ctx)
	if err != nil {
		return err
	}

	return r.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		for _, tag := range rule.Tags {
			found := false
			for _, t := range existingTags {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				if err := w.WriteTag(ctx, tag); err != nil {
					return err
				}
			}
		}
		// write rule
		return w.WriteRule(ctx, rule, false)
	})
}

func (r *RuleService) UpdateOrCreate(ctx context.Context, rule entity.Rule) (bool, error) {
	existingRule, err := r.GetRule(ctx, rule.Name)
	if err != nil {
		return false, err
	}

	tagSrv := NewTagService(r.dt)
	existingTags, err := tagSrv.GetTags(ctx)
	if err != nil {
		return false, err
	}

	update := existingRule != nil
	err = r.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		for _, tag := range rule.Tags {
			found := false
			for _, t := range existingTags {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				if err := w.WriteTag(ctx, tag); err != nil {
					return err
				}
			}
		}
		return w.WriteRule(ctx, rule, update)
	})

	return update, err
}

func (r *RuleService) DeleteRule(ctx context.Context, name string) error {
	return r.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		return w.DeleteRule(ctx, name)
	})
}
