package pg

// copyright SpiceDB

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg/models"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryFilter func(original sq.SelectBuilder) sq.SelectBuilder

type Writer interface {
	WriteTransaction(ctx context.Context, transaction entity.Transaction) (int64, error)
	DeleteTransaction(ctx context.Context, id int64) error
	WriteLabel(ctx context.Context, label entity.Label) error
	DeleteLabel(ctx context.Context, label entity.Label) error
	WriteRule(ctx context.Context, rule entity.Rule, update bool) error
	DeleteRule(ctx context.Context, id string) error
}

type TxUserFunc func(context.Context, Writer) error

type Datastore struct {
	pool ConnPooler
}

// NewPostgresDatastore creates a new Postgres datastore instance with the given configuration options.
// It establishes a connection pool and sets up query interceptors for logging and monitoring.
func NewPostgresDatastore(ctx context.Context, url string, options ...Option) (*Datastore, error) {
	pgOptions := newPostgresConfig(options)

	pgxConfig, err := pgOptions.PgxConfig(url)
	if err != nil {
		return nil, err
	}

	pgPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	if err := pgPool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Datastore{pool: MustNewInterceptorPooler(pgPool, newLogInterceptor())}, nil
}

// QueryTransactions retrieves transactions from the database based on the provided query filters.
func (d *Datastore) QueryTransactions(ctx context.Context, filterFn ...QueryFilter) ([]entity.Transaction, error) {
	query := psql.Select(
		colID,
		colDate,
		colTransactionAccount,
		colTransactionType,
		colTransactionContent,
		colTransactionAmount,
		"label_id",
		colRuleID,
		colHash).
		From(transactionTable).
		LeftJoin("transactions_labels ON transactions_labels.transaction_id = transactions.id")

	for _, fn := range filterFn {
		query = fn(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return []entity.Transaction{}, fmt.Errorf(errUnableToReadRule, err)
	}

	rows, err := d.pool.Query(ctx, sql, args...)
	if err != nil {
		return []entity.Transaction{}, fmt.Errorf(errUnableToReadRule, err)
	}
	defer rows.Close()

	tRows := make(models.Transactions)
	rs := pgxscan.NewRowScanner(rows)

	for rows.Next() {
		tRow := models.Transaction{}
		err := rs.Scan(&tRow)
		if err != nil {
			return []entity.Transaction{}, fmt.Errorf(errUnableToReadRule, err)
		}
		tRows.Add(tRow)
	}

	return tRows.ToEntity(), nil
}

// QueryRules retrieves rules from the database based on the provided query filters.
func (d *Datastore) QueryRules(ctx context.Context, filter ...QueryFilter) ([]entity.Rule, error) {
	query := selectRulesStmt

	for _, f := range filter {
		query = f(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return []entity.Rule{}, fmt.Errorf(errUnableToReadRule, err)
	}

	rows, err := d.pool.Query(ctx, sql, args...)
	if err != nil {
		return []entity.Rule{}, fmt.Errorf(errUnableToReadRule, err)
	}
	defer rows.Close()

	ruleRows := make(models.RuleRows)
	rs := pgxscan.NewRowScanner(rows)

	for rows.Next() {
		ruleRow := models.RuleRow{}
		err := rs.Scan(&ruleRow)
		if err != nil {
			return []entity.Rule{}, fmt.Errorf(errUnableToReadRule, err)
		}
		ruleRows.Add(ruleRow)
	}

	return ruleRows.ToEntity(), nil
}

// QueryLabels retrieves labels from the database based on the provided query filters.
func (d *Datastore) QueryLabels(ctx context.Context, filter ...QueryFilter) ([]entity.Label, error) {
	query := selectLabelsStmt

	for _, f := range filter {
		query = f(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return []entity.Label{}, fmt.Errorf(errUnableToReadLabel, err)
	}

	rows, err := d.pool.Query(ctx, sql, args...)
	if err != nil {
		return []entity.Label{}, fmt.Errorf(errUnableToReadLabel, err)
	}

	labels := make(models.Labels)
	rs := pgxscan.NewRowScanner(rows)

	for rows.Next() {
		row := models.Label{}
		err := rs.Scan(&row)
		if err != nil {
			return []entity.Label{}, fmt.Errorf(errUnableToReadLabel, err)
		}
		labels.Add(row)
	}

	return labels.ToEntity(), nil
}

// CountTransactions returns transaction statistics grouped by labels for reporting purposes.
func (d *Datastore) CountTransactions(ctx context.Context) ([]entity.TransactionStat, error) {
	sql, args, err := countTransactionsPerLabelPerRuleStmt.ToSql()
	if err != nil {
		return []entity.TransactionStat{}, fmt.Errorf(errUnableToReadLabel, err)
	}

	rows, err := d.pool.Query(ctx, sql, args...)
	if err != nil {
		return []entity.TransactionStat{}, fmt.Errorf(errUnableToReadLabel, err)
	}

	stats := make([]entity.TransactionStat, 0)
	rs := pgxscan.NewRowScanner(rows)

	for rows.Next() {
		row := models.TransactionCountRow{}
		err := rs.Scan(&row)
		if err != nil {
			return []entity.TransactionStat{}, fmt.Errorf(errUnableToReadLabel, err)
		}
		stats = append(stats, entity.TransactionStat{LabelID: row.LabelID, RuleID: row.RuleID, Count: row.Count})
	}

	return stats, nil
}

// WriteTx executes a write transaction with the provided user function.
// It manages transaction lifecycle and provides a Writer interface for data modifications.
func (d *Datastore) WriteTx(ctx context.Context, txFn TxUserFunc) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}

	writer := &writerTx{tx: tx}

	if err := txFn(ctx, writer); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// Close gracefully shuts down the datastore and closes all database connections.
func (d *Datastore) Close() {
	d.pool.Close()
}
