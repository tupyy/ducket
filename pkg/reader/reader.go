package reader

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

// FileReader defines the interface for reading transaction data from different file formats
type FileReader interface {
	Read(r io.Reader) ([]entity.Transaction, error)
}

// NewFileReader creates a new FileReader based on the file extension
func NewFileReader(filename string) (FileReader, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".xlsx", ".xls":
		return &ExcelReader{}, nil
	case ".csv":
		return &CSVReader{}, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s. Supported types: .xlsx, .xls, .csv", ext)
	}
}

// ReadTransactionsFromFile reads transactions from a file using the appropriate reader
func ReadTransactionsFromFile(filename string, r io.Reader) ([]entity.Transaction, error) {
	reader, err := NewFileReader(filename)
	if err != nil {
		return nil, err
	}

	return reader.Read(r)
}
