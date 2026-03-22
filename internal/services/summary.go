package services

import (
	"context"

	"github.com/tupyy/ducket/internal/store"
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
