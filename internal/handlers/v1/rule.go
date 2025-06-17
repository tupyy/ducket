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

func RulesHandlers(r *gin.RouterGroup, validator *validator.Validate) {
	validator.RegisterStructValidation(inbound.RuleFormValidation, inbound.RuleForm{})
	validator.RegisterStructValidation(inbound.UpdateRuleFormValidation, inbound.UpdateRuleForm{})

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

		if err := validator.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dt := dtContext.MustFromContext(c)
		ruleSrv := services.NewRuleService(dt)

		ruleToCreate := entity.NewRule(name, form.Pattern, form.Tags...)
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

}
