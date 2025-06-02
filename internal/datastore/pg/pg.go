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

type SortBy int8

const (
	CAPTURED_AT SortBy = iota
)

const (
	rulesTable                = "rules"
	tagsTable                 = "tags"
	transactionTable          = "transactions"
	transactionsTagsTable     = "transactions_tags"
	colID                     = "id"
	colDate                   = "date"
	colTransactionType        = "transaction_type"
	colTransactionDescription = "description"
	colTransactionAmount      = "amount"
	colTransactionID          = "transaction_id"
	colTagID                  = "tag_id"
	colRuleName               = "name"
	colRulPattern             = "pattern"
	colTagValue               = "value"
	colTagRuleID              = "rule_id"

	errUnableToReadRule = "unable to read rule: %w"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	selectTransactionsStmt = psql.Select(
		"id",
		"date",
		"transaction_type",
		"description",
		"amount",
		"tag",
	).
		From(transactionTable).
		LeftJoin("transactions_tags ON transactions_tags.transaction_id = transactions.id")

	selectRulesStmt = psql.Select(
		"id",
		"name",
		"pattern",
		"value",
	).
		From(rulesTable).
		LeftJoin("tags ON tags.rule_id = rules.id")

	selectTagsStmt = psql.Select(
		"value",
		"rule_id",
	).From(tagsTable).
		InnerJoin("rules on rules.id = tags.rule_id")

	selectTransactionTagsStmt = psql.Select("*").From(transactionsTagsTable)

	insertTransaction = psql.Insert(transactionTable).
				Columns(
			colID,
			colDate,
			colTransactionType,
			colTransactionDescription,
			colTransactionAmount,
		)

	insertTransactionTag = psql.Insert(transactionsTagsTable).Columns(colTransactionID, colTagID)
	insertRule           = psql.Insert(rulesTable).Columns(colRuleName, colRulPattern)
	insertTag            = psql.Insert(tagsTable).Columns(colTagValue, colTagRuleID)
)

//go:generate go run github.com/ecordell/optgen -output zz_generated.transaction_filter.go . TransactionFilter
type TransactionFilter struct {
	After  *time.Time `debugmap:"visible"`
	Before *time.Time `debugmap:"visible"`
}

//xgo:generate go run github.com/ecordell/optgen -output zz_generated.rule_filter.go . RuleFilter
type RuleFilter struct {
}

//go:generate go run github.com/ecordell/optgen -output zz_generated.query_transaction_opts.go . QueryTransactionOptions
type QueryTransactionOptions struct {
	Limit  int `debugmap:"visible"`
	Offset int `debugmap:"visible"`
}

//xgo:generate go run github.com/ecordell/optgen -output zz_generated.query_rule_opts.go . QueryRuleOptions
type QueryRuleOptions struct {
}

type TxUserFunc func(context.Context, TransactionWriter, RuleWriter) error

type TransactionWriter interface {
	Write(ctx context.Context, transaction entity.Transaction) error
	Delete(ctx context.Context, id string) error
}

type RuleWriter interface {
	WriteRule(ctx context.Context, rule entity.Rule) error
	WriteTag(ctx context.Context, value string, ruleID int) error
	DeleteRule(ctx context.Context, id string) error
	DeleteTag(ctx context.Context, value string) error
}

type Datastore struct {
	pool *pgxpool.Pool
}

func NewPostgresDatastore(ctx context.Context, url string, options ...Option) (*Datastore, error) {
	pgOptions := newPostgresConfig(options)

	initializationContext, cancelInit := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelInit()

	pgxConfig, err := pgOptions.PgxConfig(url)
	if err != nil {
		return nil, err
	}

	pgPool, err := pgxpool.NewWithConfig(initializationContext, pgxConfig)
	if err != nil {
		return nil, err
	}

	return &Datastore{pool: pgPool}, nil
}

func (d *Datastore) QueryTransactions(ctx context.Context, filter TransactionFilter, opts *QueryTransactionOptions) ([]entity.Transaction, error) {
	query := selectTransactionsStmt
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

func (d *Datastore) WriteTx(ctx context.Context, txFn TxUserFunc) error {
	return nil
}

func (d *Datastore) Close() {
	d.pool.Close()
}
