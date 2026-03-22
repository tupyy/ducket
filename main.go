package main

import (
	"fmt"
	"os"

	"github.com/tupyy/ducket/cmd"
	"github.com/tupyy/ducket/internal/config"
	"github.com/spf13/cobra"
)

var sha string

func main() {
	cfg := config.NewConfigWithOptionsAndDefaults(
		config.WithDatabase(config.NewDatabaseWithOptions(
			config.WithURI("ducket.db"),
		)),
		config.WithServerPort(8080),
		config.WithLogFormat("console"),
		config.WithLogLevel("debug"),
		config.WithGinMode("debug"),
	)

	fmt.Printf("Build from commit: %s\n", sha)

	var rootCmd = &cobra.Command{
		Use:   "ducket",
		Short: "Personal finance tracker",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	registerLoggingFlags(rootCmd, cfg)

	rootCmd.AddCommand(cmd.NewRunCommand(cfg))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func registerLoggingFlags(cmd *cobra.Command, config *config.Config) {
	cmd.PersistentFlags().StringVar(&config.LogFormat, "log-format", config.LogFormat, "format of the logs: console or json")
	cmd.PersistentFlags().StringVar(&config.LogLevel, "log-level", config.LogLevel, "log level")
}
