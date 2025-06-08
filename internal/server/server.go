package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/server/middlewares"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RunnableServerConfig struct {
	Datastore          *pg.Datastore
	GraceTimeout       time.Duration
	Port               int
	RegisterHandlersFn func(router *gin.RouterGroup)
	CloseCb            func() error
	GinMode            string
}

type runnableServer struct {
	srv         *http.Server
	cfg         *RunnableServerConfig
	engine      *gin.Engine
	closePostCb func() error
}

func NewRunnableServer(cfg *RunnableServerConfig) *runnableServer {
	gin.SetMode(cfg.GinMode)
	engine := gin.New()

	router := engine.Group("/api/v1/")
	router.Use(
		middlewares.Logger(),
		middlewares.DatastoreMiddleware(cfg.Datastore),
		ginzap.RecoveryWithZap(zap.S().Desugar(), true),
	)

	// register handlers
	cfg.RegisterHandlersFn(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: engine,
	}

	return &runnableServer{srv: srv, cfg: cfg, closePostCb: cfg.CloseCb}
}

func (r *runnableServer) Run(ctx context.Context) {
	go func() {
		if err := r.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalw("server closed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Infow("shutdown server", "grace timeout", fmt.Sprintf("%s", r.cfg.GraceTimeout))

	newCtx, cancel := context.WithTimeout(ctx, r.cfg.GraceTimeout)
	defer cancel()
	go func() {
		if err := r.srv.Shutdown(newCtx); err != nil {
			zap.S().Errorw("server shutdown", "error", err)
		}
	}()

	<-newCtx.Done()
	zap.S().Info("server exiting")

	_ = r.closePostCb()
}
