package services

import (
	"context"

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

func (r *RuleService) Create(ctx context.Context, rule entity.Rule) error {
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

func (r *RuleService) DeleteRule(ctx context.Context, rule entity.Rule) error {
	return r.dt.WriteTx(ctx, func(ctx context.Context, w pg.Writer) error {
		return w.DeleteRule(ctx, rule.ID)
	})
}
