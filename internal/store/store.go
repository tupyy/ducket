package store

import (
	"context"
	"database/sql"

	"git.tls.tupangiu.ro/cosmin/finante/internal/store/migrations"
)

type Store struct {
	db         *sql.DB
	qi         *queryInterceptor
	transactor *DBTransactor
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:         db,
		qi:         newQueryInterceptor(db),
		transactor: newTransactor(db),
	}
}

func (s *Store) Migrate(ctx context.Context) error {
	return migrations.Run(ctx, s.db)
}

func (s *Store) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.transactor.WithTx(ctx, fn)
}

func (s *Store) Checkpoint() error {
	_, err := s.db.Exec("FORCE CHECKPOINT")
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}
