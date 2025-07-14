package v1

import (
	"fmt"
	"net/http"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SummaryHandlers registers all summary-related HTTP handlers with the provided router group.
// This includes endpoints for retrieving aggregated statistics and summaries.
func SummaryHandlers(r *gin.RouterGroup) {
	r.GET("/summary", func(c *gin.Context) {
		now := time.Now()
		start, err := parseTimestamp(c.Query("startDate"), time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			zap.S().Warnw("failed to parse starting date timestamp. defaults to first day of the current month", "error", err, "url", c.Request.URL)
		}

		end, err := parseTimestamp(c.Query("endDate"), now)
		if err != nil {
			zap.S().Warnw("failed to parse ending date timestamp. defaults to now", "error", err, "url", c.Request.URL)
		}

		// Validate that start date is before end date
		if start.After(end) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "startDate must be before endDate",
				"startDate": start.Format(time.RFC3339),
				"endDate":   end.Format(time.RFC3339),
			})
			return
		}

		dt := dtContext.MustFromContext(c)
		srv := services.NewTransactionService(dt)
		transactions, err := srv.GetTransactions(c.Request.Context(), services.NewTransactionFilterWithOptions(
			services.WithStart(&start),
			services.WithEnd(&end),
		))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		total := float32(0)
		labelAmount := make(map[string]float32)

		// Get label service to lookup label information
		labelSrv := services.NewLabelService(dt)
		allLabels, err := labelSrv.GetLabels(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get labels"})
			return
		}

		// Create a map for quick label lookup by ID
		labelMap := make(map[int]entity.Label)
		for _, label := range allLabels {
			labelMap[label.ID] = label
		}

		for _, t := range transactions {
			total += t.Amount
			for labelID := range t.Labels {
				if label, exists := labelMap[labelID]; exists {
					labelKey := fmt.Sprintf("%s:%s", label.Key, label.Value)
					amount, ok := labelAmount[labelKey]
					if ok {
						amount += t.Amount
						labelAmount[labelKey] = amount
						continue
					}
					labelAmount[labelKey] = t.Amount
				}
			}
		}

		summary := outbound.Summary{
			StartingDate: start.Format("02/01/2006"),
			EndingDate:   end.Format("02/01/2006"),
			Items:        labelAmount,
			TotalAmmount: total,
		}

		c.JSON(http.StatusAccepted, summary)
	})
}
