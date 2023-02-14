/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tupyy/finance/internal/parser"
	"github.com/tupyy/finance/internal/reader"
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
		if rulesFile == "" || file == "" {
			return errors.New("rules file or/and file to be parse is missing")
		}
		// read excel file
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		reader := &reader.ExcelReader{}
		records, err := reader.Read(f)
		if err != nil {
			return fmt.Errorf("unable to read records from file %q: %w", file, err)
		}

		parsedRecords, err := parser.Parse(ctx, records, rules)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringVarP(&rulesFile, "rules", "r", "", "path to rules files")
	parseCmd.Flags().StringVarP(&file, "file", "f", "", "file to be parsed")
}
