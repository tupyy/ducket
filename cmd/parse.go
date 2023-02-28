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
	"github.com/tupyy/finance/internal/entity"
	"github.com/tupyy/finance/internal/parser"
	"github.com/tupyy/finance/internal/reader"
	"github.com/tupyy/finance/internal/writer/json"
)

var (
	rulesFiles []string
	file       string
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(rulesFiles) == 0 || file == "" {
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

		// read rules files
		rules := []entity.Rule{}
		for _, f := range rulesFiles {
			content, err := os.Open(f)
			if err != nil {
				return err
			}

			r, err := reader.ReadRules(content)
			if err != nil {
				return err
			}
			rules = append(rules, r...)
		}

		transactions := parser.Parse(context.Background(), records, rules)

		_ = json.Write(transactions)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringSliceVarP(&rulesFiles, "rules", "r", []string{}, "path to rules files")
	parseCmd.Flags().StringVarP(&file, "file", "f", "", "file to be parsed")
}
