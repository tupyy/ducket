package store_test

import (
	"context"
	"database/sql"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
)

var _ = Describe("Dynamic Tag CTE", func() {
	var (
		ctx context.Context
		s   *store.Store
		db  *sql.DB
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = store.NewDB(":memory:")
		Expect(err).NotTo(HaveOccurred())

		s = store.NewStore(db)
		err = s.Migrate(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	date := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	}

	insertTxn := func(content string, kind entity.TransactionKind, amount float64) int64 {
		id, err := s.CreateTransaction(ctx, entity.Transaction{
			Hash:    content + time.Now().String(),
			Date:    date(2025, 1, 15),
			Account: 1001,
			Kind:    kind,
			Amount:  amount,
			Content: content,
		})
		Expect(err).NotTo(HaveOccurred())
		return id
	}

	insertRule := func(name, filter string, tags []string) {
		_, err := s.CreateRule(ctx, entity.Rule{
			Name:   name,
			Filter: filter,
			Tags:   tags,
		})
		Expect(err).NotTo(HaveOccurred())
	}

	Context("when no rules exist", func() {
		It("should return empty tags", func() {
			id := insertTxn("KAUFLAND BUCURESTI", entity.DebitTransaction, 150)

			txn, err := s.GetTransaction(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.Tags).To(BeEmpty())
		})
	})

	Context("when a single rule matches", func() {
		It("should assign that rule's tags", func() {
			id := insertTxn("KAUFLAND BUCURESTI", entity.DebitTransaction, 150)
			insertRule("grocery", "content ~ /KAUFLAND|LIDL/", []string{"food", "groceries"})

			txn, err := s.GetTransaction(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.Tags).To(ConsistOf("food", "groceries"))
		})
	})

	Context("when multiple rules match the same transaction", func() {
		It("should merge and deduplicate tags", func() {
			id := insertTxn("KAUFLAND BUCURESTI", entity.DebitTransaction, 150)
			insertRule("grocery", "content ~ /KAUFLAND/", []string{"food", "groceries"})
			insertRule("big-spend", "amount > 100", []string{"expensive", "food"})

			txn, err := s.GetTransaction(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.Tags).To(ConsistOf("food", "groceries", "expensive"))
		})
	})

	Context("when a rule does not match", func() {
		It("should not assign tags from that rule", func() {
			id := insertTxn("SHELL PETROL", entity.DebitTransaction, 80)
			insertRule("grocery", "content ~ /KAUFLAND|LIDL/", []string{"food"})

			txn, err := s.GetTransaction(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.Tags).To(BeEmpty())
		})
	})

	Context("when a rule has an invalid filter", func() {
		It("should skip that rule and not error", func() {
			id := insertTxn("KAUFLAND BUCURESTI", entity.DebitTransaction, 150)

			// Insert a rule with invalid filter directly via SQL
			_, err := db.ExecContext(ctx,
				"INSERT INTO rules (name, filter, tags) VALUES (?, ?, ['broken'])",
				"bad-rule", "invalid %%% filter",
			)
			Expect(err).NotTo(HaveOccurred())

			insertRule("grocery", "content ~ /KAUFLAND/", []string{"food"})

			txn, err := s.GetTransaction(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.Tags).To(ConsistOf("food"))
		})
	})

	Context("with ListTransactions", func() {
		It("should include tags in listed results", func() {
			insertTxn("KAUFLAND BUCURESTI", entity.DebitTransaction, 150)
			insertTxn("SHELL PETROL", entity.DebitTransaction, 80)
			insertTxn("SALARIU", entity.CreditTransaction, 5000)

			insertRule("grocery", "content ~ /KAUFLAND/", []string{"food"})
			insertRule("fuel", "content ~ /SHELL/", []string{"transport"})

			txns, err := s.ListTransactions(ctx, "kind = 'debit'", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(2))

			tagMap := map[string][]string{}
			for _, t := range txns {
				tagMap[t.Content] = t.Tags
			}
			Expect(tagMap["KAUFLAND BUCURESTI"]).To(ConsistOf("food"))
			Expect(tagMap["SHELL PETROL"]).To(ConsistOf("transport"))
		})
	})

	Context("filtering on multiple transaction columns", func() {
		It("should support kind + amount + content filters in rules", func() {
			id := insertTxn("KAUFLAND MEGA", entity.DebitTransaction, 500)
			insertTxn("KAUFLAND SMALL", entity.DebitTransaction, 30)

			insertRule("big-grocery", "content ~ /KAUFLAND/ and amount > 100 and kind = 'debit'", []string{"big-grocery"})

			txn, err := s.GetTransaction(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.Tags).To(ConsistOf("big-grocery"))

			// The small one should not match
			txns, err := s.ListTransactions(ctx, "content ~ /KAUFLAND SMALL/", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(1))
			Expect(txns[0].Tags).To(BeEmpty())
		})
	})
})
