package stdout

import (
	"fmt"

	"github.com/tupyy/finance/internal/entity"
)

func Write(transactions []*entity.Transaction) error {
	for _, t := range transactions {
		fmt.Println(t)
	}
	return nil
}
