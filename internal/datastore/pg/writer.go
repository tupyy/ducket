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
	labelsToDissociate := []int{}
	labelsToAssociate := []entity.Label{}
	for _, label := range rule.Labels {
		labelsToAssociate = append(labelsToAssociate, label)
	}

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
			labelsToDissociate = append(labelsToDissociate, *oldRule.LabelID)
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

	if len(labelsToDissociate) > 0 {
		deleteLabelAssociationStmt := psql.Delete(rulesLabelsTable).Where(sq.Eq{colRuleID: rule.Name})

		sql, arg, err := deleteLabelAssociationStmt.ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteLabel, err)
		}
		if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
			return fmt.Errorf(errUnableToWriteLabel, err)
		}
	}

	if len(labelsToAssociate) > 0 {
		addLabelAssociationStmt := psql.Insert(rulesLabelsTable).Columns(colRuleID, colLabelID)
		for _, label := range labelsToAssociate {
			// First, make sure the label exists (create if it doesn't)
			if err := w.WriteLabel(ctx, label); err != nil {
				return fmt.Errorf(errUnableToWriteLabel, err)
			}

			// Get the label ID from the database
			var labelID int
			labelSelectSQL, labelSelectArgs, err := psql.Select(colID).
				From(labelsTable).
				Where(sq.And{
					sq.Eq{colLabelKey: label.Key},
					sq.Eq{colLabelValue: label.Value},
				}).
				ToSql()
			if err != nil {
				return fmt.Errorf(errUnableToWriteLabel, err)
			}

			err = w.tx.QueryRow(ctx, labelSelectSQL, labelSelectArgs...).Scan(&labelID)
			if err != nil {
				return fmt.Errorf(errUnableToWriteLabel, err)
			}

			addLabelAssociationStmt = addLabelAssociationStmt.Values(rule.Name, labelID)
		}

		// add on conflict statement
		addLabelAssociationStmt = addLabelAssociationStmt.Suffix("ON CONFLICT DO NOTHING")

		sql, arg, err := addLabelAssociationStmt.ToSql()
		if err != nil {
			return fmt.Errorf(errUnableToWriteLabel, err)
		}
		if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
			return fmt.Errorf(errUnableToWriteLabel, err)
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

func (w *writerTx) WriteLabel(ctx context.Context, label entity.Label) error {
	sql, args, err := psql.Insert(labelsTable).
		Columns(colLabelKey, colLabelValue).
		Values(label.Key, label.Value).
		Suffix("ON CONFLICT (key, value) DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToWriteLabel, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToWriteLabel, err)
	}

	return nil
}

func (w *writerTx) DeleteLabel(ctx context.Context, label entity.Label) error {
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

func (w *writerTx) WriteTransaction(ctx context.Context, transaction entity.Transaction) (int64, error) {
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

	// remove old label associations
	sql, arg, err := psql.Delete(transactionsLabelsTable).Where(sq.Eq{colTransactionID: transactionID}).ToSql()
	if err != nil {
		return 0, fmt.Errorf(errUnableToWriteLabel, err)
	}
	if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
		return 0, fmt.Errorf(errUnableToWriteLabel, err)
	}

	if len(transaction.Labels) > 0 {
		addLabelStmt := psql.Insert(transactionsLabelsTable).Columns(colTransactionID, colLabelID, colRuleID)
		for labelID, ruleID := range transaction.Labels {
			addLabelStmt = addLabelStmt.Values(transactionID, labelID, ruleID)
		}

		sql, arg, err := addLabelStmt.ToSql()
		if err != nil {
			return 0, fmt.Errorf(errUnableToWriteLabel, err)
		}
		if _, err := w.tx.Exec(ctx, sql, arg...); err != nil {
			return 0, fmt.Errorf(errUnableToWriteLabel, err)
		}
	}

	return int64(transactionID), nil

}

func (w *writerTx) DeleteTransaction(ctx context.Context, id int64) error {
	sql, args, err := psql.Delete(transactionTable).Where(sq.Eq{colID: id}).ToSql()
	if err != nil {
		return fmt.Errorf(errUnableToDeleteTransaction, err)
	}

	if _, err := w.tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf(errUnableToDeleteTransaction, err)
	}

	return nil
}
