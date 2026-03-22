package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tupyy/ducket/internal/config"
	"github.com/tupyy/ducket/internal/handlers"
	"github.com/tupyy/ducket/internal/server"
	"github.com/tupyy/ducket/internal/services"
	"github.com/tupyy/ducket/internal/store"
	"github.com/tupyy/ducket/pkg/logger"
	"github.com/ecordell/optgen/helpers"
	"github.com/fatih/color"
	"github.com/go-extras/cobraflags"
	"github.com/jzelinskie/cobrautil/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

// NewServeCommand creates a new cobra command for starting the server.
func NewServeCommand(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "serve",
		Short:        "Serve the server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logger.SetupLogger(config)
			defer logger.Sync()

			undo := zap.ReplaceGlobals(logger)
			defer undo()

			zap.S().Info("using configuration", "config", helpers.Flatten(config.DebugMap()))

			if config.Mode == "prod" && config.StaticsFolder == "" {
				return fmt.Errorf("statics folder should be provided in prod mode")
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
			defer cancel()

			// Database
			db, err := store.NewDB(config.Database.URI)
			if err != nil {
				return err
			}

			st := store.NewStore(db)
			defer st.Close()

			if err := st.Migrate(ctx); err != nil {
				return fmt.Errorf("running migrations: %w", err)
			}

			// Services
			txnSvc := services.NewTransactionService(st)
			ruleSvc := services.NewRuleService(st)
			summarySvc := services.NewSummaryService(st)

			// Handler
			h := handlers.NewHandler(txnSvc, ruleSvc, summarySvc)

			// HTTP server
			srv := server.New(config.ServerPort, config.GinMode, h)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
					cancel()
				}()
				zap.S().Infof("starting HTTP server on port %d", config.ServerPort)
				if err := srv.Start(); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						zap.S().Errorw("failed to start http server", "error", err)
					}
				}
			}()

			go func() {
				<-ctx.Done()
				stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer stopCancel()
				srv.Stop(stopCtx)
			}()

			<-ctx.Done()
			wg.Wait()

			zap.S().Info("server shutdown complete")
			return nil
		},
	}

	registerFlags(cmd, config)
	cobraflags.CobraOnInitialize("FINANTE", cmd)

	return cmd
}

func registerFlags(cmd *cobra.Command, config *config.Config) {
	nfs := cobrautil.NewNamedFlagSets(cmd)

	dbFlagSet := nfs.FlagSet(color.New(color.FgCyan, color.Bold).Sprint("database"))
	registerDatabaseFlags(dbFlagSet, config.Database)

	serverFlagSet := nfs.FlagSet(color.New(color.FgCyan, color.Bold).Sprint("server"))
	registerServerFlags(serverFlagSet, config)

	nfs.AddFlagSets(cmd)
}

func registerDatabaseFlags(flagSet *pflag.FlagSet, config *config.Database) {
	flagSet.StringVar(&config.URI, "db-uri", config.URI, `path to DuckDB database file (e.g. "./finante.db" or ":memory:")`)
}

func registerServerFlags(flagSet *pflag.FlagSet, config *config.Config) {
	flagSet.IntVar(&config.ServerPort, "server-port", config.ServerPort, "port on which the server is listening")
	flagSet.StringVar(&config.GinMode, "server-gin-mode", config.GinMode, "gin mode: either release or debug. It applies only on server-type web")
	flagSet.StringVar(&config.Mode, "server-mode", config.Mode, "server mode: dev or prod")
	flagSet.StringVar(&config.StaticsFolder, "statics-folder", config.StaticsFolder, "path to statics")
}
