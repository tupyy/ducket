package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/inbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	queryDateFormat = "02/01/2006"
)

// TransactionHandlers registers all transaction-related HTTP handlers with the provided router group.
// This includes endpoints for CRUD operations on transactions.
func TransactionHandlers(r *gin.RouterGroup) {
	r.GET("/transactions", func(c *gin.Context) {
		now := time.Now()
		start, err := parseTimestamp(c.Query("startDate"), time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			zap.S().Warnw("failed to parse starting date timestamp. defaults to first day of the current month", "error", err, "url", c.Request.URL)
		}

		end, err := parseTimestamp(c.Query("endDate"), time.Date(now.Year(), now.Month(), 31, 0, 0, 0, 0, time.UTC))
		if err != nil {
			zap.S().Warnw("failed to parse ending date timestamp. defaults to end of current month", "error", err, "url", c.Request.URL)
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

		t := outbound.NewTransactions(len(transactions), start, end)
		for _, transaction := range transactions {
			t.Items = append(t.Items, outbound.FromEntity(transaction))
		}

		c.JSON(http.StatusOK, t)
	})

	r.POST("/transactions", func(c *gin.Context) {
		var form inbound.TransactionForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validator := validator.New()
		validator.RegisterStructValidation(inbound.TransactionFormValidation, inbound.TransactionForm{})
		if err := validator.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tEntity, err := form.Entity()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dt := dtContext.MustFromContext(c)
		tSrv := services.NewTransactionService(dt)

		existingTransaction, err := tSrv.GetTransaction(c.Request.Context(), tEntity.Hash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if existingTransaction != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("transaction %d already exists", existingTransaction.ID)})
			return
		}

		t, err := tSrv.CreateOrUpdate(c.Request.Context(), tEntity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, outbound.FromEntity(t))
	})

	r.PUT("/transactions/:id", func(c *gin.Context) {
		idParam := c.Param("id")

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an int"})
			return
		}

		var form inbound.TransactionForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validator := validator.New()
		validator.RegisterStructValidation(inbound.TransactionFormValidation, inbound.TransactionForm{})

		if err := validator.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tEntity, err := form.Entity()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tEntity.ID = id

		dt := dtContext.MustFromContext(c)
		tSrv := services.NewTransactionService(dt)
		t, err := tSrv.CreateOrUpdate(c.Request.Context(), tEntity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, outbound.FromEntity(t))
	})

	r.DELETE("/transactions/:id", func(c *gin.Context) {
		idParam := c.Param("id")

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an int"})
			return
		}

		dt := dtContext.MustFromContext(c)
		tSrv := services.NewTransactionService(dt)
		if err := tSrv.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	})

}

// parseTimestamp parses a timestamp string (milliseconds since epoch) and returns the corresponding time.
// Used for parsing query parameters that represent timestamps from the frontend.
func parseTimestamp(sTimestamp string, defaultTime time.Time) (time.Time, error) {
	if sTimestamp == "" {
		return defaultTime, nil
	}
	timestamp, err := strconv.ParseInt(sTimestamp, 10, 64)
	if err != nil {
		return defaultTime, err
	}
	// Convert milliseconds to seconds for time.Unix
	return time.Unix(timestamp/1000, (timestamp%1000)*1000000), nil
}
