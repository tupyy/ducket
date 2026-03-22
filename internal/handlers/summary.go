package handlers

import (
	"net/http"

	v1 "github.com/tupyy/ducket/api/v1"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSummaryOverview(c *gin.Context, params v1.GetSummaryOverviewParams) {
	filter := ""
	if params.Filter != nil {
		filter = *params.Filter
	}

	overview, err := h.summarySvc.Overview(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, v1.SummaryOverview{
		TotalTransactions: overview.TotalTransactions,
		TotalDebit:        overview.TotalDebit,
		TotalCredit:       overview.TotalCredit,
		Balance:           overview.Balance,
		UniqueAccounts:    overview.UniqueAccounts,
		UniqueTags:        overview.UniqueTags,
	})
}

func (h *Handler) GetSummaryByTag(c *gin.Context, params v1.GetSummaryByTagParams) {
	filter := ""
	if params.Filter != nil {
		filter = *params.Filter
	}

	tags, err := h.summarySvc.ByTag(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}

	result := make([]v1.TagSummary, 0, len(tags))
	for _, t := range tags {
		result = append(result, v1.TagSummary{
			Tag:        t.Tag,
			TotalDebit: t.TotalDebit,
			TotalCredit: t.TotalCredit,
			Count:      t.Count,
		})
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetSummaryBalanceTrend(c *gin.Context, params v1.GetSummaryBalanceTrendParams) {
	filter := ""
	if params.Filter != nil {
		filter = *params.Filter
	}

	trend, err := h.summarySvc.BalanceTrend(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}

	result := make([]v1.BalanceTrendPoint, 0, len(trend))
	for _, p := range trend {
		result = append(result, v1.BalanceTrendPoint{
			Month:   p.Month,
			Debit:   p.Debit,
			Credit:  p.Credit,
			Balance: p.Balance,
		})
	}
	c.JSON(http.StatusOK, result)
}
