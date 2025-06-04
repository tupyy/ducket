package pg

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg/models"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type writerTx struct {
	tx pgx.Tx
}

func (w *writerTx) WriteRule(ctx context.Context, rule entity.Rule, update bool) error {
	tagsToDissociate := []string{}
	tagsToAssociate := []string{}
	for _, tag := range rule.Tags {
		tagsToAssociate = append(tagsToAssociate, tag)
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
			if oldRule.Tag == nil {
				continue
			}
			tagsToDissociate = append(tagsToDissociate, *oldRule.Tag)
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
	transactionID := transaction.ID
	stmt := selectTransactionStmp.Where(sq.Eq{colID: transactionID})

	sql, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToWriteTransaction, err)
	}

	rows, err := w.tx.Query(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf(errUnableToWriteTransaction, err)
	}

	tagsToDissociate := []string{}
	update := false
	rs := pgxscan.NewRowScanner(rows)
	for rows.Next() {
		t := models.Transaction{}
		err := rs.Scan(&t)
		if err != nil {
			return fmt.Errorf(errUnableToWriteTransaction, err)
		}
		if t.Tag != nil && t.RuleID != nil {
			tagsToDissociate = append(tagsToDissociate, *t.Tag, *t.RuleID)
		}
		update = true
	}
	rows.Close()

	switch update {
	case false:
		// insert transaction
		sql, args, err = insertTransaction.Values(
			transaction.Date,
			transaction.Hash,
			transaction.Kind,
			transaction.RawContent,
			transaction.Amount,
		).Suffix("RETURNING id").ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteTransaction, err)
		}

		rows, err := w.tx.Query(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf(errUnableToWriteTransaction, err)
		}

		if rows.Next() {
			if err := rows.Scan(&transactionID); err != nil {
				return fmt.Errorf(errUnableToWriteTransaction, err)
			}
		}
		rows.Close()
	case true:
		sql, args, err = psql.Update(transactionTable).
			Set(colDate, transaction.Date).
			Set(colHash, transaction.Hash).
			Set(colTransactionContent, transaction.RawContent).
			Set(colAmount, transaction.Amount).
			Set(colTransactionType, transaction.Kind).
			Where(sq.Eq{colID: transactionID}).
			ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteTransaction, err)
		}

		_, err := w.tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf(errUnableToWriteTransaction, err)
		}
	}

	// remove old tags associations
	if len(tagsToDissociate) > 0 {
		sql, arg, err := psql.Delete(transactionsTagsTable).Where(sq.Eq{colTransactionID: transactionID}).ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
		if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
	}

	if len(transaction.Tags) > 0 {
		addTagStmt := psql.Insert(transactionsTagsTable).Columns(colTransactionID, colTagID, colRuleID)
		for tag, ruleID := range transaction.Tags {
			addTagStmt = addTagStmt.Values(transactionID, tag, ruleID)
		}

		sql, arg, err := addTagStmt.ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
		if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
			return fmt.Errorf(errUnableToWriteTag, err)
		}
	}

	return nil

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
