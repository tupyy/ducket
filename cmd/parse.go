/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	dryRun bool
	file   string
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := setupLogger()
		defer logger.Sync()

		undo := zap.ReplaceGlobals(logger)
		defer undo()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "dry-run")
	parseCmd.Flags().StringVarP(&file, "file", "f", "", "file to be parsed")
}
