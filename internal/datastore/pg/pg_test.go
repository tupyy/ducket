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

const (
	pgUri           = "postgresql://postgres:postgres@localhost:5432/postgres"
	insertRulesStmt = "INSERT INTO rules VALUES ('%s', '%s', '%s');"
)

var (
	psql                 = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertRule           = psql.Insert("rules").Columns("id", "name", "pattern")
	insertTagStmt        = psql.Insert("tags").Columns("value")
	insertTransaction    = psql.Insert("transactions").Columns("id", "date", "kind", "content", "amount")
	insertTransactionTag = psql.Insert("transactions_tags").Columns("transaction_id", "tag_id", "rule_id")
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
	})

	AfterAll(func() {
		dt.Close()
	})

	Context("rule", func() {
		It("reads successfully rules -- empty response", func() {
			rules, err := dt.QueryRules(context.TODO(), pg.RuleFilter{}, &pg.QueryRuleOptions{})
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(0))
		})

		It("reads successfully rules -- without tags", func() {
			sql, args, err := insertRule.
				Values("rule1", "rule1", "pattern").
				Values("rule2", "rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rules, err := dt.QueryRules(context.TODO(), pg.RuleFilter{}, &pg.QueryRuleOptions{})
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(2))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[0].ID))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[1].ID))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[0].Name))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[1].Name))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[0].Pattern))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[1].Pattern))
		})

		It("reads successfully rules -- with tags", func() {
			sql, args, err := insertRule.
				Values("rule1", "rule1", "pattern").
				Values("rule2", "rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTagStmt.
				Values("tag1").
				Values("tag2").
				Values("tag3").
				Values("tag4").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = psql.Insert("rules_tags").Columns("rule_id", "tag").
				Values("rule1", "tag1").
				Values("rule1", "tag2").
				Values("rule2", "tag3").
				Values("rule2", "tag4").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rules, err := dt.QueryRules(context.TODO(), pg.RuleFilter{}, &pg.QueryRuleOptions{})
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(2))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[0].ID))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[1].ID))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[0].Name))
			Expect([]string{"rule1", "rule2"}).To(ContainElement(rules[1].Name))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[0].Pattern))
			Expect([]string{"pattern", "pattern2"}).To(ContainElement(rules[1].Pattern))

			Expect(rules[0].Tags).To(HaveLen(2))
			Expect(rules[1].Tags).To(HaveLen(2))
			for tag := range rules[0].Tags {
				Expect([]string{"tag1", "tag2", "tag3", "tag4"}).To(ContainElement(tag))
			}
			for tag := range rules[1].Tags {
				Expect([]string{"tag1", "tag2", "tag3", "tag4"}).To(ContainElement(tag))
			}

		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM tags;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("write rule", func() {
		It("write rule successfully -- no tags", func() {
			rule := entity.NewRule("rule1", "rule1", "pattern1")
			err := dt.WriteTx(context.TODO(), func(ctx context.Context, w pg.Writer) error {
				return w.WriteRule(context.TODO(), rule, false)
			})

			var count int

			rows, err := pgPool.Query(context.TODO(), "select count(*) from rules;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))
			}
		})

		It("update rule successfully -- no tags", func() {
			sql, args, err := insertRule.
				Values("rule1", "rule1", "pattern").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rule := entity.NewRule("rule1", "updated_name", "updated_pattern")
			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w pg.Writer) error {
				return w.WriteRule(context.TODO(), rule, true)
			})

			rows, err := pgPool.Query(context.TODO(), "select name,pattern from rules limit 1;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				var (
					pattern string
					name    string
				)
				err := rows.Scan(&name, &pattern)
				Expect(err).To(BeNil())
				Expect(pattern).To(Equal("updated_pattern"))
				Expect(name).To(Equal("updated_name"))
			}
		})

		It("write rule successfully -- with tags", func() {
			sql, args, err := insertTagStmt.
				Values("tag1").
				Values("tag2").
				Values("tag3").
				Values("tag4").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rule := entity.NewRule("rule1", "rule1", "pattern1", entity.Tag{Value: "tag1"}, entity.Tag{Value: "tag2"})
			err = dt.WriteTx(context.TODO(), func(ctx context.Context, w pg.Writer) error {
				return w.WriteRule(context.TODO(), rule, false)
			})
			Expect(err).To(BeNil())

			var count int

			rows, err := pgPool.Query(context.TODO(), "select count(*) from rules;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))
			}

			count = 0
			rows, err = pgPool.Query(context.TODO(), "select count(*) from rules_tags;")
			Expect(err).To(BeNil())
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(2))
			}
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM tags;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("transaction", func() {
		It("query successfully transactions -- empty response", func() {
			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{}, &pg.QueryTransactionOptions{})
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(0))
		})

		It("query successfully transactions -- no tags", func() {
			sql, args, err := insertTransaction.
				Values("1", time.Now(), "credit", "transaction", "1.1").
				Values("2", time.Now(), "debit", "transaction", "2.1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{}, &pg.QueryTransactionOptions{})
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(2))

			Expect([]string{"1", "2"}).To(ContainElement(transactions[0].ID))
			Expect([]string{"1", "2"}).To(ContainElement(transactions[1].ID))

			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[0].Kind)))
			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[1].Kind)))
		})

		It("query successfully transactions -- with tags", func() {
			sql, args, err := insertRule.
				Values("rule1", "rule1", "pattern").
				Values("rule2", "rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTagStmt.
				Values("tag1").
				Values("tag2").
				Values("tag3").
				Values("tag4").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTransaction.
				Values("1", time.Now(), "credit", "transaction", "1.1").
				Values("2", time.Now(), "debit", "transaction", "2.1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			sql, args, err = insertTransactionTag.
				Values("1", "tag1", "rule1").
				Values("1", "tag2", "rule2").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{}, &pg.QueryTransactionOptions{})
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(2))

			Expect([]string{"1", "2"}).To(ContainElement(transactions[0].ID))
			Expect([]string{"1", "2"}).To(ContainElement(transactions[1].ID))

			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[0].Kind)))
			Expect([]string{"credit", "debit"}).To(ContainElement(string(transactions[1].Kind)))

			for _, t := range transactions {
				if t.ID == "1" {
					Expect(t.Tags).To(HaveLen(2))
					Expect(t.Tags["tag1"]).To(Equal(entity.Tag{Value: "tag1", RuleID: "rule1"}))
					Expect(t.Tags["tag2"]).To(Equal(entity.Tag{Value: "tag2", RuleID: "rule2"}))
				}

				if t.ID == "2" {
					Expect(t.Tags).To(HaveLen(0))
				}
			}
		})

		It("query successfully transactions -- with before filter", func() {
			now := time.Now()
			sql, args, err := insertTransaction.
				Values("1", time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), "credit", "transaction", "1.1").
				Values("2", time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.UTC), "debit", "transaction", "2.1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			beforeTime := time.Date(now.Year(), 1, 3, 0, 0, 0, 0, time.UTC)
			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{
				Before: &beforeTime,
			}, &pg.QueryTransactionOptions{})
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))

			Expect(transactions[0].ID).To(Equal("1"))
		})

		It("query successfully transactions -- with after filter", func() {
			now := time.Now()
			sql, args, err := insertTransaction.
				Values("1", time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), "credit", "transaction", "1.1").
				Values("2", time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.UTC), "debit", "transaction", "2.1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			afterTime := time.Date(now.Year(), 1, 3, 0, 0, 0, 0, time.UTC)
			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{
				After: &afterTime,
			}, &pg.QueryTransactionOptions{})
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))

			Expect(transactions[0].ID).To(Equal("2"))
		})

		It("query successfully transactions -- with before and after filter", func() {
			now := time.Now()
			sql, args, err := insertTransaction.
				Values("1", time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), "credit", "transaction", "1.1").
				Values("2", time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.UTC), "debit", "transaction", "2.1").
				Values("3", time.Date(now.Year(), 3, 1, 0, 0, 0, 0, time.UTC), "debit", "transaction", "2.1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			afterTime := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			beforeTime := time.Date(now.Year(), 2, 2, 0, 0, 0, 0, time.UTC)
			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{
				Before: &beforeTime,
				After:  &afterTime,
			}, &pg.QueryTransactionOptions{})
			Expect(err).To(BeNil())

			Expect(transactions).To(HaveLen(2))
			Expect([]string{"1", "2"}).To(ContainElement(transactions[0].ID))
			Expect([]string{"1", "2"}).To(ContainElement(transactions[1].ID))
		})

		It("query successfully transactions -- with offset", func() {
			sql, args, err := insertTransaction.
				Values("1", time.Now(), "credit", "transaction", "1.1").
				Values("2", time.Now(), "debit", "transaction", "2.1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{}, &pg.QueryTransactionOptions{Offset: 1})
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))
		})

		It("query successfully transactions -- with limit", func() {
			sql, args, err := insertTransaction.
				Values("1", time.Now(), "credit", "transaction", "1.1").
				Values("2", time.Now(), "debit", "transaction", "2.1").
				ToSql()

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transactions, err := dt.QueryTransactions(context.TODO(), pg.TransactionFilter{}, &pg.QueryTransactionOptions{Limit: 1})
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(1))
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM tags;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})
})
