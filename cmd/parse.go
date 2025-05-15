/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"git.tls.tupangiu.ro/cosmin/finante/internal/parser"
	"git.tls.tupangiu.ro/cosmin/finante/internal/reader"
	"git.tls.tupangiu.ro/cosmin/finante/internal/repo"
	"git.tls.tupangiu.ro/cosmin/finante/internal/writer/postgres"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	rulesFile string
	file      string
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

		if rulesFile == "" || file == "" {
			return errors.New("rules file or/and file to be parse is missing")
		}
		// read excel file
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		excelReader := &reader.ExcelReader{}
		records, err := excelReader.Read(f)
		if err != nil {
			return fmt.Errorf("unable to read records from file %q: %w", file, err)
		}

		rules, err := reader.ReadRules(rulesFile)
		if err != nil {
			return err
		}

		zap.S().Info(rules)
		transactions := parser.Parse(context.Background(), records, rules)

		pgClient, err := repo.New(repo.ClientParams{
			Host:     "fedorasrv",
			Port:     5432,
			DBName:   "finance",
			User:     "postgres",
			Password: "postgres",
		})
		if err != nil {
			zap.S().Fatal(err)
		}

		pgRepo, err := repo.NewRepo(pgClient)
		if err != nil {
			zap.S().Fatal(err)
		}

		pgWriter := postgres.NewPgWriter(pgRepo)
		err = pgWriter.Write(context.Background(), transactions)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&rulesFile, "rules", "r", "", "path to rules files")
	parseCmd.Flags().StringVarP(&file, "file", "f", "", "file to be parsed")
}
