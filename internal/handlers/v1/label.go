package v1

import (
	"net/http"

	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/outbound"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
