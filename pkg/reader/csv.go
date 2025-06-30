package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"go.uber.org/zap"
)

type CSVReader struct{}

// Read parses a CSV file from the provided io.Reader and extracts transaction data.
// It expects the first row to contain headers and starts reading transactions after finding
// a row that begins with "date". Returns a slice of Transaction entities.
func (c *CSVReader) Read(r io.Reader) ([]entity.Transaction, error) {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		return []entity.Transaction{}, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return []entity.Transaction{}, fmt.Errorf("CSV file is empty")
	}

	startRead := false
	transactions := make([]entity.Transaction, 0, len(records))

	for _, record := range records {
		if len(record) > 0 && strings.ToLower(strings.TrimSpace(record[0])) == "date" {
			startRead = true
			continue // Skip the header row
		}

		if startRead {
			t, err := parseCSVRow(record)
			if err != nil {
				zap.S().Errorw("failed to parse CSV row", "error", err, "row", record)
				continue
			}
			transactions = append(transactions, *t)
		}
	}

	return transactions, nil
}

func parseCSVRow(record []string) (*entity.Transaction, error) {
	if len(record) < 3 {
		return nil, fmt.Errorf("CSV row does not contain enough columns: %q", record)
	}

	// Parse date from first column
	date, err := parseCSVDate(strings.TrimSpace(record[0]))
	if err != nil {
		return nil, fmt.Errorf("cannot parse date %q: %w", record[0], err)
	}

	// Parse content from second column
	content := formatContent(strings.TrimSpace(record[1]))

	// Determine transaction type and amount based on number of columns
	var sum string
	kind := entity.CreditTransaction

	if len(record) == 3 {
		// Format: date, content, amount (debit)
		kind = entity.DebitTransaction
		sum = strings.TrimSpace(record[2])
	} else if len(record) >= 4 {
		// Format: date, content, debit_amount, credit_amount
		debitAmount := strings.TrimSpace(record[2])
		creditAmount := strings.TrimSpace(record[3])

		if debitAmount != "" && debitAmount != "0" {
			kind = entity.DebitTransaction
			sum = debitAmount
		} else if creditAmount != "" && creditAmount != "0" {
			kind = entity.CreditTransaction
			sum = creditAmount
		} else {
			return nil, fmt.Errorf("no valid amount found in row: %q", record)
		}
	}

	floatSum, err := parseCSVSum(sum)
	if err != nil {
		return nil, fmt.Errorf("cannot parse amount %q: %w", sum, err)
	}

	return entity.NewTransaction(kind, date, floatSum, content), nil
}

func parseCSVDate(s string) (time.Time, error) {
	// Try multiple date formats commonly used in CSV files
	formats := []string{
		"02/01/2006", // DD/MM/YYYY
		"01/02/2006", // MM/DD/YYYY
		"2006-01-02", // YYYY-MM-DD
		"2006/01/02", // YYYY/MM/DD
		"02-01-2006", // DD-MM-YYYY
		"01-02-2006", // MM-DD-YYYY
	}

	for _, format := range formats {
		if date, err := time.Parse(format, s); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date with any known format")
}

func parseCSVSum(s string) (float32, error) {
	// Clean the string: remove currency symbols, spaces, and handle different decimal separators
	cleaned := strings.ReplaceAll(s, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "£", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	// Handle negative values
	if strings.HasPrefix(cleaned, "-") {
		f, err := strconv.ParseFloat(cleaned[1:], 32)
		if err != nil {
			return 0, err
		}
		return -float32(f), nil
	}

	f, err := strconv.ParseFloat(cleaned, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}
