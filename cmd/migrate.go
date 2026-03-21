package cmd

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/config"
	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
	"git.tls.tupangiu.ro/cosmin/finante/pkg/logger"
	"github.com/fatih/color"
	"github.com/go-extras/cobraflags"
	"github.com/jzelinskie/cobrautil/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewMigrateCommand creates a new cobra command for running database migrations.
func NewMigrateCommand(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "migrate",
		Short:        "Run database migrations",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logger.SetupLogger(config)
			defer logger.Sync()

			undo := zap.ReplaceGlobals(logger)
			defer undo()

			zap.S().Infow("starting database migration", "db_uri", config.Database.URI)

			db, err := store.NewDB(config.Database.URI)
			if err != nil {
				zap.S().Errorw("failed to connect to database", "error", err)
				return err
			}
			defer db.Close()

			zap.S().Info("connected to database successfully")

			st := store.NewStore(db)
			if err := st.Migrate(context.Background()); err != nil {
				zap.S().Errorw("migration failed", "error", err)
				return err
			}

			zap.S().Info("migrations completed successfully")
			return nil
		},
	}

	registerMigrateFlags(cmd, config)
	cobraflags.CobraOnInitialize("FINANTE", cmd)

	return cmd
}

func registerMigrateFlags(cmd *cobra.Command, config *config.Config) {
	nfs := cobrautil.NewNamedFlagSets(cmd)

	dbFlagSet := nfs.FlagSet(color.New(color.FgCyan, color.Bold).Sprint("database"))
	registerDatabaseFlags(dbFlagSet, config.Database)

	nfs.AddFlagSets(cmd)
}
