package handlers

import (
	"net/http"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/inbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func summaryHandlers(r *gin.RouterGroup) {
	validate.RegisterStructValidation(inbound.TransactionFormValidation, inbound.TransactionForm{})

	r.GET("/summary", func(c *gin.Context) {
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

		total := float32(0)
		tagAmount := make(map[string]float32)
		for _, t := range transactions {
			total += t.Amount
			for tag := range t.Tags {
				amount, ok := tagAmount[tag]
				if ok {
					amount += t.Amount
					tagAmount[tag] = amount
					continue
				}
				tagAmount[tag] = t.Amount
			}
		}

		summary := outbound.Summary{
			StartingDate: start.Format("02/01/2006"),
			EndingDate:   end.Format("02/01/2006"),
			Items:        tagAmount,
			TotalAmmount: total,
		}

		c.JSON(http.StatusAccepted, summary)
	})
}
