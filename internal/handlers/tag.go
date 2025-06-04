package handlers

import (
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/models"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func tagHandlers(r *gin.RouterGroup) {
	r.GET("/tags", func(c *gin.Context) {
		dt := MustFromContext(c)

		// get tags from tagSrv
		tagSrv := services.NewTagService(dt)
		tags, err := tagSrv.GetTags(c.Request.Context())
		if err != nil {
			zap.S().Errorw("failed to get tags", "error", err)
			c.JSON(500, err)
			return
		}

		ruleSrv := services.NewRuleService(dt)
		rules, err := ruleSrv.GetRules(c.Request.Context())
		if err != nil {
			zap.S().Errorw("failed to get rules", "error", err)
			c.JSON(500, err)
			return
		}

		c.JSON(200, models.NewTags(tags, rules))
	})
}
