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

var _ = Describe("TransactionStore", func() {
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

	createTxn := func(kind entity.TransactionKind, dt time.Time, amount float64, content string) entity.Transaction {
		info := "some info"
		t := entity.Transaction{
			Hash:    content + dt.String(),
			Date:    dt,
			Account: 1001,
			Kind:    kind,
			Amount:  amount,
			Content: content,
			Info:    &info,
		}
		id, err := s.CreateTransaction(ctx, t)
		Expect(err).NotTo(HaveOccurred())
		t.ID = id
		return t
	}

	Describe("CreateTransaction", func() {
		It("should create a transaction and return its ID", func() {
			id, err := s.CreateTransaction(ctx, entity.Transaction{
				Hash:    "abc123",
				Date:    date(2025, 1, 15),
				Account: 1001,
				Kind:    entity.DebitTransaction,
				Amount:  42.50,
				Content: "KAUFLAND BUCURESTI",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(BeNumerically(">", 0))
		})

		It("should reject duplicate hashes", func() {
			createTxn(entity.DebitTransaction, date(2025, 1, 1), 10, "DUPLICATE")

			_, err := s.CreateTransaction(ctx, entity.Transaction{
				Hash:    "DUPLICATE" + date(2025, 1, 1).String(),
				Date:    date(2025, 1, 2),
				Account: 1001,
				Kind:    entity.DebitTransaction,
				Amount:  20,
				Content: "OTHER",
			})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetTransaction", func() {
		It("should return a transaction by ID", func() {
			created := createTxn(entity.DebitTransaction, date(2025, 3, 10), 99.99, "SHELL PETROL")

			txn, err := s.GetTransaction(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn).NotTo(BeNil())
			Expect(txn.Content).To(Equal("SHELL PETROL"))
			Expect(txn.Amount).To(BeNumerically("~", 99.99, 0.01))
			Expect(txn.Kind).To(Equal(entity.DebitTransaction))
			Expect(txn.Info).NotTo(BeNil())
			Expect(*txn.Info).To(Equal("some info"))
		})

		It("should return nil for non-existent ID", func() {
			txn, err := s.GetTransaction(ctx, 999)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn).To(BeNil())
		})
	})

	Describe("GetTransactionByHash", func() {
		It("should find by hash", func() {
			created := createTxn(entity.CreditTransaction, date(2025, 5, 1), 3000, "SALARIU")

			txn, err := s.GetTransactionByHash(ctx, created.Hash)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn).NotTo(BeNil())
			Expect(txn.ID).To(Equal(created.ID))
		})

		It("should return nil for non-existent hash", func() {
			txn, err := s.GetTransactionByHash(ctx, "nonexistent")
			Expect(err).NotTo(HaveOccurred())
			Expect(txn).To(BeNil())
		})
	})

	Describe("ListTransactions", func() {
		BeforeEach(func() {
			createTxn(entity.DebitTransaction, date(2025, 1, 5), 150, "KAUFLAND BUCURESTI")
			createTxn(entity.DebitTransaction, date(2025, 1, 10), 80, "LIDL SECTOR 3")
			createTxn(entity.DebitTransaction, date(2025, 2, 1), 200, "SHELL PETROL")
			createTxn(entity.CreditTransaction, date(2025, 1, 15), 5000, "SALARIU IANUARIE")
			createTxn(entity.DebitTransaction, date(2025, 3, 1), 1500, "CHIRIE MARTIE")
		})

		It("should list all without filter", func() {
			txns, err := s.ListTransactions(ctx, store.NoFilter, 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(5))
		})

		It("should filter by kind", func() {
			txns, err := s.ListTransactions(ctx, "kind = 'debit'", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(4))
		})

		It("should filter by content regex", func() {
			txns, err := s.ListTransactions(ctx, "content ~ /KAUFLAND|LIDL/", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(2))
		})

		It("should filter by amount", func() {
			txns, err := s.ListTransactions(ctx, "amount > 1000", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(2))
		})

		It("should filter by date range", func() {
			txns, err := s.ListTransactions(ctx, "date >= '2025-01-01' and date < '2025-02-01'", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(3))
		})

		It("should combine filters", func() {
			txns, err := s.ListTransactions(ctx, "kind = 'debit' and amount > 100", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(3))
		})

		It("should paginate with limit", func() {
			txns, err := s.ListTransactions(ctx, store.NoFilter, 2, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(2))
		})

		It("should paginate with offset", func() {
			txns, err := s.ListTransactions(ctx, store.NoFilter, 2, 3)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(2))
		})

		It("should order by date descending", func() {
			txns, err := s.ListTransactions(ctx, store.NoFilter, 0, 0)
			Expect(err).NotTo(HaveOccurred())
			for i := 1; i < len(txns); i++ {
				Expect(txns[i-1].Date).To(BeTemporally(">=", txns[i].Date))
			}
		})

		It("should reject invalid filter", func() {
			_, err := s.ListTransactions(ctx, "bogus %%% field", 0, 0)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("UpdateTransaction", func() {
		It("should update transaction fields", func() {
			created := createTxn(entity.DebitTransaction, date(2025, 1, 1), 100, "ORIGINAL")

			newInfo := "updated info"
			err := s.UpdateTransaction(ctx, entity.Transaction{
				ID:      created.ID,
				Hash:    created.Hash,
				Date:    date(2025, 1, 2),
				Account: 1002,
				Kind:    entity.CreditTransaction,
				Amount:  200,
				Content: "UPDATED",
				Info:    &newInfo,
			})
			Expect(err).NotTo(HaveOccurred())

			updated, err := s.GetTransaction(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Content).To(Equal("UPDATED"))
			Expect(updated.Amount).To(BeNumerically("~", 200, 0.01))
			Expect(updated.Kind).To(Equal(entity.CreditTransaction))
			Expect(*updated.Info).To(Equal("updated info"))
		})
	})

	Describe("DeleteTransaction", func() {
		It("should delete a transaction", func() {
			created := createTxn(entity.DebitTransaction, date(2025, 1, 1), 50, "TO DELETE")

			err := s.DeleteTransaction(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())

			txn, err := s.GetTransaction(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn).To(BeNil())
		})
	})
})
