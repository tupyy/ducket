package reader

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/tupyy/ducket/internal/entity"
)

// FileReader defines the interface for reading transaction data from different file formats
type FileReader interface {
	Read(r io.Reader) ([]entity.Transaction, error)
}

// NewFileReader creates a new FileReader based on the file extension.
// The account parameter is used by CSV readers to assign an account number
// to imported transactions (Excel files contain the account number in the sheet).
func NewFileReader(filename string, account int64) (FileReader, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".xlsx", ".xls":
		return &ExcelReader{}, nil
	case ".csv":
		return &CSVReader{Account: account}, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s. Supported types: .xlsx, .xls, .csv", ext)
	}
}

// ReadTransactionsFromFile reads transactions from a file using the appropriate reader
func ReadTransactionsFromFile(filename string, account int64, r io.Reader) ([]entity.Transaction, error) {
	reader, err := NewFileReader(filename, account)
	if err != nil {
		return nil, err
	}

	return reader.Read(r)
}
