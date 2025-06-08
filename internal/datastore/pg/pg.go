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
	WriteTag(ctx context.Context, value string) error
	DeleteTag(ctx context.Context, value string) error
	WriteRule(ctx context.Context, rule entity.Rule, update bool) error
	DeleteRule(ctx context.Context, id string) error
}

type TxUserFunc func(context.Context, Writer) error

type Datastore struct {
	pool ConnPooler
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

	return &Datastore{pool: MustNewInterceptorPooler(pgPool, newLogInterceptor())}, nil
}

func (d *Datastore) QueryTransactions(ctx context.Context, filterFn ...QueryFilter) ([]entity.Transaction, error) {
	query := psql.Select(colID, colDate, colTransactionType, colTransactionContent, colAmount, colTagID, colRuleID, colHash).
		From(transactionTable).
		LeftJoin("transactions_tags ON transactions_tags.transaction_id = transactions.id")

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

func (d *Datastore) QueryTags(ctx context.Context, filter ...QueryFilter) ([]string, error) {
	query := selectTagsStmt

	for _, f := range filter {
		query = f(query)
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
