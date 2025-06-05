package handlers

import (
	"net/http"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/inbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func rulesHandlers(r *gin.RouterGroup) {
	validate.RegisterStructValidation(inbound.TagFormValidation, inbound.RuleForm{})

	r.GET("/rules", func(c *gin.Context) {
		dt := MustFromContext(c)

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

		if err := validate.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dt := MustFromContext(c)
		ruleSrv := services.NewRuleService(dt)
		if err := ruleSrv.Create(c.Request.Context(), inbound.FormToEntity(form)); err != nil {
			zap.S().Errorw("failed to create rule", "error", err.Error(), "form", form)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{})
	})

}
