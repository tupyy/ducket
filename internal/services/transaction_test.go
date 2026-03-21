package services_test

import (
	"context"
	"database/sql"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
	srvErrors "git.tls.tupangiu.ro/cosmin/finante/pkg/errors"
)

var _ = Describe("TransactionService", func() {
	var (
		ctx context.Context
		svc *services.TransactionService
		db  *sql.DB
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = store.NewDB(":memory:")
		Expect(err).NotTo(HaveOccurred())

		st := store.NewStore(db)
		err = st.Migrate(ctx)
		Expect(err).NotTo(HaveOccurred())

		svc = services.NewTransactionService(st)
	})

	AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	date := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	}

	makeTxn := func(content string, amount float64, dt time.Time) entity.Transaction {
		return *entity.NewTransaction(entity.DebitTransaction, 1001, dt, amount, content)
	}

	Describe("Create", func() {
		It("should create a transaction and return it with ID", func() {
			txn, err := svc.Create(ctx, makeTxn("KAUFLAND", 150, date(2025, 1, 15)))
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.ID).To(BeNumerically(">", 0))
		})

		It("should reject duplicate hash", func() {
			txn := makeTxn("DUPLICATE", 100, date(2025, 1, 1))
			_, err := svc.Create(ctx, txn)
			Expect(err).NotTo(HaveOccurred())

			_, err = svc.Create(ctx, txn)
			Expect(err).To(HaveOccurred())
			Expect(srvErrors.IsDuplicateResourceError(err)).To(BeTrue())
		})
	})

	Describe("Get", func() {
		It("should return a transaction by ID", func() {
			created, err := svc.Create(ctx, makeTxn("SHELL", 80, date(2025, 3, 10)))
			Expect(err).NotTo(HaveOccurred())

			txn, err := svc.Get(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(txn.Content).To(Equal("SHELL"))
		})

		It("should return not found for non-existent ID", func() {
			_, err := svc.Get(ctx, 999)
			Expect(err).To(HaveOccurred())
			Expect(srvErrors.IsResourceNotFoundError(err)).To(BeTrue())
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			_, err := svc.Create(ctx, makeTxn("KAUFLAND", 150, date(2025, 1, 5)))
			Expect(err).NotTo(HaveOccurred())
			_, err = svc.Create(ctx, makeTxn("SHELL", 80, date(2025, 1, 10)))
			Expect(err).NotTo(HaveOccurred())
			_, err = svc.Create(ctx, makeTxn("CHIRIE", 1500, date(2025, 2, 1)))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should list all transactions", func() {
			txns, err := svc.List(ctx, store.NoFilter, 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(3))
		})

		It("should filter transactions", func() {
			txns, err := svc.List(ctx, "amount > 100", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(2))
		})

		It("should paginate", func() {
			txns, err := svc.List(ctx, store.NoFilter, 2, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(txns).To(HaveLen(2))
		})
	})

	Describe("Update", func() {
		It("should update a transaction", func() {
			created, err := svc.Create(ctx, makeTxn("ORIGINAL", 100, date(2025, 1, 1)))
			Expect(err).NotTo(HaveOccurred())

			created.Content = "UPDATED"
			created.Amount = 200
			err = svc.Update(ctx, created)
			Expect(err).NotTo(HaveOccurred())

			updated, err := svc.Get(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Content).To(Equal("UPDATED"))
			Expect(updated.Amount).To(BeNumerically("~", 200, 0.01))
		})

		It("should return not found for non-existent transaction", func() {
			txn := makeTxn("GHOST", 100, date(2025, 1, 1))
			txn.ID = 999
			err := svc.Update(ctx, txn)
			Expect(err).To(HaveOccurred())
			Expect(srvErrors.IsResourceNotFoundError(err)).To(BeTrue())
		})
	})

	Describe("Delete", func() {
		It("should delete a transaction", func() {
			created, err := svc.Create(ctx, makeTxn("TO DELETE", 50, date(2025, 1, 1)))
			Expect(err).NotTo(HaveOccurred())

			err = svc.Delete(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = svc.Get(ctx, created.ID)
			Expect(err).To(HaveOccurred())
			Expect(srvErrors.IsResourceNotFoundError(err)).To(BeTrue())
		})
	})
})
