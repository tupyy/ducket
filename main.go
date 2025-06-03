package main

import (
	"os"

	"git.tls.tupangiu.ro/cosmin/finante/cmd"
	"git.tls.tupangiu.ro/cosmin/finante/internal/config"
	"git.tls.tupangiu.ro/cosmin/finante/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "finance",
		Short: "Manage my finances",
	}

	cfg := config.NewConfigWithOptionsAndDefaults(
		config.WithDatabase(config.NewDatabaseWithOptions(
			config.WithURI("postgres://postgres:postgres@localhost:5432/photos?sslmode=disable"),
		)),
		config.WithServerPort(8080),
		config.WithLogFormat("console"),
		config.WithLogLevel("debug"),
		config.WithGinMode("debug"),
	)

	logger := logger.SetupLogger(cfg)
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	serveCmd := cmd.NewServeCommand(cfg)
	cmd.RegisterFlags(serveCmd, cfg)

	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
