package reader

import (
	"io"

	"github.com/tupyy/finance/internal/entity"
)

type ExcelReader struct{}

func (e *ExcelReader) Read(r io.Reader) ([]*entity.Transaction, error) {
	return []*entity.Transaction{}, nil
}
