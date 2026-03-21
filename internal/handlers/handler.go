package handlers

import (
	"net/http"
	"strings"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
	srvErrors "git.tls.tupangiu.ro/cosmin/finante/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	txnSvc     *services.TransactionService
	ruleSvc    *services.RuleService
	summarySvc *services.SummaryService
}

func NewHandler(txnSvc *services.TransactionService, ruleSvc *services.RuleService, summarySvc *services.SummaryService) *Handler {
	return &Handler{txnSvc: txnSvc, ruleSvc: ruleSvc, summarySvc: summarySvc}
}

// -- Transactions --

var validTransactionSortFields = map[string]bool{
	"date":    true,
	"amount":  true,
	"kind":    true,
	"account": true,
}

func (h *Handler) ListTransactions(c *gin.Context, params v1.ListTransactionsParams) {
	filter, limit, offset := extractListParams(params.Filter, params.Limit, params.Offset)

	var sortParams []store.SortParam
	if params.Sort != nil {
		for _, s := range *params.Sort {
			parts := strings.SplitN(s, ":", 2)
			if len(parts) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort format, expected 'field:direction' (e.g., 'date:desc')"})
				return
			}
			field, direction := parts[0], parts[1]
			if !validTransactionSortFields[field] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort field: " + field})
				return
			}
			if direction != "asc" && direction != "desc" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort direction: " + direction + ", must be 'asc' or 'desc'"})
				return
			}
			sortParams = append(sortParams, store.SortParam{Field: field, Desc: direction == "desc"})
		}
	}

	tags := c.QueryArray("tags")

	txns, total, err := h.txnSvc.List(c.Request.Context(), filter, tags, sortParams, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	items := make([]v1.Transaction, 0, len(txns))
	for _, t := range txns {
		items = append(items, v1.NewTransactionFromEntity(t))
	}
	c.JSON(http.StatusOK, v1.TransactionList{Items: items, Total: total})
}

func (h *Handler) GetTransaction(c *gin.Context, id int64) {
	txn, err := h.txnSvc.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, v1.NewTransactionFromEntity(*txn))
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var req v1.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn := entity.NewTransaction(entity.TransactionKind(req.Kind), req.Account, req.Date, req.Amount, req.Content)
	txn.Info = req.Info
	txn.Recipient = req.Recipient

	created, err := h.txnSvc.Create(c.Request.Context(), *txn)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, v1.NewTransactionFromEntity(created))
}

func (h *Handler) UpdateTransaction(c *gin.Context, id int64) {
	var req v1.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn := entity.Transaction{
		ID:        id,
		Kind:      entity.TransactionKind(req.Kind),
		Account:   req.Account,
		Date:      req.Date,
		Amount:    req.Amount,
		Content:   req.Content,
		Info:      req.Info,
		Recipient: req.Recipient,
	}

	if err := h.txnSvc.Update(c.Request.Context(), txn); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteTransaction(c *gin.Context, id int64) {
	if err := h.txnSvc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// -- Rules --

func (h *Handler) ListRules(c *gin.Context, params v1.ListRulesParams) {
	filter, limit, offset := extractListParams(params.Filter, params.Limit, params.Offset)

	rules, err := h.ruleSvc.List(c.Request.Context(), filter, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	result := make([]v1.Rule, 0, len(rules))
	for _, r := range rules {
		result = append(result, v1.NewRuleFromEntity(r))
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetRule(c *gin.Context, id int) {
	rule, err := h.ruleSvc.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, v1.NewRuleFromEntity(*rule))
}

func (h *Handler) CreateRule(c *gin.Context) {
	var req v1.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.ruleSvc.Create(c.Request.Context(), entity.Rule{
		Name:   req.Name,
		Filter: req.Filter,
		Tags:   req.Tags,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, v1.NewRuleFromEntity(rule))
}

func (h *Handler) UpdateRule(c *gin.Context, id int) {
	var req v1.UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ruleSvc.Update(c.Request.Context(), entity.Rule{
		ID:     id,
		Name:   req.Name,
		Filter: req.Filter,
		Tags:   req.Tags,
	}); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteRule(c *gin.Context, id int) {
	if err := h.ruleSvc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// -- Helpers --

const (
	defaultLimit = 100
	maxLimit     = 1000
)

func extractListParams(filter *string, limit *int, offset *int) (string, int, int) {
	f := ""
	if filter != nil {
		f = *filter
	}
	l := defaultLimit
	if limit != nil && *limit > 0 {
		l = *limit
	}
	if l > maxLimit {
		l = maxLimit
	}
	o := 0
	if offset != nil && *offset > 0 {
		o = *offset
	}
	return f, l, o
}

func handleError(c *gin.Context, err error) {
	if srvErrors.IsResourceNotFoundError(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if srvErrors.IsDuplicateResourceError(err) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if srvErrors.IsValidationError(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	zap.S().Errorw("internal error", "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
