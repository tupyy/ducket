package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	pkgErrors "git.tls.tupangiu.ro/cosmin/finante/pkg/errors"
	sq "github.com/Masterminds/squirrel"
)

const NoFilter = ""

var txnBaseColumns = []string{
	"t.id", "t.hash", "t.date", "t.account", "t.kind",
	"t.amount", "t.content", "t.info", "t.recipient",
}

func (s *Store) ListTransactions(ctx context.Context, filter string, limit, offset int) ([]entity.Transaction, error) {
	q, tagCTE, tagArgs, err := s.taggedQuery(ctx)
	if err != nil {
		return nil, err
	}

	q = q.OrderBy("t.date DESC", "t.id DESC")

	if filter != "" {
		sqlizer, err := ParseTransactionFilter(filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter: %w", err)
		}
		q = q.Where(sqlizer)
	}
	if limit > 0 {
		q = q.Limit(uint64(limit))
	}
	if offset > 0 {
		q = q.Offset(uint64(offset))
	}

	rows, err := s.execTaggedQuery(ctx, q, tagCTE, tagArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTransactions(rows)
}

func (s *Store) GetTransaction(ctx context.Context, id int64) (*entity.Transaction, error) {
	q, tagCTE, tagArgs, err := s.taggedQuery(ctx)
	if err != nil {
		return nil, err
	}

	q = q.Where(sq.Eq{"t.id": id})

	rows, err := s.execTaggedQuery(ctx, q, tagCTE, tagArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txns, err := scanTransactions(rows)
	if err != nil {
		return nil, err
	}
	if len(txns) == 0 {
		return nil, pkgErrors.NewResourceNotFoundError("transaction", fmt.Sprintf("%d", id))
	}
	return &txns[0], nil
}

func (s *Store) GetTransactionByHash(ctx context.Context, hash string) (*entity.Transaction, error) {
	q, tagCTE, tagArgs, err := s.taggedQuery(ctx)
	if err != nil {
		return nil, err
	}

	q = q.Where(sq.Eq{"t.hash": hash})

	rows, err := s.execTaggedQuery(ctx, q, tagCTE, tagArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txns, err := scanTransactions(rows)
	if err != nil {
		return nil, err
	}
	if len(txns) == 0 {
		return nil, nil
	}
	return &txns[0], nil
}

func (s *Store) CreateTransaction(ctx context.Context, t entity.Transaction) (int64, error) {
	query, args, err := sq.Insert("transactions").
		Columns("hash", "date", "account", "kind", "amount", "content", "info", "recipient").
		Values(t.Hash, t.Date, t.Account, string(t.Kind), t.Amount, t.Content, t.Info, t.Recipient).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = s.qi.QueryRowContext(ctx, query, args...).Scan(&id)
	return id, err
}

func (s *Store) UpdateTransaction(ctx context.Context, t entity.Transaction) error {
	query, args, err := sq.Update("transactions").
		Set("hash", t.Hash).
		Set("date", t.Date).
		Set("account", t.Account).
		Set("kind", string(t.Kind)).
		Set("amount", t.Amount).
		Set("content", t.Content).
		Set("info", t.Info).
		Set("recipient", t.Recipient).
		Where(sq.Eq{"id": t.ID}).
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return err
	}

	result, err := s.qi.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return pkgErrors.NewResourceNotFoundError("transaction", fmt.Sprintf("%d", t.ID))
	}
	return nil
}

func (s *Store) DeleteTransaction(ctx context.Context, id int64) error {
	query, args, err := sq.Delete("transactions").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return err
	}

	result, err := s.qi.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return pkgErrors.NewResourceNotFoundError("transaction", fmt.Sprintf("%d", id))
	}
	return nil
}

// buildTagCTE builds a CTE that matches all rules against transactions using the filter DSL.
func buildTagCTE(rules []entity.Rule) (string, []any) {
	if len(rules) == 0 {
		return "", nil
	}

	var branches []string
	var allArgs []any

	for _, r := range rules {
		sqlizer, err := ParseTransactionFilter(r.Filter)
		if err != nil {
			continue
		}

		whereSQL, whereArgs, err := sqlizer.ToSql()
		if err != nil {
			continue
		}

		tagsLiteral := duckdbArray(r.Tags)
		branch := fmt.Sprintf(
			"SELECT t.id AS transaction_id, %s AS tags FROM transactions t WHERE %s",
			tagsLiteral, whereSQL,
		)
		branches = append(branches, branch)
		allArgs = append(allArgs, whereArgs...)
	}

	if len(branches) == 0 {
		return "", nil
	}

	cte := "rule_matches AS (" + strings.Join(branches, " UNION ALL ") + "), " +
		"rule_tags AS (SELECT transaction_id, list_distinct(flatten(list(tags))) AS tags FROM rule_matches GROUP BY transaction_id)"

	return cte, allArgs
}

// taggedQuery builds a SELECT with optional tag CTE. When rules exist, it joins
// the dynamically-built rule_tags CTE. When no rules exist, it selects an empty
// array literal for tags and skips the CTE entirely.
func (s *Store) taggedQuery(ctx context.Context) (sq.SelectBuilder, string, []any, error) {
	rules, err := s.ListRules(ctx, NoFilter, 0, 0)
	if err != nil {
		return sq.SelectBuilder{}, "", nil, err
	}

	tagCTE, tagArgs := buildTagCTE(rules)

	cols := make([]string, len(txnBaseColumns))
	copy(cols, txnBaseColumns)

	q := sq.Select().From("transactions t").PlaceholderFormat(sq.Question)

	if tagCTE != "" {
		cols = append(cols, "COALESCE(rt.tags, []::VARCHAR[]) AS tags")
		q = q.Columns(cols...).LeftJoin("rule_tags rt ON rt.transaction_id = t.id")
	} else {
		cols = append(cols, "[]::VARCHAR[] AS tags")
		q = q.Columns(cols...)
	}

	return q, tagCTE, tagArgs, nil
}

func (s *Store) execTaggedQuery(ctx context.Context, q sq.SelectBuilder, tagCTE string, tagArgs []any) (*sql.Rows, error) {
	query, selectArgs, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	if tagCTE != "" {
		query = "WITH " + tagCTE + " " + query
	}

	allArgs := make([]any, 0, len(tagArgs)+len(selectArgs))
	allArgs = append(allArgs, tagArgs...)
	allArgs = append(allArgs, selectArgs...)
	return s.qi.QueryContext(ctx, query, allArgs...)
}

func scanTransactions(rows *sql.Rows) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		var info, recipient sql.NullString
		var tags StringArray
		if err := rows.Scan(
			&t.ID, &t.Hash, &t.Date, &t.Account, &t.Kind,
			&t.Amount, &t.Content, &info, &recipient, &tags,
		); err != nil {
			return nil, err
		}
		if info.Valid {
			t.Info = &info.String
		}
		if recipient.Valid {
			t.Recipient = &recipient.String
		}
		t.Tags = []string(tags)
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}
