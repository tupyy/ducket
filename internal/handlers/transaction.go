package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/inbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	queryDateFormat = "02/01/2006"
)

func transactionHandlers(r *gin.RouterGroup) {
	validate.RegisterStructValidation(inbound.TransactionFormValidation, inbound.TransactionForm{})

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

	r.POST("/transactions", func(c *gin.Context) {
		var form inbound.TransactionForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tEntity, err := form.Entity()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dt := MustFromContext(c)
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

		if err := validate.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tEntity, err := form.Entity()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tEntity.ID = id

		dt := MustFromContext(c)
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

		dt := MustFromContext(c)
		tSrv := services.NewTransactionService(dt)
		if err := tSrv.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
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
