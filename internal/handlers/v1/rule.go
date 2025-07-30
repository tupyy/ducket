package v1

import (
	"net/http"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"

	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GetRules handles GET /rules requests to retrieve all available rules.
// It fetches all rules from the database through the rule service and returns
// them as JSON. Returns HTTP 500 if there's an error retrieving the rules.
func (s *ServerImpl) GetRules(c *gin.Context) {
	dt := dtContext.MustFromContext(c)

	ruleSrv := services.NewRuleService(dt)
	rules, err := ruleSrv.GetRules(c.Request.Context())
	if err != nil {
		zap.S().Errorw("failed to get rules", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, v1.NewRules(rules))
}

// CreateRule handles POST /rules requests to create a new rule.
// It validates the request body as a RuleForm, checks business validation rules,
// then creates the rule through the rule service. Returns HTTP 400 for validation
// errors or HTTP 201 on successful creation.
func (s *ServerImpl) CreateRule(c *gin.Context) {
	var form v1.RuleForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	validator.RegisterStructValidation(v1.RuleFormValidation, v1.RuleForm{})

	if err := validator.Struct(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dt := dtContext.MustFromContext(c)
	ruleSrv := services.NewRuleService(dt)
	if err := ruleSrv.Create(c.Request.Context(), form.Entity()); err != nil {
		zap.S().Errorw("failed to create rule", "error", err.Error(), "form", form)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, form)
}

// UpdateRule handles PUT /rules/{id} requests to update an existing rule or create one if it doesn't exist.
// It validates the request body as an UpdateRuleForm, then attempts to update the rule.
// If the rule doesn't exist, it creates a new one. Returns HTTP 400 for validation errors,
// HTTP 500 for server errors, HTTP 201 for creation, or HTTP 200 for successful update.
func (s *ServerImpl) UpdateRule(c *gin.Context, id string) {
	var form v1.RuleForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	validator.RegisterStructValidation(v1.RuleFormValidation, v1.RuleForm{})

	if err := validator.Struct(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dt := dtContext.MustFromContext(c)
	ruleSrv := services.NewRuleService(dt)

	// Convert form to entity, but override the name with the ID from the URL
	ruleEntity := form.Entity()
	ruleToCreate := entity.NewRule(id, ruleEntity.Pattern, ruleEntity.Labels...)
	err := ruleSrv.Update(c.Request.Context(), ruleToCreate)
	if err != nil {
		switch err.(type) {
		case *services.ErrResourceNotFound:
			if err := ruleSrv.Create(c.Request.Context(), ruleToCreate); err != nil {
				zap.S().Errorw("failed to create rule", "error", err.Error(), "form", form)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, ruleToCreate)
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, ruleToCreate)
}

// DeleteRule handles DELETE /rules/{id} requests to remove a rule by its ID.
// It validates that the rule name is not longer than 20 characters, then deletes
// the rule through the rule service. Returns HTTP 400 for validation errors or
// HTTP 204 on successful deletion.
func (s *ServerImpl) DeleteRule(c *gin.Context, id string) {
	dt := dtContext.MustFromContext(c)
	ruleSrv := services.NewRuleService(dt)
	if err := ruleSrv.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// ProcessRules handles POST /rules/process requests to apply all rules to all transactions.
// It runs the rule applier service to automatically label transactions based on all
// configured rules. This is a bulk operation that processes all transactions.
// Returns HTTP 500 if the operation fails or HTTP 200 on success.
func (s *ServerImpl) ProcessRules(c *gin.Context) {
	dt := dtContext.MustFromContext(c)

	ruleApplierService := services.NewRuleApplierService(dt)

	zap.S().Info("Starting rule application to all transactions")

	err := ruleApplierService.ApplyAllRulesToAllTransactions(c.Request.Context())
	if err != nil {
		zap.S().Errorw("Failed to apply rules to all transactions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to apply rules to transactions",
			"success": false,
		})
		return
	}

	zap.S().Info("Successfully applied rules to all transactions")

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully applied all rules to all transactions",
		"success": true,
	})
}

// ProcessRule handles POST /rules/{id}/process requests to apply a specific rule to all transactions.
// It validates the rule name (max 20 characters), retrieves the rule, and applies it to all
// transactions in the system. Returns HTTP 400 for validation errors, HTTP 404 if the rule
// doesn't exist, HTTP 500 for server errors, or HTTP 200 on successful processing.
func (s *ServerImpl) ProcessRule(c *gin.Context, id string) {
	ruleName := id

	dt := dtContext.MustFromContext(c)

	// Get the specific rule
	ruleService := services.NewRuleService(dt)
	rule, err := ruleService.GetRule(c.Request.Context(), ruleName)
	if err != nil {
		zap.S().Errorw("Failed to get rule", "name", ruleName, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get rule",
			"success": false,
		})
		return
	}

	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Rule not found",
			"success": false,
		})
		return
	}

	// Get all transactions
	transactionService := services.NewTransactionService(dt)
	transactions, err := transactionService.GetTransactions(c.Request.Context(), &services.TransactionFilter{})
	if err != nil {
		zap.S().Errorw("Failed to get transactions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get transactions",
			"success": false,
		})
		return
	}

	zap.S().Infow("Starting rule application to all transactions",
		"rule", ruleName,
		"total_transactions", len(transactions),
	)

	// Apply the specific rule to all transactions
	ruleApplierService := services.NewRuleApplierService(dt)
	matchCount := 0

	for _, transaction := range transactions {
		// Check if this rule matches the transaction
		matched, err := ruleApplierService.MatchesRule(transaction.RawContent, rule.Pattern)
		if err != nil {
			zap.S().Warnw("Failed to check rule pattern",
				"rule", ruleName,
				"pattern", rule.Pattern,
				"error", err,
			)
			continue
		}

		if matched {
			// Apply all labels from this rule to the transaction
			updatedTransaction := transaction
			if updatedTransaction.Labels == nil {
				updatedTransaction.Labels = make([]entity.LabelAssociation, 0)
			}

			for _, label := range rule.Labels {
				updatedTransaction.Labels = append(updatedTransaction.Labels, entity.LabelAssociation{Label: label, RuleID: &ruleName})
			}

			// Update the transaction in the database
			_, err = transactionService.Update(c.Request.Context(), updatedTransaction)
			if err != nil {
				zap.S().Errorw("Failed to update transaction with rule",
					"transaction_id", transaction.ID,
					"rule", ruleName,
					"error", err,
				)
				continue
			}

			matchCount++
			zap.S().Debugw("Applied rule to transaction",
				"rule", ruleName,
				"transaction_id", transaction.ID,
				"transaction_content", transaction.RawContent,
				"labels", rule.Labels,
			)
		}
	}

	zap.S().Infow("Successfully applied rule to transactions",
		"rule", ruleName,
		"matches_found", matchCount,
		"total_transactions", len(transactions),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":            "Successfully applied rule to all transactions",
		"success":            true,
		"rule":               ruleName,
		"matches_found":      matchCount,
		"total_transactions": len(transactions),
	})
}
