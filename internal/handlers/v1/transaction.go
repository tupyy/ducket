package v1

import (
	"fmt"
	"net/http"
	"time"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"

	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	queryDateFormat = "02/01/2006"
)

// GetTransaction handles GET /transactions/{id} requests to retrieve a specific transaction by ID.
// It fetches the transaction from the database through the transaction service.
// Returns HTTP 404 if the transaction is not found, HTTP 500 for server errors,
// or HTTP 201 with the transaction data on success.
func (s *ServerImpl) GetTransaction(c *gin.Context, id int64) {
	dt := dtContext.MustFromContext(c)

	// Add the label to the transaction - this will handle transaction not found errors
	tSrv := services.NewTransactionService(dt)
	transaction, err := tSrv.GetTransactionById(c.Request.Context(), int(id))
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, v1.NewTransaction(*transaction))
}

// GetTransactions handles GET /transactions requests to retrieve a filtered list of transactions.
// It supports optional query parameters for date filtering (startDate and endDate in DD/MM/YYYY format).
// Defaults to the current month if no dates are provided. Returns HTTP 400 for invalid parameters,
// HTTP 500 for server errors, or HTTP 200 with the transaction list on success.
func (s *ServerImpl) GetTransactions(c *gin.Context, params v1.GetTransactionsParams) {
	startDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	if params.StartDate != nil {
		start, err := parseTimestamp(*params.StartDate)
		if err != nil {
			zap.S().Warnw("failed to parse starting date timestamp. defaults to first day of the current month", "error", err, "url", c.Request.URL)
		} else {
			startDate = start
		}
	}

	endDate := time.Date(time.Now().Year(), time.Now().Month(), 31, 0, 0, 0, 0, time.UTC)
	if params.EndDate != nil {
		end, err := parseTimestamp(*params.EndDate)
		if err != nil {
			zap.S().Warnw("failed to parse ending date timestamp. defaults to end of current month", "error", err, "url", c.Request.URL)
		} else {
			endDate = end
		}
	}

	// Validate that start date is before end date
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "startDate must be before endDate",
			"startDate": startDate.Format(time.RFC3339),
			"endDate":   endDate.Format(time.RFC3339),
		})
		return
	}

	dt := dtContext.MustFromContext(c)
	srv := services.NewTransactionService(dt)
	transactions, err := srv.GetTransactions(c.Request.Context(), services.NewTransactionFilterWithOptions(
		services.WithStart(&startDate),
		services.WithEnd(&endDate),
	))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t := v1.NewTransactions(len(transactions), startDate, endDate)
	for _, transaction := range transactions {
		*t.Items = append(*t.Items, v1.NewTransaction(transaction))
	}

	c.JSON(http.StatusOK, t)
}

// CreateTransaction handles POST /transactions requests to create a new transaction.
// It validates the request body as a CreateTransactionForm, checks for duplicate transactions
// by hash, then creates the transaction through the transaction service. Returns HTTP 400
// for validation errors or duplicate transactions, HTTP 500 for server errors, or HTTP 201
// with the created transaction on success.
func (s *ServerImpl) CreateTransaction(c *gin.Context) {
	var form v1.TransactionForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	validator.RegisterStructValidation(v1.CreateTransactionFormValidation, v1.TransactionForm{})
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

	t, err := tSrv.Create(c.Request.Context(), tEntity)
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceExistsAlready:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, v1.NewTransaction(t))
}

// UpdateTransaction handles PUT /transactions/{id} requests to update an existing transaction or create one if it doesn't exist.
// It validates the request body as a CreateTransactionForm, then attempts to update the transaction.
// If the transaction doesn't exist, it creates a new one with the provided ID. Returns HTTP 400
// for validation errors, HTTP 500 for server errors, HTTP 201 for creation, or HTTP 200 for successful update.
func (s *ServerImpl) UpdateTransaction(c *gin.Context, id int64) {
	var form v1.TransactionForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	validator.RegisterStructValidation(v1.CreateTransactionFormValidation, v1.TransactionForm{})

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

	updatedTransaction, err := tSrv.Update(c.Request.Context(), tEntity)
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			newTransaction, err := tSrv.Create(c.Request.Context(), tEntity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, v1.NewTransaction(newTransaction))
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, v1.NewTransaction(updatedTransaction))
}

