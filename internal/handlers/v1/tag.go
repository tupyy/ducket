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

// LabelHandlers registers all label-related HTTP handlers with the provided router group.
// This includes endpoints for CRUD operations on labels.
func LabelHandlers(r *gin.RouterGroup) {
	r.GET("/labels", func(c *gin.Context) {
		dt := dtContext.MustFromContext(c)

		labelSrv := services.NewLabelService(dt)
		labels, err := labelSrv.GetLabels(c.Request.Context())
		if err != nil {
			zap.S().Errorw("failed to get labels", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, outbound.NewLabels(labels))
	})

	r.POST("/labels", func(c *gin.Context) {
		var form inbound.LabelForm
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
		labelSrv := services.NewLabelService(dt)
		if err := labelSrv.CreateLabel(c.Request.Context(), form.ToEntity()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, outbound.NewLabel(form.Key, form.Value))
	})

	r.DELETE("/labels/:id", func(c *gin.Context) {
		// TODO: Implement label deletion by ID
		// This would require parsing the ID and calling labelSrv.Delete
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Label deletion not yet implemented"})
	})
}
