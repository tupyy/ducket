package v1

import (
	"net/http"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/inbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// RulesHandlers registers all rule-related HTTP handlers with the provided router group.
// This includes endpoints for CRUD operations on rules.
func RulesHandlers(r *gin.RouterGroup) {
	r.GET("/rules", func(c *gin.Context) {
		dt := dtContext.MustFromContext(c)

		ruleSrv := services.NewRuleService(dt)
		rules, err := ruleSrv.GetRules(c.Request.Context())
		if err != nil {
			zap.S().Errorw("failed to get rules", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}

		c.JSON(http.StatusOK, outbound.NewRules(rules))
	})

	r.POST("/rules", func(c *gin.Context) {
		var form inbound.RuleForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validator := validator.New()
		validator.RegisterStructValidation(inbound.RuleFormValidation, inbound.RuleForm{})

		if err := validator.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dt := dtContext.MustFromContext(c)
		ruleSrv := services.NewRuleService(dt)
		if err := ruleSrv.Create(c.Request.Context(), inbound.FormToEntity(form)); err != nil {
			zap.S().Errorw("failed to create rule", "error", err.Error(), "form", form)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, form)
	})

	r.PUT("/rules/:id", func(c *gin.Context) {
		name := c.Param("id")

		if len(name) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name must be have less than 20 chars"})
			return
		}

		var form inbound.UpdateRuleForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validator := validator.New()
		validator.RegisterStructValidation(inbound.UpdateRuleFormValidation, inbound.UpdateRuleForm{})

		if err := validator.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dt := dtContext.MustFromContext(c)
		ruleSrv := services.NewRuleService(dt)

		labels := make([]entity.Label, 0, len(form.Labels))
		for key, value := range form.Labels {
			labels = append(labels, entity.Label{
				Key:   key,
				Value: value,
			})
		}
		ruleToCreate := entity.NewRule(name, form.Pattern, labels...)
		updated, err := ruleSrv.UpdateOrCreate(c.Request.Context(), ruleToCreate)
		if err != nil {
			zap.S().Errorw("failed to create rule", "error", err.Error(), "form", form)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status := http.StatusCreated
		if updated {
			status = http.StatusOK
		}

		c.JSON(status, ruleToCreate)
	})

	r.DELETE("/rules/:id", func(c *gin.Context) {
		name := c.Param("id")

		if len(name) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name must be have less than 20 chars"})
			return
		}

		dt := dtContext.MustFromContext(c)
		ruleSrv := services.NewRuleService(dt)
		if err := ruleSrv.DeleteRule(c.Request.Context(), name); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	})

	r.POST("/rules/apply", func(c *gin.Context) {
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
	})

	r.POST("/rules/:id/sync", func(c *gin.Context) {
		ruleName := c.Param("id")

		if len(ruleName) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rule name is required"})
			return
		}

		if len(ruleName) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rule name must have less than 20 chars"})
			return
		}

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
					updatedTransaction.Labels = make(map[int]string)
				}

				for _, label := range rule.Labels {
					// TODO: Need to resolve label.ID properly
					// For now, using a placeholder ID since label.ID might not be set
					updatedTransaction.Labels[label.ID] = rule.Name
				}

				// Update the transaction in the database
				_, err = transactionService.CreateOrUpdate(c.Request.Context(), updatedTransaction)
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
	})

}
