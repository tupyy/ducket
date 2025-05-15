package json

import (
	"encoding/json"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

func Write(transactions []*entity.Transaction) error {
	out, err := json.Marshal(transactions)
	if err != nil {
		return nil
	}
	fmt.Printf("%s", out)
	return nil
}
