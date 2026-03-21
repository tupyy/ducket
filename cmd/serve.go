package cmd

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/config"
	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
	"git.tls.tupangiu.ro/cosmin/finante/pkg/logger"
	"github.com/ecordell/optgen/helpers"
	"github.com/fatih/color"
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
				return fmt.Errorf("statics folder should be provided in prod mod")
			}

			db, err := store.NewDB(config.Database.URI)
			if err != nil {
				return err
			}

			st := store.NewStore(db)
			defer st.Close()

			if err := st.Migrate(context.Background()); err != nil {
				return fmt.Errorf("running migrations: %w", err)
			}

			// TODO: wire up HTTP server and handlers
			zap.S().Info("store ready, server not yet implemented")
			return nil
		},
	}

	registerFlags(cmd, config)

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
	flagSet.StringVar(&config.Mode, "server-mode", config.Mode, "server mod: dev or prod")
	flagSet.StringVar(&config.StaticsFolder, "statics-folder", config.StaticsFolder, "path to statics")
}
