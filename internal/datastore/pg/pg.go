package pg

import (
	"context"
	"fmt"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg/models"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go run github.com/ecordell/optgen -output zz_generated.transaction_filter.go . TransactionFilter
type TransactionFilter struct {
	After  *time.Time `debugmap:"visible"`
	Before *time.Time `debugmap:"visible"`
}

//xgo:generate go run github.com/ecordell/optgen -output zz_generated.rule_filter.go . RuleFilter
type RuleFilter struct {
}

type TagFilter struct {
	RuleID *string
}

//go:generate go run github.com/ecordell/optgen -output zz_generated.query_transaction_opts.go . QueryTransactionOptions
type QueryTransactionOptions struct {
	Limit  int `debugmap:"visible"`
	Offset int `debugmap:"visible"`
}

//xgo:generate go run github.com/ecordell/optgen -output zz_generated.query_rule_opts.go . QueryRuleOptions
type QueryRuleOptions struct {
}

type Writer interface {
	WriteTransaction(ctx context.Context, transaction entity.Transaction) error
	DeleteTransaction(ctx context.Context, id string) error
	WriteTag(ctx context.Context, value string) error
	DeleteTag(ctx context.Context, value string) error
	WriteRule(ctx context.Context, rule entity.Rule, update bool) error
	DeleteRule(ctx context.Context, id string) error
}

type TxUserFunc func(context.Context, Writer) error

type Datastore struct {
	pool *pgxpool.Pool
}

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

	return &Datastore{pool: pgPool}, nil
}

func (d *Datastore) QueryTransactions(ctx context.Context, filter TransactionFilter, opts *QueryTransactionOptions) ([]entity.Transaction, error) {
	query := psql.Select(colID, colDate, colTransactionType, colTransactionContent, colAmount, colTagID, colRuleID, colHash).
		From(transactionTable).
		LeftJoin("transactions_tags ON transactions_tags.transaction_id = transactions.id")

	if filter.Before != nil && filter.After != nil {
		query = query.Where(
			sq.And{
				sq.GtOrEq{"date": filter.After},
				sq.LtOrEq{"date": filter.Before},
			},
		)
	} else {
		if filter.Before != nil {
			query = query.Where(sq.LtOrEq{"date": filter.Before})
		}
		if filter.After != nil {
			query = query.Where(sq.GtOrEq{"date": filter.After})
		}
	}

	if opts.Limit > 0 {
		query = query.Limit(uint64(opts.Limit))
	}

	if opts.Offset > 0 {
		query = query.Offset(uint64(opts.Offset))
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

func (d *Datastore) QueryRules(ctx context.Context, filter RuleFilter, opts *QueryRuleOptions) ([]entity.Rule, error) {
	query := selectRulesStmt
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

func (d *Datastore) QueryTags(ctx context.Context, filter TagFilter) ([]string, error) {
	query := selectTagsStmt

	if filter.RuleID != nil {
		query = query.Where(sq.Eq{colRuleID: filter.RuleID})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return []string{}, fmt.Errorf(errUnableToReadTag, err)
	}

	rows, err := d.pool.Query(ctx, sql, args...)
	if err != nil {
		return []string{}, fmt.Errorf(errUnableToReadTag, err)
	}

	tags := make(models.Tags)
	rs := pgxscan.NewRowScanner(rows)

	for rows.Next() {
		tagRow := models.Tag{}
		err := rs.Scan(&tagRow)
		if err != nil {
			return []string{}, fmt.Errorf(errUnableToReadTag, err)
		}
		tags.Add(tagRow)
	}

	return tags.ToEntity(), nil

}

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

func (d *Datastore) Close() {
	d.pool.Close()
}
