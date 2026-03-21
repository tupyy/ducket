package services

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
)

type SummaryService struct {
	st *store.Store
}

func NewSummaryService(st *store.Store) *SummaryService {
	return &SummaryService{st: st}
}

func (s *SummaryService) Overview(ctx context.Context, filter string) (*store.SummaryOverview, error) {
	return s.st.GetSummaryOverview(ctx, filter)
}

func (s *SummaryService) ByTag(ctx context.Context, filter string) ([]store.TagSummary, error) {
	return s.st.GetTagSummary(ctx, filter)
}

func (s *SummaryService) BalanceTrend(ctx context.Context, filter string) ([]store.BalanceTrendPoint, error) {
	return s.st.GetBalanceTrend(ctx, filter)
}
