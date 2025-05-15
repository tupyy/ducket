package stdout

import (
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
)

func Write(transactions []*entity.Transaction) error {
	for _, t := range transactions {
		fmt.Println(t)
	}
	return nil
}
