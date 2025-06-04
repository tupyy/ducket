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
	"go.uber.org/zap"
)

//go:generate go run github.com/ecordell/optgen -output zz_configuration.go . RunnableServerConfig
type RunnableServerConfig struct {
	PostgresURI  string        `debugmap:"visible"`
	GraceTimeout time.Duration `debugmap:"visible"`
	Port         int           `debugmap:"visible"`
}

type runnableServer struct {
	dt  *pg.Datastore
	srv *http.Server
	cfg *RunnableServerConfig
}

func NewRunnableServer(cfg *RunnableServerConfig) *runnableServer {
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", cfg.Port),
	}

	return &runnableServer{srv: srv, cfg: cfg}
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
}
