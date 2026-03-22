package server

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/tupyy/ducket/api/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	srv *http.Server
}

func New(port int, ginMode string, handler v1.ServerInterface) *Server {
	gin.SetMode(ginMode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	api := engine.Group("/api/v1")
	v1.RegisterHandlers(api, handler)

	return &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", port),
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
