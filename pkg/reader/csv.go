package reader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/tupyy/ducket/internal/entity"
	"go.uber.org/zap"
)

type CSVReader struct {
	Account int64
}

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

	transactions := make([]entity.Transaction, 0, len(records))

	for _, record := range records {
		t, err := c.parseRow(record)
		if err != nil {
			zap.S().Errorw("failed to parse CSV row", "error", err, "row", record)
			continue
		}
		transactions = append(transactions, *t)
	}

	return transactions, nil
}

func (c *CSVReader) parseRow(r []string) (*entity.Transaction, error) {
	if len(r) < 10 {
		return nil, fmt.Errorf("data is not a transaction %q", r)
	}

	if r[8] != "COMPLETED" {
		return nil, errors.New("transaction is not completed")
	}

	floatSum, err := c.parseSum(r[5])
	if err != nil {
		return nil, fmt.Errorf("cannot convert sum %q: %w", r[5], err)
	}

	kind := entity.CreditTransaction
	switch r[0] {
	case "CARD_PAYMENT":
		kind = entity.DebitTransaction
	case "TOPUP":
		kind = entity.CreditTransaction
	case "TRANSFER":
		if floatSum < 0 {
			kind = entity.DebitTransaction
		} else {
			kind = entity.CreditTransaction
		}
	}

	if floatSum < 0 {
		floatSum = floatSum * -1
	}

	date, err := c.parseDate(r[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert to date %q: %w", r[2], err)
	}
	rowContent := formatContent(r[4])

	return entity.NewTransaction(kind, c.Account, date, floatSum, rowContent), nil
}

func (c *CSVReader) parseDate(s string) (time.Time, error) {
	format := "2006-01-02" // YYYY-MM-DD

	parts := strings.Split(s, " ")

	if date, err := time.Parse(format, parts[0]); err == nil {
		return date, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse date with any known format")
}

func (c *CSVReader) parseSum(s string) (float64, error) {
	cleaned := strings.ReplaceAll(s, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "£", "")

	// Handle European format (e.g. "1.234,56") by stripping thousand-separator
	// dots before replacing the decimal comma with a dot.
	if strings.Contains(cleaned, ",") {
		cleaned = strings.ReplaceAll(cleaned, ".", "")
		cleaned = strings.ReplaceAll(cleaned, ",", ".")
	}

	f, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
