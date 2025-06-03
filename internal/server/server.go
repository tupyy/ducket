package server

import (
	"context"
	"net/http"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/config"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateRunnableServer(ctx context.Context, config *config.Config) *gin.Engine {
	engine := gin.New()
	gin.SetMode(config.GinMode)

	privateRouter := engine.Group("/")

	privateRouter.Use(
		ginzap.Ginzap(zap.S().Desugar(), time.RFC3339, true),
		ginzap.RecoveryWithZap(zap.S().Desugar(), true),
	)

	engine.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	return engine
}
