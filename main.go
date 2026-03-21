package main

import (
	"fmt"
	"os"

	"git.tls.tupangiu.ro/cosmin/finante/cmd"
	"git.tls.tupangiu.ro/cosmin/finante/internal/config"
	"github.com/spf13/cobra"
)

var sha string

func main() {
	cfg := config.NewConfigWithOptionsAndDefaults(
		config.WithDatabase(config.NewDatabaseWithOptions(
			config.WithURI("finante.db"),
		)),
		config.WithServerPort(8080),
		config.WithLogFormat("console"),
		config.WithLogLevel("debug"),
		config.WithGinMode("debug"),
	)

	fmt.Printf("Build from commit: %s\n", sha)

	var rootCmd = &cobra.Command{
		Use:   "finance",
		Short: "Manage my finances",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	registerLoggingFlags(rootCmd, cfg)

	rootCmd.AddCommand(cmd.NewServeCommand(cfg))
	rootCmd.AddCommand(cmd.NewMigrateCommand(cfg))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func registerLoggingFlags(cmd *cobra.Command, config *config.Config) {
	cmd.PersistentFlags().StringVar(&config.LogFormat, "log-format", config.LogFormat, "format of the logs: console or json")
	cmd.PersistentFlags().StringVar(&config.LogLevel, "log-level", config.LogLevel, "log level")
}
