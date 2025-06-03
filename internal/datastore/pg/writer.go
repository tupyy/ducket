package pg

import (
	"context"
	"errors"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg/models"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

const (
	errUnableToWriteTag          = "unable to write tag: %w"
	errUnableToDeleteTag         = "unable to delete tag: %w"
	errUnableToDeleteRule        = "unable to delete rule: %w"
	errUnableToWriteRule         = "unable to write rule: %w"
	errUnableToDeleteTransaction = "unable to delete transaction: %w"
	rulesTagsTable               = "rules_tags"
)

var (
	insertRule = psql.Insert(rulesTable).Columns("id", "name", "pattern")
	updateRule = psql.Update(rulesTable)
	insertTag  = psql.Insert("tags").Columns("value")
)

type writerTx struct {
	tx pgx.Tx
}

func (w *writerTx) WriteRule(ctx context.Context, rule entity.Rule, update bool) error {
	tagsToDissociate := []string{}
	tagsToAssociate := []string{}
	for _, tag := range rule.Tags {
		tagsToAssociate = append(tagsToAssociate, tag.Value)
	}

	switch update {
	case false:
		query := insertRule
		sql, args, err := query.Values(rule.ID, rule.Name, rule.Pattern).ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}

		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}
	case true:
		// get the old one
		sql, args, err := selectRulesStmt.Where(sq.Eq{"id": rule.ID}).ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}

		rows, err := w.tx.Query(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}

		rs := pgxscan.NewRowScanner(rows)
		for rows.Next() {
			oldRule := models.RuleRow{}
			err := rs.Scan(&oldRule)
			if err != nil {
				return fmt.Errorf(errUnableToWriteRule, err)
			}
			if oldRule.TagValue == nil {
				continue
			}
			tagsToDissociate = append(tagsToDissociate, *oldRule.TagValue)
		}
		rows.Close()

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

		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}
	}

	if len(tagsToDissociate) > 0 {
		deleteTagAssociationStmt := psql.Delete(rulesTagsTable)

		or := sq.Or{}
		for _, tag := range tagsToDissociate {
			or = append(or, sq.Eq{colTag: tag})
		}
		deleteTagAssociationStmt.Where(or)

		sql, arg, err := deleteTagAssociationStmt.ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
		if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
	}

	if len(tagsToAssociate) > 0 {
		addTagAssociationStmt := psql.Insert(rulesTagsTable).Columns(colRuleID, colTag)
		for _, tag := range tagsToAssociate {
			addTagAssociationStmt = addTagAssociationStmt.Values(rule.ID, tag)
		}

		// add on conflict statement
		addTagAssociationStmt = addTagAssociationStmt.Suffix("ON CONFLICT DO NOTHING")

		sql, arg, err := addTagAssociationStmt.ToSql()
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
	sql, args, err := psql.Delete(rulesTable).Where(sq.Eq{colID: id}).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToDeleteRule, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToDeleteRule, err)
	}

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
