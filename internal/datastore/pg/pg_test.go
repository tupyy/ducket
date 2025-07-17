package pg_test

import (
	"context"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	pgUri           = "postgresql://postgres:postgres@localhost:5432/postgres"
	insertRulesStmt = "INSERT INTO rules VALUES ('%s', '%s', '%s');"
)

var (
	psql                   = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertRule             = psql.Insert("rules").Columns("id", "pattern")
	insertLabelStmt        = psql.Insert("labels").Columns("key", "value")
	insertLabelWithIDStmt  = psql.Insert("labels").Columns("id", "key", "value")
	insertTransaction      = psql.Insert("transactions").Columns("id", "date", "account", "kind", "content", "amount", "hash")
	insertTransactionLabel = psql.Insert("transactions_labels").Columns("transaction_id", "label_id", "rule_id")
)

var _ = Describe("query", Ordered, func() {
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

	Context("rule", func() {
		It("reads successfully rules -- empty response", func() {
			rules, err := dt.QueryRules(context.TODO())
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(0))
		})

		It("reads successfully rules -- without labels", func() {
			sql, args, err := insertRule.
				Values("rule1", "pattern").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rules, err := dt.QueryRules(context.TODO())
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(2))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[0].Name))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[1].Name))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[0].Pattern))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[1].Pattern))
		})

		It("reads successfully rules -- with labels", func() {
			sql, args, err := insertRule.
				Values("rule1", "pattern").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertLabelStmt.
				Values("category", "label1").
				Values("category", "label2").
				Values("category", "label3").
				Values("category", "label4").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Get the actual label IDs from the database
			var labelIDs []int
			rows, err := pgPool.Query(context.TODO(), "SELECT id FROM labels ORDER BY id")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				var labelID int
				err := rows.Scan(&labelID)
				Expect(err).To(BeNil())
				labelIDs = append(labelIDs, labelID)
			}
			Expect(labelIDs).To(HaveLen(4))

			sql, args, err = psql.Insert("rules_labels").Columns("rule_id", "label_id").
				Values("rule1", labelIDs[0]).
				Values("rule1", labelIDs[1]).
				Values("rule2", labelIDs[2]).
				Values("rule2", labelIDs[3]).
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rules, err := dt.QueryRules(context.TODO())
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(2))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[0].Name))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[1].Name))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[0].Pattern))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[1].Pattern))

			Expect(rules[0].Labels).To(HaveLen(2))
			Expect(rules[1].Labels).To(HaveLen(2))
			for _, label := range rules[0].Labels {
				Expect([]string{"label1", "label2", "label3", "label4"}).To(ContainElement(label.Value))
			}
			for _, label := range rules[1].Labels {
				Expect([]string{"label1", "label2", "label3", "label4"}).To(ContainElement(label.Value))
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

	Context("transaction", func() {
		It("query successfully transactions -- empty response", func() {
			transactions, err := dt.QueryTransactions(context.TODO())
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(0))
		})

		It("query successfully transactions -- no labels", func() {
			sql, args, err := insertTransaction.
				Values(1, time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				Values(2, time.Now(), 1002, "debit", "transaction", "2.1", "hash2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO())
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(2))

			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[0].Kind)))
			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[1].Kind)))
		})

		It("query successfully transactions -- with labels", func() {
			sql, args, err := insertRule.
				Values("rule1", "pattern").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertLabelStmt.
				Values("category", "label1").
				Values("category", "label2").
				Values("category", "label3").
				Values("category", "label4").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				Values(2, time.Now(), 1002, "debit", "transaction", "2.1", "hash2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Get the actual label IDs from the database
			var label1ID, label2ID int
			err = pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'label1'").Scan(&label1ID)
			Expect(err).To(BeNil())
			err = pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'label2'").Scan(&label2ID)
			Expect(err).To(BeNil())

			sql, args, err = insertTransactionLabel.
				Values(1, label1ID, "rule1").
				Values(1, label2ID, "rule2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO())
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(2))

			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[0].Kind)))
			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[1].Kind)))
			Expect([]string{"hash1", "hash2"}).To(ContainElement(string(transactions[0].Hash)))
			Expect([]string{"hash1", "hash2"}).To(ContainElement(string(transactions[1].Hash)))

			for _, t := range transactions {
				if t.ID == 1 {
					Expect(t.Labels).To(HaveLen(2))
					// Check that we have the expected labels
					labelValues := make([]string, len(t.Labels))
					for i, labelAssoc := range t.Labels {
						labelValues[i] = labelAssoc.Label.Value
					}
					Expect(labelValues).To(ContainElement("label1"))
					Expect(labelValues).To(ContainElement("label2"))
				}

				if t.ID == 2 {
					Expect(t.Labels).To(HaveLen(0))
				}
			}
		})

		It("query successfully transactions -- with start filter", func() {
			now := time.Now()
			sql, args, err := insertTransaction.
				Values(1, time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), 1001, "credit", "transaction", "1.1", "hash1").
				Values(2, time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.UTC), 1002, "debit", "transaction", "2.1", "hash2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			start := time.Date(now.Year(), 1, 3, 0, 0, 0, 0, time.UTC)
			transactions, err := dt.QueryTransactions(context.TODO(), pg.IntervalDateQueryFilter(&start, nil))
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))

			Expect(transactions[0].ID).To(Equal(int64(2)))
		})

		It("query successfully transactions -- with end filter", func() {
			now := time.Now()
			sql, args, err := insertTransaction.
				Values(1, time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), 1001, "credit", "transaction", "1.1", "hash1").
				Values(2, time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.UTC), 1001, "debit", "transaction", "2.1", "hash2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			end := time.Date(now.Year(), 1, 3, 0, 0, 0, 0, time.UTC)
			transactions, err := dt.QueryTransactions(context.TODO(), pg.IntervalDateQueryFilter(nil, &end))
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))

			Expect(transactions[0].ID).To(Equal(int64(1)))
		})

		It("query successfully transactions -- with before and after filter", func() {
			now := time.Now()
			sql, args, err := insertTransaction.
				Values(1, time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), 1001, "credit", "transaction", "1.1", "hash1").
				Values(2, time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.UTC), 1002, "debit", "transaction", "2.1", "hash2").
				Values(3, time.Date(now.Year(), 3, 1, 0, 0, 0, 0, time.UTC), 1003, "debit", "transaction", "2.1", "hash3").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			end := time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.UTC)
			start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			transactions, err := dt.QueryTransactions(context.TODO(), pg.IntervalDateQueryFilter(&start, &end))
			Expect(err).To(BeNil())

			Expect(transactions).To(HaveLen(2))
		})

		It("query successfully transactions -- with offset", func() {
			sql, args, err := insertTransaction.
				Values(1, time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				Values(2, time.Now(), 1002, "debit", "transaction", "2.1", "hash2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO(), pg.OffsetQueryFilter(1))
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))
		})

		It("query successfully transactions -- with limit", func() {
			sql, args, err := insertTransaction.
				Values(1, time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				Values(2, time.Now(), 1001, "debit", "transaction", "2.1", "hash2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO(), pg.LimitQueryFilter(1))
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))
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

	Context("labels", func() {
		It("successfully query labels -- no rule association", func() {
			sql, args, err := insertLabelStmt.
				Values("category", "label1").
				Values("category", "label2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			tags, err := dt.QueryLabels(context.TODO())
			Expect(err).To(BeNil())

			Expect(tags).To(HaveLen(2))
		})
		It("successfully query labels", func() {
			sql, args, err := insertRule.
				Values("rule1", "pattern").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertLabelStmt.
				Values("category", "label1").
				Values("category", "label2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTransaction.
				Values("1", time.Now(), 1001, "credit", "transaction", "1.1", "hash1").
				Values("2", time.Now(), 1001, "debit", "transaction", "2.1", "hash2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Get the actual label IDs from the database
			var label1ID, label2ID int
			err = pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'label1'").Scan(&label1ID)
			Expect(err).To(BeNil())
			err = pgPool.QueryRow(context.TODO(), "SELECT id FROM labels WHERE key = 'category' AND value = 'label2'").Scan(&label2ID)
			Expect(err).To(BeNil())

			sql, args, err = psql.Insert("rules_labels").Columns("rule_id", "label_id").
				Values("rule1", label1ID).
				Values("rule1", label2ID).
				Values("rule2", label1ID).
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			tags, err := dt.QueryLabels(context.TODO())
			Expect(err).To(BeNil())

			Expect(tags).To(HaveLen(2))
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

	Context("count", func() {
		It("count all transactions successfully", func() {
			// Create test rules
			sql, args, err := insertRule.
				Values("rule1", "pattern1").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test labels
			sql, args, err = insertLabelWithIDStmt.
				Values(1, "category", "food").
				Values(2, "category", "transport").
				Values(3, "merchant", "grocery").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = psql.Insert("rules_labels").
				Columns("rule_id", "label_id").
				Values("rule1", 1).
				Values("rule1", 3).
				Values("rule2", 2).
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test transactions
			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "debit", "grocery store", "50.0", "hash1").
				Values(2, time.Now(), 1002, "debit", "bus ticket", "5.0", "hash2").
				Values(3, time.Now(), 1003, "debit", "restaurant", "30.0", "hash3").
				Values(4, time.Now(), 1004, "debit", "taxi", "15.0", "hash4").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction-label relationships
			sql, args, err = insertTransactionLabel.
				Values(1, 1, "rule1"). // transaction 1 -> label 1 (category:food) -> rule1
				Values(1, 3, "rule1"). // transaction 1 -> label 3 (merchant:grocery) -> rule1
				Values(2, 2, "rule2"). // transaction 2 -> label 2 (category:transport) -> rule2
				Values(3, 1, "rule1"). // transaction 3 -> label 1 (category:food) -> rule1
				Values(4, 2, "rule2"). // transaction 4 -> label 2 (category:transport) -> rule2
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Test CountTransactions
			count, err := dt.CountTransactions(context.TODO())
			Expect(err).To(BeNil())
			Expect(count).To(Equal(4))
		})

		It("count all transactions successfully -- by label id", func() {
			// Create test rules
			sql, args, err := insertRule.
				Values("rule1", "pattern1").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test labels
			sql, args, err = insertLabelWithIDStmt.
				Values(1, "category", "food").
				Values(2, "category", "transport").
				Values(3, "merchant", "grocery").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = psql.Insert("rules_labels").
				Columns("rule_id", "label_id").
				Values("rule1", 1).
				Values("rule1", 3).
				Values("rule2", 2).
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test transactions
			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "debit", "grocery store", "50.0", "hash1").
				Values(2, time.Now(), 1002, "debit", "bus ticket", "5.0", "hash2").
				Values(3, time.Now(), 1003, "debit", "restaurant", "30.0", "hash3").
				Values(4, time.Now(), 1004, "debit", "taxi", "15.0", "hash4").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction-label relationships
			sql, args, err = insertTransactionLabel.
				Values(1, 1, "rule1"). // transaction 1 -> label 1 (category:food) -> rule1
				Values(1, 3, "rule1"). // transaction 1 -> label 3 (merchant:grocery) -> rule1
				Values(2, 2, "rule2"). // transaction 2 -> label 2 (category:transport) -> rule2
				Values(3, 1, "rule1"). // transaction 3 -> label 1 (category:food) -> rule1
				Values(4, 2, "rule2"). // transaction 4 -> label 2 (category:transport) -> rule2
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Test CountTransactions
			count, err := dt.CountTransactions(context.TODO(), pg.CountByLabelIDQueryFilter(1))
			Expect(err).To(BeNil())
			Expect(count).To(Equal(2))
		})

		It("should count 0 transactions where label not found", func() {
			// Create test rules
			sql, args, err := insertRule.
				Values("rule1", "pattern1").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test labels
			sql, args, err = insertLabelWithIDStmt.
				Values(1, "category", "food").
				Values(2, "category", "transport").
				Values(3, "merchant", "grocery").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = psql.Insert("rules_labels").
				Columns("rule_id", "label_id").
				Values("rule1", 1).
				Values("rule1", 3).
				Values("rule2", 2).
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test transactions
			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "debit", "grocery store", "50.0", "hash1").
				Values(2, time.Now(), 1002, "debit", "bus ticket", "5.0", "hash2").
				Values(3, time.Now(), 1003, "debit", "restaurant", "30.0", "hash3").
				Values(4, time.Now(), 1004, "debit", "taxi", "15.0", "hash4").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction-label relationships
			sql, args, err = insertTransactionLabel.
				Values(1, 1, "rule1"). // transaction 1 -> label 1 (category:food) -> rule1
				Values(1, 3, "rule1"). // transaction 1 -> label 3 (merchant:grocery) -> rule1
				Values(2, 2, "rule2"). // transaction 2 -> label 2 (category:transport) -> rule2
				Values(3, 1, "rule1"). // transaction 3 -> label 1 (category:food) -> rule1
				Values(4, 2, "rule2"). // transaction 4 -> label 2 (category:transport) -> rule2
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Test CountTransactions
			count, err := dt.CountTransactions(context.TODO(), pg.CountByLabelIDQueryFilter(10))
			Expect(err).To(BeNil())
			Expect(count).To(Equal(0))
		})

		It("count all transactions successfully -- by rule id", func() {
			// Create test rules
			sql, args, err := insertRule.
				Values("rule1", "pattern1").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test labels
			sql, args, err = insertLabelWithIDStmt.
				Values(1, "category", "food").
				Values(2, "category", "transport").
				Values(3, "merchant", "grocery").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = psql.Insert("rules_labels").
				Columns("rule_id", "label_id").
				Values("rule1", 1).
				Values("rule1", 3).
				Values("rule2", 2).
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create test transactions
			sql, args, err = insertTransaction.
				Values(1, time.Now(), 1001, "debit", "grocery store", "50.0", "hash1").
				Values(2, time.Now(), 1002, "debit", "bus ticket", "5.0", "hash2").
				Values(3, time.Now(), 1003, "debit", "restaurant", "30.0", "hash3").
				Values(4, time.Now(), 1004, "debit", "taxi", "15.0", "hash4").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction-label relationships
			sql, args, err = insertTransactionLabel.
				Values(1, 1, "rule1"). // transaction 1 -> label 1 (category:food) -> rule1
				Values(1, 3, "rule1"). // transaction 1 -> label 3 (merchant:grocery) -> rule1
				Values(2, 2, "rule2"). // transaction 2 -> label 2 (category:transport) -> rule2
				Values(3, 1, "rule1"). // transaction 3 -> label 1 (category:food) -> rule1
				Values(4, 2, "rule2"). // transaction 4 -> label 2 (category:transport) -> rule2
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Test CountTransactions
			count, err := dt.CountTransactions(context.TODO(), pg.CountByRuleIDQueryFilter("rule1"))
			Expect(err).To(BeNil())
			Expect(count).To(Equal(3))
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
