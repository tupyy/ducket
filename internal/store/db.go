package store

import (
	"database/sql"

	_ "github.com/duckdb/duckdb-go/v2"
)

func NewDB(path string) (*sql.DB, error) {
	conn, err := sql.Open("duckdb", path)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(1)

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}
