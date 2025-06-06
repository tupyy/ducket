package handlers

import (
	"net/http"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	queryDateFormat = "02/01/2006"
)

func transactionHandlers(r *gin.RouterGroup) {
	r.GET("/transactions", func(c *gin.Context) {
		now := time.Now()
		start, err := parseTime(c.Query("start"), time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			zap.S().Warnw("failed to parse starting date. defaults to first day of the current month", "error", err, "url", c.Request.URL)
		}

		end, err := parseTime(c.Query("end"), now)
		if err != nil {
			zap.S().Warnw("failed to parse ending date. defaults to now", "error", err, "url", c.Request.URL)
		}

		dt := MustFromContext(c)
		srv := services.NewTransactionService(dt)
		transactions, err := srv.GetTransactions(c.Request.Context(), services.NewTransactionFilterWithOptions(
			services.WithStart(&start),
			services.WithEnd(&end),
		))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		t := outbound.NewTransactions(len(transactions), start, end)
		for _, transaction := range transactions {
			t.Items = append(t.Items, outbound.FromEntity(transaction))
		}

		c.JSON(http.StatusAccepted, t)
	})
}

func parseTime(sTime string, defaultTime time.Time) (time.Time, error) {
	if sTime == "" {
		return defaultTime, nil
	}
	startTime, err := time.Parse(queryDateFormat, sTime)
	if err != nil {
		return defaultTime, err
	}
	return startTime, nil
}
