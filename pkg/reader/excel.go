package reader

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	excel "github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

const (
	dateFormat = "02/01/2006"
)

type ExcelReader struct{}

// Read parses an Excel file from the provided io.Reader and extracts transaction data.
// It looks for a sheet named "Sheet0" and starts reading transactions after finding
// a row that begins with "date". Returns a slice of Transaction entities.
func (e *ExcelReader) Read(r io.Reader) ([]*entity.Transaction, error) {
	f, err := excel.OpenReader(r, excel.Options{})
	if err != nil {
		return []*entity.Transaction{}, err
	}

	rows, err := f.GetRows("Sheet0")
	startRead := false
	transactions := make([]*entity.Transaction, 0, len(rows))
	for _, row := range rows {
		if len(row) > 0 && strings.ToLower(row[0]) == "date" {
			startRead = true
		}
		if startRead {
			t, err := parseRow(row)
			if err != nil {
				zap.S().Error(err)
				continue
			}
			transactions = append(transactions, t)
		}
	}

	return transactions, nil
}

func parseRow(r []string) (*entity.Transaction, error) {
	if len(r) < 3 {
		return nil, fmt.Errorf("data is not a transaction %q", r)
	}

	date, err := parseDate(r[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert to date %q: %w", r[0], err)
	}
	rowContent := formatContent(r[1])

	var sum string
	kind := entity.CreditTransaction
	if len(r) == 3 {
		kind = entity.DebitTransaction
		sum = r[2]
	} else {
		sum = r[3]
	}

	floatSum, err := parseSum(sum)
	if err != nil {
		return nil, fmt.Errorf("cannot convert sum %q: %w", sum, err)
	}

	return entity.NewTransaction(kind, date, floatSum, rowContent), nil
}

func parseDate(s string) (time.Time, error) {
	return time.Parse(dateFormat, s)
}

func parseSum(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func formatContent(s string) string {
	ss := strings.ReplaceAll(s, "\n", " ")
	parts := strings.Split(ss, " ")
	trimmedParts := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		trimmedParts = append(trimmedParts, strings.TrimSpace(p))
	}
	return strings.ToLower(strings.Join(trimmedParts, " "))
}
