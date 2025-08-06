package pg

import (
	"context"
	"fmt"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg/models"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Writer struct {
	tx pgx.Tx
}

func (w *Writer) WriteRule(ctx context.Context, rule entity.Rule, update bool) error {
	switch update {
	case false:
		query := insertRule
		sql, args, err := query.Values(rule.Name, rule.Pattern).ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}

		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}
	case true:
		// get the old one
		sql, args, err := selectRulesStmt.Where(sq.Eq{"rule_id": rule.Name}).ToSql()
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
			if oldRule.LabelID == nil {
				continue
			}
		}
		rows.Close()

		query := updateRule
		sql, args, err = query.
			Set("pattern", rule.Pattern).
			Where(sq.Eq{"id": rule.Name}).
			ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}

		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf(errUnableToWriteRule, err)
		}
	}

	return nil
}

func (w *Writer) DeleteRule(ctx context.Context, id string) error {
	sql, args, err := psql.Delete(rulesTable).Where(sq.Eq{colID: id}).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToDeleteRule, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToDeleteRule, err)
	}

	return nil
}

func (w *Writer) WriteLabel(ctx context.Context, label entity.Label) (int, error) {
	sql, args, err := psql.Insert(labelsTable).
		Columns(colLabelKey, colLabelValue).
		Values(label.Key, label.Value).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf(errUnableToWriteLabel, err)
	}

	id := 0
	if err := w.tx.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, fmt.Errorf(errUnableToWriteLabel, err)
	}

	return id, nil
}

func (w *Writer) DeleteLabel(ctx context.Context, label entity.Label) error {
	sql, args, err := psql.Delete(labelsTable).Where(sq.And{
		sq.Eq{colLabelKey: label.Key},
		sq.Eq{colLabelValue: label.Value},
	}).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToDeleteLabel, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToDeleteLabel, err)
	}

	return nil
}

func (w *Writer) WriteTransaction(ctx context.Context, transaction entity.Transaction) (int64, error) {
	transactionID := transaction.ID
	stmt := selectTransactionStmp.Where(sq.Eq{colID: transactionID})

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, fmt.Errorf(errUnableToWriteTransaction, err)
	}

	rows, err := w.tx.Query(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf(errUnableToWriteTransaction, err)
	}

	update := false
	rs := pgxscan.NewRowScanner(rows)
	for rows.Next() {
		t := models.Transaction{}
		err := rs.Scan(&t)
		if err != nil {
			return 0, fmt.Errorf(errUnableToWriteTransaction, err)
		}
		update = true
	}
	rows.Close()

	switch update {
	case false:
		// insert transaction
		sql, args, err = insertTransaction.Values(
			transaction.Date,
			transaction.Account,
			transaction.Hash,
			transaction.Kind,
			transaction.RawContent,
			transaction.Amount,
			transaction.Info,
		).Suffix("RETURNING id").ToSql()
		if err != nil {
			return 0, fmt.Errorf(errUnableToWriteTransaction, err)
		}

		rows, err := w.tx.Query(ctx, sql, args...)
		if err != nil {
			return 0, fmt.Errorf(errUnableToWriteTransaction, err)
		}

		if rows.Next() {
			if err := rows.Scan(&transactionID); err != nil {
				return 0, fmt.Errorf(errUnableToWriteTransaction, err)
			}
		}
		rows.Close()
	case true:
		sql, args, err = psql.Update(transactionTable).
			Set(colDate, transaction.Date).
			Set(colTransactionAccount, transaction.Account).
			Set(colHash, transaction.Hash).
			Set(colTransactionContent, transaction.RawContent).
			Set(colTransactionAmount, transaction.Amount).
			Set(colTransactionType, transaction.Kind).
			Set(colTransactionInfo, transaction.Info).
			Where(sq.Eq{colID: transactionID}).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf(errUnableToWriteTransaction, err)
		}

		_, err := w.tx.Exec(ctx, sql, args...)
		if err != nil {
			return 0, fmt.Errorf(errUnableToWriteTransaction, err)
		}
	}

	return int64(transactionID), nil
}

