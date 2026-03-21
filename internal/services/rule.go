package services

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
)

type RuleService struct {
	st *store.Store
}

func NewRuleService(st *store.Store) *RuleService {
	return &RuleService{st: st}
}

func (r *RuleService) List(ctx context.Context, filter string, limit, offset int) ([]entity.Rule, error) {
	return r.st.ListRules(ctx, filter, limit, offset)
}

func (r *RuleService) Get(ctx context.Context, id int) (*entity.Rule, error) {
	return r.st.GetRule(ctx, id)
}

func (r *RuleService) Create(ctx context.Context, rule entity.Rule) (entity.Rule, error) {
	if _, err := store.ParseTransactionFilter(rule.Filter); err != nil {
		return rule, fmt.Errorf("invalid filter: %w", err)
	}

	var id int
	if err := r.st.WithTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.st.CreateRule(ctx, rule)
		return err
	}); err != nil {
		return rule, err
	}

	rule.ID = id
	return rule, nil
}

func (r *RuleService) Update(ctx context.Context, rule entity.Rule) error {
	existing, err := r.st.GetRule(ctx, rule.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return NewErrRuleNotFound(rule.ID)
	}

	if _, err := store.ParseTransactionFilter(rule.Filter); err != nil {
		return fmt.Errorf("invalid filter: %w", err)
	}

	return r.st.WithTx(ctx, func(ctx context.Context) error {
		return r.st.UpdateRule(ctx, rule)
	})
}

func (r *RuleService) Delete(ctx context.Context, id int) error {
	return r.st.WithTx(ctx, func(ctx context.Context) error {
		return r.st.DeleteRule(ctx, id)
	})
}
