package v1

import (
	"net/http"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/inbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func TagHandlers(r *gin.RouterGroup) {
	r.GET("/tags", func(c *gin.Context) {
		dt := dtContext.MustFromContext(c)

		tagSrv := services.NewTagService(dt)
		tags, err := tagSrv.GetTags(c.Request.Context())
		if err != nil {
			zap.S().Errorw("failed to get tags", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}

		c.JSON(http.StatusOK, outbound.NewTags(tags))
	})

	r.POST("/tags", func(c *gin.Context) {
		var form inbound.TagForm
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
		tagSrv := services.NewTagService(dt)
		if err := tagSrv.Create(c.Request.Context(), form.Value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, outbound.NewTag(form.Value))
	})

	r.DELETE("/tags/:id", func(c *gin.Context) {
		id := c.Param("id")

		dt := dtContext.MustFromContext(c)
		tagSrv := services.NewTagService(dt)
		if err := tagSrv.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	})
}
