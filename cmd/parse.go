/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tupyy/finance/internal/parser"
	"github.com/tupyy/finance/internal/reader"
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
		_ = parser.Parse(context.Background(), records, rules)

		// token := "MPAmH3mZnTnJ1PTBHniof1FQhBzKPnnJ7ngHkyqJZgWU6ct8qHdrjZ6ZBFNSlZW-obSZuk0Mb5mH-UrmxAgZrA=="
		// bucket := "finance"
		// org := "home"

		// influxWriter := influxdb.InfluxWriter{
		// 	Url:    "http://localhost:8086",
		// 	Token:  token,
		// 	Org:    org,
		// 	Bucket: bucket,
		// }

		// _ = influxWriter.Write(transactions)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&rulesFile, "rules", "r", "", "path to rules files")
	parseCmd.Flags().StringVarP(&file, "file", "f", "", "file to be parsed")
}
