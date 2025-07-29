package v1

import (
	"net/http"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetLabels handles GET /labels requests to retrieve all available labels.
// It fetches all labels from the database through the label service and returns
// them as JSON. Returns HTTP 500 if there's an error retrieving the labels.
func (s *ServerImpl) GetLabels(c *gin.Context) {
	dt := dtContext.MustFromContext(c)

	labelSrv := services.NewLabelService(dt)
	labels, err := labelSrv.GetLabels(c.Request.Context())
	if err != nil {
		zap.S().Errorw("failed to get labels", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, outbound.NewLabels(labels))
}
