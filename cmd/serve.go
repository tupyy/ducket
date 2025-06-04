package cmd

import (
	"context"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/config"
	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers"
	"git.tls.tupangiu.ro/cosmin/finante/internal/server"
	"git.tls.tupangiu.ro/cosmin/finante/pkg/logger"
	"github.com/ecordell/optgen/helpers"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/jzelinskie/cobrautil/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

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

			// init datastore
			dt, err := pg.NewPostgresDatastore(context.Background(), config.Database.URI)
			if err != nil {
				return err
			}

			server := server.NewRunnableServer(
				server.NewRunnableServerConfigWithOptionsAndDefaults(
					server.WithDatastore(dt),
					server.WithGraceTimeout(1*time.Second),
					server.WithPort(config.ServerPort),
					server.WithRegisterHandlersFn(func(r *gin.RouterGroup) {
						handlers.RegisterHandlers(r)
					}),
					server.WithCloseCallback(func() error {
						zap.S().Info("close datastore")
						dt.Close()
						return nil
					}),
				),
			)

			server.Run(context.Background())

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
	flagSet.StringVar(&config.URI, "db-conn-uri", config.URI, `connection string used by remote databases (e.g. "postgres://postgres:password@localhost:5432/photos")`)
	flagSet.BoolVar(&config.SSL, "db-ssl-mode", config.SSL, "ssl mode")
}

func registerServerFlags(flagSet *pflag.FlagSet, config *config.Config) {
	flagSet.IntVar(&config.ServerPort, "server-port", config.ServerPort, "port on which the server is listening")
	flagSet.StringVar(&config.GinMode, "server-gin-mode", config.GinMode, "gin mode: either release or debug. It applies only on server-type web")
}
