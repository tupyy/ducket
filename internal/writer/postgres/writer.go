package postgres

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/repo"
	"go.uber.org/zap"
)

type PgWriter struct {
	pg *repo.TransationRepo
}

func NewPgWriter(pgRepo *repo.TransationRepo) *PgWriter {
	return &PgWriter{pgRepo}
}

func (w *PgWriter) Write(ctx context.Context, transactions []*entity.Transaction) error {
	for _, t := range transactions {
		if err := w.pg.Write(ctx, t); err != nil {
			zap.S().Error(err)
			continue
		}
	}
	return nil
}
