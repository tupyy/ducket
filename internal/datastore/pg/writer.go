package pg

import (
	"context"
	"errors"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	errUnableToWriteTag          = "unable to write tag: %w"
	errUnableToDeleteTag         = "unable to delete tag: %w"
	errUnableToWriteRule         = "unable to write rule: %w"
	errUnableToDeleteTransaction = "unable to delete transaction: %w"
	rulesTagsTable               = "rules_tags"
)

var (
	insertRule = psql.Insert("rules").Columns("id", "name", "pattern")
	updateRule = psql.Update("rules")
	insertTag  = psql.Insert("tags").Columns("value")
)

type writerTx struct {
	tx pgx.Tx
}

func (w *writerTx) WriteRule(ctx context.Context, rule entity.Rule, update bool) error {
	var (
		sql  string
		args []any
		err  error
	)

	switch update {
	case false:
		query := insertRule
		sql, args, err = query.Values(rule.ID, rule.Name, rule.Pattern).ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}
	case true:
		query := updateRule
		sql, args, err = query.
			Set("id", rule.ID).
			Set("name", rule.Name).
			Set("pattern", rule.Pattern).
			Where(sq.Eq{"id": rule.ID}).
			ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToWriteRule, err)
	}

	// associate rules and tags if any
	if len(rule.Tags) > 0 {
		queryTag := psql.Insert(rulesTagsTable).Columns(colRuleID, colTag)
		for _, tag := range rule.Tags {
			queryTag = queryTag.Values(rule.ID, tag.Value)
		}

		// add on conflict statement
		queryTag = queryTag.Suffix("ON CONFLICT DO NOTHING")

		sql, arg, err := queryTag.ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
		if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
	}

	return nil
}

func (w *writerTx) DeleteRule(ctx context.Context, id string) error {
	return nil
}

func (w *writerTx) WriteTag(ctx context.Context, value string) error {
	sql, args, err := psql.Insert(tagsTable).Columns(colValue).Values(value).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToWriteTag, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToWriteTag, err)
	}

	return nil
}

func (w *writerTx) DeleteTag(ctx context.Context, value string) error {
	sql, args, err := psql.Delete(tagsTable).Where(sq.Eq{colValue: value}).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToDeleteTag, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToDeleteTag, err)
	}

	return nil
}

func (w *writerTx) WriteTransaction(ctx context.Context, transaction entity.Transaction) error {
	return errors.New("not implementated")
}

func (w *writerTx) DeleteTransaction(ctx context.Context, id string) error {
	sql, args, err := psql.Delete(transactionTable).Where(sq.Eq{colID: id}).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToDeleteTransaction, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToDeleteTransaction, err)
	}

	return nil
}
