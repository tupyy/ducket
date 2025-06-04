package handlers

import (
	"net/http"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/models"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func tagHandlers(r *gin.RouterGroup) {
	validate := validator.New()
	validate.RegisterStructValidation(models.TagFormValidation, models.TagForm{})

	r.GET("/tags", func(c *gin.Context) {
		dt := MustFromContext(c)

		// get tags from tagSrv
		tagSrv := services.NewTagService(dt)
		tags, err := tagSrv.GetTags(c.Request.Context())
		if err != nil {
			zap.S().Errorw("failed to get tags", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}

		ruleSrv := services.NewRuleService(dt)
		rules, err := ruleSrv.GetRules(c.Request.Context())
		if err != nil {
			zap.S().Errorw("failed to get rules", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}

		c.JSON(http.StatusOK, models.NewTags(tags, rules))
	})

	r.POST("/tags", func(c *gin.Context) {
		var form models.TagForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dt := MustFromContext(c)
		tagSrv := services.NewTagService(dt)
		if err := tagSrv.Create(c.Request.Context(), form.Value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.NewTag(form.Value))
	})

	r.DELETE("/tags/:id", func(c *gin.Context) {
		id := c.Param("id")

		dt := MustFromContext(c)
		tagSrv := services.NewTagService(dt)
		if err := tagSrv.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	})
}
