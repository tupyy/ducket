package pg

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"github.com/jackc/pgx/v5"
)

type ruleWriterTx struct {
	tx pgx.Tx
}

func (r *ruleWriterTx) WriteRule(ctx context.Context, rule entity.Rule) error {
	return nil
}

func (r *ruleWriterTx) DeleteRule(ctx context.Context, id string) error {
	return nil
}

func (r *ruleWriterTx) WriteTag(ctx context.Context, value string, ruleID int) error {
	return nil
}

func (r *ruleWriterTx) DeleteTag(ctx context.Context, value string) error {
	return nil
}
