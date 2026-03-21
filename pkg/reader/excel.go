package reader

import (
	"fmt"
	"io"
	"regexp"
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

var (
	accountRegexp = regexp.MustCompile(`^Compte.*([0-9]{11})`)
)

type ExcelReader struct{}

// Read parses an Excel file from the provided io.Reader and extracts transaction data.
// It looks for a sheet named "Sheet0" and starts reading transactions after finding
// a row that begins with "date". Returns a slice of Transaction entities.
func (e *ExcelReader) Read(r io.Reader) ([]entity.Transaction, error) {
	f, err := excel.OpenReader(r, excel.Options{})
	if err != nil {
		return []entity.Transaction{}, err
	}
	defer f.Close()

	accountNumber := int64(0)
	rows, err := f.GetRows("Sheet0")
	if err != nil {
		return []entity.Transaction{}, fmt.Errorf("failed to get rows from Sheet0: %w", err)
	}

	startReadTransactions := false
	transactions := make([]entity.Transaction, 0, len(rows))
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		if matched := accountRegexp.MatchString(row[0]); matched {
			sAccountNumber := accountRegexp.FindStringSubmatch(row[0])
			if len(sAccountNumber) > 0 {
				if a, err := strconv.ParseInt(sAccountNumber[1], 10, 64); err == nil {
					accountNumber = a
				}
			}
		}

		if strings.ToLower(row[0]) == "date" {
			startReadTransactions = true
			continue
		}

		if startReadTransactions {
			t, err := parseRow(row)
			if err != nil {
				zap.S().Errorw("failed to parse Excel row", "error", err, "row", row)
				continue
			}
			t.Account = accountNumber
			transactions = append(transactions, *t)
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
	if len(r) >= 4 && r[3] != "" {
		sum = r[3]
	} else if len(r) >= 3 && r[2] != "" {
		kind = entity.DebitTransaction
		sum = r[2]
	} else {
		return nil, fmt.Errorf("no amount found in row %q", r)
	}

	floatSum, err := parseSum(sum)
	if err != nil {
		return nil, fmt.Errorf("cannot convert sum %q: %w", sum, err)
	}

	return entity.NewTransaction(kind, 0, date, floatSum, rowContent), nil
}

func parseDate(s string) (time.Time, error) {
	return time.Parse(dateFormat, s)
}

func parseSum(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
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