func (w *Writer) DeleteTransaction(ctx context.Context, id int64) error {
	sql, args, err := psql.Delete(transactionTable).Where(sq.Eq{colID: id}).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToDeleteTransaction, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToDeleteTransaction, err)
	}

	return nil
}

func (w *Writer) WriteRelationships(ctx context.Context, relationships []entity.Relationship) error {
	qInsertTransactionLabelRule := insertTransactionLabelRuleRelationship
	hasTransactionLabelRuleRelationship := false

	qInsertTransactionLabel := insertTransactionLabelRelationship
	hasTransactionLabelRelationship := false

	qInsertLabelRule := insertLabelRuleRelationship
	hasLabelRuleRelationship := false

	for _, r := range relationships {
		switch r.Kind {
		case entity.RelationshipLabelRule:
			qInsertLabelRule = qInsertLabelRule.Values(r.LabelID, r.RuleID)
			hasLabelRuleRelationship = true
		case entity.RelationshipLabelTransaction:
			qInsertTransactionLabel = qInsertTransactionLabel.Values(r.TransactionID, r.LabelID)
			hasTransactionLabelRelationship = true
		case entity.RelationshipLabelRuleTransaction:
			qInsertTransactionLabelRule = qInsertTransactionLabelRule.Values(r.TransactionID, r.LabelID, r.RuleID)
			hasTransactionLabelRuleRelationship = true
		}
	}

	if hasTransactionLabelRuleRelationship {
		sql, args, err := qInsertTransactionLabelRule.Suffix("ON CONFLICT DO NOTHING").ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRelationship, err)
		}
		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf(errUnableToWriteRelationship, err)
		}
	}

	if hasTransactionLabelRelationship {
		sql, args, err := qInsertTransactionLabel.Suffix("ON CONFLICT DO NOTHING").ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRelationship, err)
		}
		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf(errUnableToWriteRelationship, err)
		}
	}

	if hasLabelRuleRelationship {
		sql, args, err := qInsertLabelRule.Suffix("ON CONFLICT DO NOTHING").ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteRelationship, err)
		}
		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf(errUnableToWriteRelationship, err)
		}
	}

	return nil
}

func (w *Writer) DeleteRelationships(ctx context.Context, relationships []entity.Relationship) error {
	execQuery := func(q squirrel.DeleteBuilder) error {
		sql, args, err := q.ToSql()
		if err != nil {
			return err
		}
		if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
			return err
		}
		return nil
	}

	for _, r := range relationships {
		switch r.Kind {
		case entity.RelationshipLabelRule:
			queries := []sq.DeleteBuilder{}
			queries = append(queries, deleteLabelRuleRelationship.Where(sq.And{sq.Eq{colRuleID: r.RuleID}, sq.Eq{colLabelID: r.LabelID}}))
			queries = append(queries, deleteTransactionLabelRuleRelationship.Where(sq.And{sq.Eq{colRuleID: r.RuleID}, sq.Eq{colLabelID: r.LabelID}}))

			for _, q := range queries {
				if err := execQuery(q); err != nil {
					return fmt.Errorf(errUnableToDeleteRelationship, err)
				}
			}
		case entity.RelationshipLabelTransaction:
			query := deleteTransactionLabelRuleRelationship.Where(sq.And{sq.Eq{colTransactionID: r.TransactionID}, sq.Eq{colLabelID: r.LabelID}})
			if err := execQuery(query); err != nil {
				return fmt.Errorf(errUnableToDeleteRelationship, err)
			}
		case entity.RelationshipLabelRuleTransaction:
			query := deleteTransactionLabelRuleRelationship.Where(sq.And{sq.Eq{colTransactionID: r.TransactionID}, sq.Eq{colLabelID: r.LabelID}, sq.Eq{colRuleID: r.RuleID}})
			if err := execQuery(query); err != nil {
				return fmt.Errorf(errUnableToDeleteRelationship, err)
			}
		}
	}

	return nil
}
