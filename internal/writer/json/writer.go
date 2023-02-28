package json

import (
	"encoding/json"
	"fmt"

	"github.com/tupyy/finance/internal/entity"
)

func Write(transactions []*entity.Transaction) error {
	out, err := json.Marshal(transactions)
	if err != nil {
		return nil
	}
	fmt.Printf("%s", out)
	return nil
}
