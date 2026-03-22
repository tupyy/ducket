package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tupyy/ducket/internal/entity"
	"github.com/tupyy/ducket/internal/store"
)

var _ = Describe("WithTx", func() {
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

	It("should commit on success", func() {
		err := s.WithTx(ctx, func(ctx context.Context) error {
			_, err := s.CreateRule(ctx, entity.Rule{
				Name:   "committed",
				Filter: "content ~ /TEST/",
				Tags:   []string{"test"},
			})
			return err
		})
		Expect(err).NotTo(HaveOccurred())

		rules, err := s.ListRules(ctx, "name = 'committed'", 0, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(rules).To(HaveLen(1))
	})

	It("should rollback on error", func() {
		err := s.WithTx(ctx, func(ctx context.Context) error {
			_, err := s.CreateRule(ctx, entity.Rule{
				Name:   "should-not-exist",
				Filter: "content ~ /TEST/",
				Tags:   []string{"test"},
			})
			if err != nil {
				return err
			}
			return fmt.Errorf("forced error")
		})
		Expect(err).To(HaveOccurred())

		rules, err := s.ListRules(ctx, "name = 'should-not-exist'", 0, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(rules).To(BeEmpty())
	})

	It("should be atomic across multiple operations", func() {
		err := s.WithTx(ctx, func(ctx context.Context) error {
			_, err := s.CreateRule(ctx, entity.Rule{
				Name:   "rule-in-tx",
				Filter: "content ~ /TX/",
				Tags:   []string{"tx-test"},
			})
			if err != nil {
				return err
			}

			_, err = s.CreateTransaction(ctx, entity.Transaction{
				Hash:    "tx-hash-123",
				Date:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Account: 1001,
				Kind:    entity.DebitTransaction,
				Amount:  100,
				Content: "TX TEST",
			})
			return err
		})
		Expect(err).NotTo(HaveOccurred())

		rules, err := s.ListRules(ctx, store.NoFilter, 0, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(rules).To(HaveLen(1))

		txns, _, err := s.ListTransactions(ctx, store.NoFilter, nil, nil, 0, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(txns).To(HaveLen(1))
	})
})
