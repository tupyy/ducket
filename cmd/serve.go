package cmd

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/config"
	"git.tls.tupangiu.ro/cosmin/finante/internal/server"
	"github.com/fatih/color"
	"github.com/jzelinskie/cobrautil/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewServeCommand(config *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:          "serve",
		Short:        "Serve the server either web or admin",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			server := server.CreateRunnableServer(ctx, config)

			// run server
			server.Run(fmt.Sprintf("0.0.0.0:%d", config.ServerPort))

			return nil
		},
	}
}

func RegisterFlags(cmd *cobra.Command, config *config.Config) {
	nfs := cobrautil.NewNamedFlagSets(cmd)

	dbFlagSet := nfs.FlagSet(color.New(color.FgCyan, color.Bold).Sprint("database"))
	registerDatabaseFlags(dbFlagSet, config.Database)

	serverFlagSet := nfs.FlagSet(color.New(color.FgCyan, color.Bold).Sprint("server"))
	registerServerFlags(serverFlagSet, config)

	logFlagSet := nfs.FlagSet(color.New(color.FgCyan, color.Bold).Sprint("log"))
	registerLoggingFlags(logFlagSet, config)

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

func registerLoggingFlags(flagSet *pflag.FlagSet, config *config.Config) {
	flagSet.StringVar(&config.LogFormat, "log-format", config.LogFormat, "format of the logs: console or json")
	flagSet.StringVar(&config.LogLevel, "log-level", config.LogLevel, "log level")
}