// DeleteTransaction handles DELETE /transactions/{id} requests to remove a transaction by its ID.
// It permanently deletes the transaction from the database through the transaction service.
// Returns HTTP 400 if the deletion fails or HTTP 204 on successful deletion.
func (s *ServerImpl) DeleteTransaction(c *gin.Context, id int64) {
	dt := dtContext.MustFromContext(c)
	tSrv := services.NewTransactionService(dt)
	if err := tSrv.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// PatchTransactionInfo handles PATCH /transactions/{id} requests to update the info field of a transaction.
// It validates the request body as an UpdateTransactionInfoForm, then updates only the info field
// through the transaction service. Returns HTTP 400 for validation errors, HTTP 404 if the transaction
// is not found, HTTP 500 for server errors, or HTTP 200 with the updated transaction on success.
func (s *ServerImpl) PatchTransactionInfo(c *gin.Context, id int64) {
	var form v1.TransactionInfoForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	validator.RegisterStructValidation(v1.UpdateTransactionInfoFormValidation, v1.TransactionInfoForm{})
	if err := validator.Struct(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dt := dtContext.MustFromContext(c)
	tSrv := services.NewTransactionService(dt)

	info := ""
	if form.Info != nil {
		info = *form.Info
	}
	updatedTransaction, err := tSrv.UpdateInfo(c.Request.Context(), id, info)
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, v1.NewTransaction(updatedTransaction))
}

// Transaction Labels endpoints

// GET /api/v1/transactions/{id}/labels - Get all labels for a transaction
func (s *ServerImpl) GetTransactionLabels(c *gin.Context, id int64) {
	dt := dtContext.MustFromContext(c)
	tSrv := services.NewTransactionService(dt)

	// Get labels for the transaction - this will handle transaction not found errors
	labels, err := tSrv.Labels(c.Request.Context(), int(id))
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, v1.NewLabels(labels))
}

// AddTransactionLabel handles POST /transactions/{id}/labels requests to add a label to a transaction.
// It validates the request body as a LabelForm, retrieves the transaction, adds the label to it,
// and updates the transaction. Returns HTTP 400 for validation errors, HTTP 404 if the transaction
// is not found, HTTP 500 for server errors, or HTTP 200 with the updated transaction on success.
func (s *ServerImpl) AddTransactionLabel(c *gin.Context, id int64) {
	var form v1.LabelForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	if err := validator.Struct(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dt := dtContext.MustFromContext(c)

	// Add the label to the transaction - this will handle transaction not found errors
	tSrv := services.NewTransactionService(dt)
	transaction, err := tSrv.GetTransactionById(c.Request.Context(), int(id))
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	transaction.Labels = append(transaction.Labels, entity.LabelAssociation{
		Label: form.Entity(),
	})

	updated, err := tSrv.Update(c.Request.Context(), *transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, v1.NewTransaction(updated))
}

// DELETE /api/v1/transactions/{id}/labels - Remove all labels from a transaction
func (s *ServerImpl) RemoveTransactionLabels(c *gin.Context, id int64) {
	dt := dtContext.MustFromContext(c)

	// Check if transaction exists by trying to get its labels
	tSrv := services.NewTransactionService(dt)
	// Add the label to the transaction - this will handle transaction not found errors
	transaction, err := tSrv.GetTransactionById(c.Request.Context(), int(id))
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	transaction.Labels = []entity.LabelAssociation{}
	updated, err := tSrv.Update(c.Request.Context(), *transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, v1.NewTransaction(updated))
}

// RemoveTransactionLabel handles DELETE /transactions/{id}/labels/{labelId} requests to remove a specific label from a transaction.
// It retrieves the transaction, finds and removes the label with the specified ID, then updates the transaction.
// Returns HTTP 404 if the transaction or label is not found, HTTP 500 for server errors,
// or HTTP 200 with the updated transaction on success.
func (s *ServerImpl) RemoveTransactionLabel(c *gin.Context, id int64, labelId int) {
	dt := dtContext.MustFromContext(c)

	// Check if transaction exists by trying to get its labels
	tSrv := services.NewTransactionService(dt)
	// Add the label to the transaction - this will handle transaction not found errors
	transaction, err := tSrv.GetTransactionById(c.Request.Context(), int(id))
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	found := false
	for _, a := range transaction.Labels {
		if a.Label.ID == labelId {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("label %d not found on transaction %d", labelId, id)})
		return
	}

	newLabels := []entity.LabelAssociation{}
	for _, a := range transaction.Labels {
		if a.Label.ID != labelId {
			newLabels = append(newLabels, a)
		}
	}

	transaction.Labels = newLabels
	updated, err := tSrv.Update(c.Request.Context(), *transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, v1.NewTransaction(updated))
}

// parseTimestamp parses a timestamp string (milliseconds since epoch) and returns the corresponding time.
// Used for parsing query parameters that represent timestamps from the frontend.
func parseTimestamp(timestamp int64) (time.Time, error) {
	// Convert milliseconds to seconds for time.Unix
	return time.Unix(timestamp/1000, (timestamp%1000)*1000000), nil
}
