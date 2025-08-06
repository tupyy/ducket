package pg_test

import (
	"context"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Writer", Ordered, func() {
	var (
		dt     *pg.Datastore
		pgPool *pgxpool.Pool
	)

	BeforeAll(func() {
		pgDt, err := pg.NewPostgresDatastore(context.TODO(), pgUri)
		Expect(err).To(BeNil())
		Expect(pgDt).ToNot(BeNil())

		pgxConfig, err := pgxpool.ParseConfig(pgUri)
		Expect(err).To(BeNil())

		pool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
		Expect(err).To(BeNil())
		Expect(pool).ToNot(BeNil())

		dt = pgDt
		pgPool = pool

		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		dt.Close()
	})

	Context("WriteRule", func() {
		It("write rule successfully -- no labels", func() {
			rule := entity.NewRule("rule1", "pattern1")
			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.WriteRule(context.TODO(), rule, false)
			})
			Expect(err).To(BeNil())

			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from rules;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})

		It("update rule successfully -- no labels", func() {
			sql, args, err := insertRule.
				Values("rule1", "pattern").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rule := entity.NewRule("rule1", "updated_pattern")
			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.WriteRule(context.TODO(), rule, true)
			})
			Expect(err).To(BeNil())

			var pattern string
			err = pgPool.QueryRow(context.TODO(), "select pattern from rules limit 1;").Scan(&pattern)
			Expect(err).To(BeNil())
			Expect(pattern).To(Equal("updated_pattern"))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("DeleteRule", func() {
		It("delete successfully rule", func() {
			sql, args, err := insertRule.
				Values("rule1", "pattern").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteRule(ctx, "rule1")
			})
			Expect(err).To(BeNil())

			var count int
			rows, err := pgPool.Query(context.TODO(), "select count(*) from rules;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(0))
			}
		})

		It("delete successfully rule -- rule not found", func() {
			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteRule(ctx, "rule1")
			})
			Expect(err).To(BeNil())

			var count int
			rows, err := pgPool.Query(context.TODO(), "select count(*) from rules;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(0))
			}
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("WriteTransaction", func() {
		It("write transaction successfully -- without labels", func() {
			t := entity.NewTransaction(entity.CreditTransaction, 1001, time.Now(), 1.1, "content")

			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				_, err := w.WriteTransaction(ctx, *t)
				return err
			})
			Expect(err).To(BeNil())

			var count int
			rows, err := pgPool.Query(context.TODO(), "select count(*) from transactions;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))
			}
		})

		It("update successfully transactions -- no relationship handling", func() {
			sql, args, err := insertRule.
				Values("rule1", "pattern").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertLabelStmt.
				Values("category", "label1").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Get the actual label ID from the database
			var labelID int
			err = pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'label1'").Scan(&labelID)
			Expect(err).To(BeNil())

			// Manually create a relationship to test that WriteTransaction doesn't interfere with it
			sql, args, err = psql.Insert("transactions_labels").
				Columns("transaction_id", "label_id", "rule_id").
				Values(1, labelID, "rule1").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			t := entity.NewTransaction(entity.DebitTransaction, 1001, time.Now(), 2, "content")
			t.ID = 1 // fix ID

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				_, err := w.WriteTransaction(ctx, *t)
				return err
			})
			Expect(err).To(BeNil())

			// Verify the transaction was updated
			content := ""
			kind := ""
			hash := ""
			sql, args, _ = psql.Select("content", "kind", "hash").From("transactions").Where(sq.Eq{"id": 1}).ToSql()
			err = pgPool.QueryRow(context.TODO(), sql, args...).Scan(&content, &kind, &hash)
			Expect(err).To(BeNil())

			Expect(content).To(Equal("content"))
			Expect(kind).To(Equal("debit"))

			// Verify that WriteTransaction doesn't affect relationships - they should still exist
			count := 0
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})

		It("update successfully transactions -- only transaction data", func() {
			sql, args, err := insertTransaction.
				Values(1, time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			t := entity.NewTransaction(entity.DebitTransaction, 1001, time.Now(), 2, "content")
			t.ID = 1 // fix ID

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				_, err := w.WriteTransaction(ctx, *t)
				return err
			})
			Expect(err).To(BeNil())

			content := ""
			account := 0
			kind := ""
			hash := ""
			sql, args, _ = psql.Select("content", "kind", "hash", "account").From("transactions").Where(sq.Eq{"id": 1}).ToSql()
			err = pgPool.QueryRow(context.TODO(), sql, args...).Scan(&content, &kind, &hash, &account)
			Expect(err).To(BeNil())

			Expect(account).To(Equal(1001))
			Expect(content).To(Equal("content"))
			Expect(kind).NotTo(Equal("hash1"))
			Expect(kind).To(Equal("debit"))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("WriteLabel", func() {
		It("write labels successfully", func() {
			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				_, err := w.WriteLabel(ctx, entity.Label{Key: "category", Value: "label1"})
				return err
			})
			Expect(err).To(BeNil())

			var count int
			rows, err := pgPool.Query(context.TODO(), "select count(*) from labels;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))
			}
		})

		It("write labels successfully -- returns correct ID", func() {
			var labelID int
			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				id, err := w.WriteLabel(ctx, entity.Label{Key: "category", Value: "label1"})
				labelID = id
				return err
			})
			Expect(err).To(BeNil())
			Expect(labelID).To(BeNumerically(">", 0))

			var key, value string
			err = pgPool.QueryRow(context.TODO(), "SELECT key, value FROM labels WHERE id = $1", labelID).Scan(&key, &value)
			Expect(err).To(BeNil())
			Expect(key).To(Equal("category"))
			Expect(value).To(Equal("label1"))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("DeleteLabel", func() {
		It("delete labels successfully", func() {
			// First insert a label
			sql, args, err := insertLabelStmt.
				Values("category", "label1").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Verify it exists
			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))

			// Delete the label
			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteLabel(ctx, entity.Label{Key: "category", Value: "label1"})
			})
			Expect(err).To(BeNil())

			// Verify it's deleted
			err = pgPool.QueryRow(context.TODO(), "select count(*) from labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("delete labels successfully -- label not found", func() {
			// Try to delete a non-existent label
			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteLabel(ctx, entity.Label{Key: "category", Value: "nonexistent"})
			})
			Expect(err).To(BeNil())

			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("DeleteTransaction", func() {
		It("delete transaction successfully", func() {
			// First insert a transaction
			sql, args, err := insertTransaction.
				Values(1, time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Verify it exists
			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))

			// Delete the transaction
			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteTransaction(ctx, 1)
			})
			Expect(err).To(BeNil())

			// Verify it's deleted
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("delete transaction successfully -- transaction not found", func() {
			// Try to delete a non-existent transaction
			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteTransaction(ctx, 999)
			})
			Expect(err).To(BeNil())

			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("WriteRelationships", func() {
		BeforeEach(func() {
			// Setup some test data
			sql, args, err := insertRule.
				Values("rule1", "pattern1").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertLabelStmt.
				Values("category", "food").
				Values("category", "transport").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "debit", "grocery store", "50.0", "hash1").
				Values(2, time.Now(), 1002, "debit", "bus ticket", "5.0", "hash2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())
		})

		It("write label-rule relationships successfully", func() {
			// Get label IDs
			var label1ID, label2ID int
			err := pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'food'").Scan(&label1ID)
			Expect(err).To(BeNil())
			err = pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'transport'").Scan(&label2ID)
			Expect(err).To(BeNil())

			relationships := []entity.Relationship{
				{Kind: entity.RelationshipLabelRule, LabelID: label1ID, RuleID: "rule1"},
				{Kind: entity.RelationshipLabelRule, LabelID: label2ID, RuleID: "rule2"},
			}

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.WriteRelationships(ctx, relationships)
			})
			Expect(err).To(BeNil())

			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from rules_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(2))
		})

		It("write transaction-label relationships successfully", func() {
			// Get label ID
			var labelID int
			err := pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'food'").Scan(&labelID)
			Expect(err).To(BeNil())

			transactionID := 1
			relationships := []entity.Relationship{
				{Kind: entity.RelationshipLabelTransaction, LabelID: labelID, TransactionID: transactionID},
			}

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.WriteRelationships(ctx, relationships)
			})
			Expect(err).To(BeNil())

			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})

		It("write transaction-label-rule relationships successfully", func() {
			// Get label ID
			var labelID int
			err := pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'food'").Scan(&labelID)
			Expect(err).To(BeNil())

			transactionID := 1
			relationships := []entity.Relationship{
				{Kind: entity.RelationshipLabelRuleTransaction, LabelID: labelID, RuleID: "rule1", TransactionID: transactionID},
			}

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.WriteRelationships(ctx, relationships)
			})
			Expect(err).To(BeNil())

			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("DeleteRelationships", func() {
		BeforeEach(func() {
			// Setup some test data
			sql, args, err := insertRule.
				Values("rule1", "pattern1").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertLabelStmt.
				Values("category", "food").
				Values("category", "transport").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "debit", "grocery store", "50.0", "hash1").
				Values(2, time.Now(), 1002, "debit", "bus ticket", "5.0", "hash2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())
		})

		It("delete label-rule relationships successfully", func() {
			// Get label IDs
			var label1ID, label2ID int
			err := pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'food'").Scan(&label1ID)
			Expect(err).To(BeNil())
			err = pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'transport'").Scan(&label2ID)
			Expect(err).To(BeNil())

			// First create some relationships
			sql, args, err := psql.Insert("rules_labels").
				Columns("rule_id", "label_id").
				Values("rule1", label1ID).
				Values("rule2", label2ID).
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Verify they exist
			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from rules_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(2))

			// Delete one relationship
			relationships := []entity.Relationship{
				{Kind: entity.RelationshipLabelRule, LabelID: label1ID, RuleID: "rule1"},
			}

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteRelationships(ctx, relationships)
			})
			Expect(err).To(BeNil())

			// Verify only one is left
			err = pgPool.QueryRow(context.TODO(), "select count(*) from rules_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})

		It("delete transaction-label relationships successfully", func() {
			// Get label ID
			var labelID int
			err := pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'food'").Scan(&labelID)
			Expect(err).To(BeNil())

			// First create a relationship
			sql, args, err := insertTransactionLabel.
				Values(1, labelID, "rule1").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Verify it exists
			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))

			// Delete the relationship
			transactionID := 1
			relationships := []entity.Relationship{
				{Kind: entity.RelationshipLabelTransaction, LabelID: labelID, TransactionID: transactionID},
			}

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteRelationships(ctx, relationships)
			})
			Expect(err).To(BeNil())

			// Verify it's deleted
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("delete transaction-label-rule relationships successfully", func() {
			// Get label ID
			var labelID int
			err := pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'food'").Scan(&labelID)
			Expect(err).To(BeNil())

			// First create a relationship
			sql, args, err := insertTransactionLabel.
				Values(1, labelID, "rule1").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Verify it exists
			var count int
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))

			// Delete the relationship
			transactionID := 1
			relationships := []entity.Relationship{
				{Kind: entity.RelationshipLabelRuleTransaction, LabelID: labelID, RuleID: "rule1", TransactionID: transactionID},
			}

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteRelationships(ctx, relationships)
			})
			Expect(err).To(BeNil())

			// Verify it's deleted
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("delete label-rule relationships cleans up both rules_labels and transactions_labels tables", func() {
			// Get label IDs
			var labelID int
			err := pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'food'").Scan(&labelID)
			Expect(err).To(BeNil())

			// Create a rule-label relationship
			sql, args, err := psql.Insert("rules_labels").
				Columns("rule_id", "label_id").
				Values("rule1", labelID).
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction-label relationships that were created by the rule
			sql, args, err = insertTransactionLabel.
				Values(1, labelID, "rule1").
				Values(2, labelID, "rule1").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Verify both relationships exist
			var rulesLabelsCount, transactionsLabelsCount int
			err = pgPool.QueryRow(context.TODO(), "SELECT count(*) FROM rules_labels WHERE rule_id = 'rule1' AND label_id = $1", labelID).Scan(&rulesLabelsCount)
			Expect(err).To(BeNil())
			Expect(rulesLabelsCount).To(Equal(1))

			err = pgPool.QueryRow(context.TODO(), "SELECT count(*) FROM transactions_labels WHERE rule_id = 'rule1' AND label_id = $1", labelID).Scan(&transactionsLabelsCount)
			Expect(err).To(BeNil())
			Expect(transactionsLabelsCount).To(Equal(2))

			// Delete the rule-label relationship - this should clean up both tables
			relationships := []entity.Relationship{
				{Kind: entity.RelationshipLabelRule, LabelID: labelID, RuleID: "rule1"},
			}

			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w *pg.Writer) error {
				return w.DeleteRelationships(ctx, relationships)
			})
			Expect(err).To(BeNil())

			// Verify both tables are cleaned up
			err = pgPool.QueryRow(context.TODO(), "SELECT count(*) FROM rules_labels WHERE rule_id = 'rule1' AND label_id = $1", labelID).Scan(&rulesLabelsCount)
			Expect(err).To(BeNil())
			Expect(rulesLabelsCount).To(Equal(0))

			err = pgPool.QueryRow(context.TODO(), "SELECT count(*) FROM transactions_labels WHERE rule_id = 'rule1' AND label_id = $1", labelID).Scan(&transactionsLabelsCount)
			Expect(err).To(BeNil())
			Expect(transactionsLabelsCount).To(Equal(0))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})
})
