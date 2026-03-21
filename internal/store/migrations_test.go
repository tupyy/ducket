package store_test

import (
	"context"
	"database/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"git.tls.tupangiu.ro/cosmin/finante/internal/store"
)

var _ = Describe("Migrations", func() {
	var (
		ctx context.Context
		db  *sql.DB
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = store.NewDB(":memory:")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	It("should create all tables", func() {
		st := store.NewStore(db)
		err := st.Migrate(ctx)
		Expect(err).NotTo(HaveOccurred())

		// Verify rules table exists
		_, err = db.ExecContext(ctx, "SELECT id, name, filter, tags, created_at FROM rules LIMIT 0")
		Expect(err).NotTo(HaveOccurred())

		// Verify transactions table exists
		_, err = db.ExecContext(ctx, "SELECT id, hash, date, account, kind, amount, content, info, recipient, created_at FROM transactions LIMIT 0")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be idempotent", func() {
		st := store.NewStore(db)

		err := st.Migrate(ctx)
		Expect(err).NotTo(HaveOccurred())

		err = st.Migrate(ctx)
		Expect(err).NotTo(HaveOccurred())
	})
})
