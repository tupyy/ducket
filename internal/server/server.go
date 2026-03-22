package server

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	v1 "github.com/tupyy/ducket/api/v1"
	"github.com/tupyy/ducket/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	srv *http.Server
}

func New(cfg *config.Config, handler v1.ServerInterface) *Server {
	gin.SetMode(cfg.GinMode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	api := engine.Group("/api/v1")
	v1.RegisterHandlers(api, handler)

	if cfg.StaticsFolder != "" {
		engine.Static("/static", path.Join(cfg.StaticsFolder, "static"))
		engine.Static("/images", path.Join(cfg.StaticsFolder, "images"))
		engine.StaticFile("/", path.Join(cfg.StaticsFolder, "index.html"))

		engine.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
				return
			}
			c.File(path.Join(cfg.StaticsFolder, "index.html"))
		})
	}

	return &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.ServerPort),
			Handler: engine,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		zap.S().Errorw("server shutdown error", "error", err)
	}
}
