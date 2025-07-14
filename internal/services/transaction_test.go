package services_test

import (
	"context"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TransactionService", Ordered, func() {
	var (
		transactionService *services.TransactionService
		datastore          *pg.Datastore
		pgPool             *pgxpool.Pool
		testTime           time.Time
	)

	BeforeAll(func() {
		// Setup database connection
		pgDt, err := pg.NewPostgresDatastore(context.TODO(), pgUri)
		Expect(err).To(BeNil())
		Expect(pgDt).ToNot(BeNil())

		pgxConfig, err := pgxpool.ParseConfig(pgUri)
		Expect(err).To(BeNil())

		pool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
		Expect(err).To(BeNil())
		Expect(pool).ToNot(BeNil())

		datastore = pgDt
		pgPool = pool
		transactionService = services.NewTransactionService(datastore)
		testTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Clean up database
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		datastore.Close()
	})

	Context("GetTransactions", func() {
		It("should return empty slice when no transactions exist", func() {
			filter := services.NewTransactionFilterWithOptions()
			transactions, err := transactionService.GetTransactions(context.TODO(), filter)
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(0))
		})

		It("should return transactions when they exist", func() {
			// Insert test transactions
			insertTransactionStmt := psql.Insert("transactions").
				Columns("id", "date", "account", "kind", "content", "amount", "hash")
			sql, args, err := insertTransactionStmt.
				Values(1, testTime, 1001, "debit", "Test tx 1", -100.50, "hash1").
				Values(2, testTime.Add(24*time.Hour), 1002, "credit", "Test tx 2", 200.75, "hash2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			filter := services.NewTransactionFilterWithOptions()
			transactions, err := transactionService.GetTransactions(context.TODO(), filter)
			Expect(err).To(BeNil())
			Expect(transactions).To(HaveLen(2))

			// Verify transaction details
			var tx1, tx2 *entity.Transaction
			for i, tx := range transactions {
				if tx.Hash == "hash1" {
					tx1 = &transactions[i]
				} else if tx.Hash == "hash2" {
					tx2 = &transactions[i]
				}
			}

			Expect(tx1).ToNot(BeNil())
			Expect(tx2).ToNot(BeNil())
			Expect(tx1.Kind).To(Equal(entity.DebitTransaction))
			Expect(tx2.Kind).To(Equal(entity.CreditTransaction))
			Expect(tx1.Amount).To(Equal(float32(-100.50)))
			Expect(tx2.Amount).To(Equal(float32(200.75)))
		})

		It("should filter transactions by date range", func() {
			// Insert transactions with different dates
			insertTransactionStmt := psql.Insert("transactions").
				Columns("id", "date", "account", "kind", "content", "amount", "hash")

			oldDate := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
			newDate := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

			sql, args, err := insertTransactionStmt.
				Values(10, oldDate, 1001, "debit", "Old tx", -50.0, "old_hash").
				Values(11, newDate, 1002, "credit", "New tx", 100.0, "new_hash").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Filter for transactions after 2024-01-01
			startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			filter := services.NewTransactionFilterWithOptions(
				services.WithStart(&startDate),
			)
			transactions, err := transactionService.GetTransactions(context.TODO(), filter)
			Expect(err).To(BeNil())

			// Should only return the new transaction
			found := false
			for _, tx := range transactions {
				if tx.Hash == "new_hash" {
					found = true
				}
				if tx.Hash == "old_hash" {
					Fail("Should not return old tx")
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should apply limit and offset filters", func() {
			// Insert multiple transactions
			insertTransactionStmt := psql.Insert("transactions").
				Columns("id", "date", "account", "kind", "content", "amount", "hash")

			sql, args, err := insertTransactionStmt.
				Values(20, testTime, 1001, "debit", "Transaction 1", -10.0, "limit_hash_1").
				Values(21, testTime, 1001, "debit", "Transaction 2", -20.0, "limit_hash_2").
				Values(22, testTime, 1001, "debit", "Transaction 3", -30.0, "limit_hash_3").
				Values(23, testTime, 1001, "debit", "Transaction 4", -40.0, "limit_hash_4").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Test limit
			filter := services.NewTransactionFilterWithOptions(
				services.WithLimit(2),
			)
			transactions, err := transactionService.GetTransactions(context.TODO(), filter)
			Expect(err).To(BeNil())
			Expect(len(transactions)).To(BeNumerically("<=", 2))

			// Test offset
			filter = services.NewTransactionFilterWithOptions(
				services.WithOffset(2),
				services.WithLimit(10),
			)
			transactions, err = transactionService.GetTransactions(context.TODO(), filter)
			Expect(err).To(BeNil())
			// Should have transactions but not the first 2
		})
	})

	Context("GetTransaction", func() {
		It("should return nil when transaction doesn't exist", func() {
			transaction, err := transactionService.GetTransaction(context.TODO(), "no-exist-hash")
			Expect(err).To(BeNil())
			Expect(transaction).To(BeNil())
		})

		It("should return transaction when it exists", func() {
			// Insert test transaction
			insertTransactionStmt := psql.Insert("transactions").
				Columns("id", "date", "account", "kind", "content", "amount", "hash")
			sql, args, err := insertTransactionStmt.
				Values(100, testTime, 1001, "credit", "Single tx test", 150.25, "single_hash").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			transaction, err := transactionService.GetTransaction(context.TODO(), "single_hash")
			Expect(err).To(BeNil())
			Expect(transaction).ToNot(BeNil())
			Expect(transaction.Hash).To(Equal("single_hash"))
			Expect(transaction.Kind).To(Equal(entity.CreditTransaction))
			Expect(transaction.Amount).To(Equal(float32(150.25)))
			Expect(transaction.RawContent).To(Equal("Single tx test"))
		})
	})

	Context("CreateOrUpdate", func() {
		It("should create a new transaction", func() {
			newTransaction := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-75.50,
				"New tx content",
			)

			result, err := transactionService.CreateOrUpdate(context.TODO(), *newTransaction)
			Expect(err).To(BeNil())
			Expect(result.ID).To(BeNumerically(">", 0))
			Expect(result.Hash).To(Equal(newTransaction.Hash))

			// Verify transaction was created
			retrieved, err := transactionService.GetTransaction(context.TODO(), newTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Amount).To(Equal(float32(-75.50)))
			Expect(retrieved.Kind).To(Equal(entity.DebitTransaction))
		})

		It("should update existing transaction", func() {
			// First, create a transaction
			originalTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				100.0,
				"Original",
			)

			created, err := transactionService.CreateOrUpdate(context.TODO(), *originalTransaction)
			Expect(err).To(BeNil())
			originalID := created.ID

			// Update the same transaction (same hash, different content)
			updatedTransaction := &entity.Transaction{
				Hash:       originalTransaction.Hash, // Same hash
				Kind:       entity.CreditTransaction,
				Date:       testTime,
				Account:    1001,
				Amount:     200.0,     // Different amount
				RawContent: "Updated", // Different content
				Labels:     make(map[int]string),
			}

			result, err := transactionService.CreateOrUpdate(context.TODO(), *updatedTransaction)
			Expect(err).To(BeNil())
			Expect(result.ID).To(Equal(originalID)) // Should keep same ID
			Expect(result.Amount).To(Equal(float32(200.0)))

			// Verify transaction was updated
			retrieved, err := transactionService.GetTransaction(context.TODO(), originalTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.ID).To(Equal(originalID))
			Expect(retrieved.Amount).To(Equal(float32(200.0)))
			Expect(retrieved.RawContent).To(Equal("Updated"))
		})

		It("should handle transactions with labels", func() {
			// First, insert required labels
			insertLabelStmt := psql.Insert("labels").Columns("id", "key", "value")
			sql, args, err := insertLabelStmt.
				Values(1, "category", "expense").
				Values(2, "type", "recurring").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Insert required rules
			insertRuleStmt := psql.Insert("rules").Columns("id", "pattern")
			sql, args, err = insertRuleStmt.
				Values("rule1", "pattern1").
				Values("rule2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Link rules to labels
			insertRuleLabelStmt := psql.Insert("rules_labels").Columns("rule_id", "label_id")
			sql, args, err = insertRuleLabelStmt.
				Values("rule1", 1).
				Values("rule2", 2).
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction with labels
			transactionWithLabels := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-50.0,
				"Tx with labels",
			)
			transactionWithLabels.Labels = map[int]string{
				1: "rule1", // label_id -> rule_id
				2: "rule2",
			}

			result, err := transactionService.CreateOrUpdate(context.TODO(), *transactionWithLabels)
			Expect(err).To(BeNil())
			Expect(result.ID).To(BeNumerically(">", 0))

			// Verify transaction was created
			retrieved, err := transactionService.GetTransaction(context.TODO(), transactionWithLabels.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Labels).To(HaveLen(2))
		})
	})

	Context("Delete", func() {
		It("should delete existing transaction successfully", func() {
			// First, create a transaction
			transactionToDelete := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-25.0,
				"Tx to delete",
			)

			created, err := transactionService.CreateOrUpdate(context.TODO(), *transactionToDelete)
			Expect(err).To(BeNil())
			transactionID := created.ID

			// Verify transaction exists
			retrieved, err := transactionService.GetTransaction(context.TODO(), transactionToDelete.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())

			// Delete the transaction
			err = transactionService.Delete(context.TODO(), transactionID)
			Expect(err).To(BeNil())

			// Verify transaction was deleted
			retrieved, err = transactionService.GetTransaction(context.TODO(), transactionToDelete.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).To(BeNil())
		})

		It("should handle deletion of non-existent transaction gracefully", func() {
			err := transactionService.Delete(context.TODO(), 99999)
			Expect(err).To(BeNil()) // Should not error
		})
	})

	Context("TransactionFilter", func() {
		It("should generate correct query filters", func() {
			startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			endTime := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

			filter := services.NewTransactionFilterWithOptions(
				services.WithStart(&startTime),
				services.WithEnd(&endTime),
				services.WithLimit(50),
				services.WithOffset(10),
			)

			queryFilters := filter.QueriesFn()
			Expect(queryFilters).To(HaveLen(3)) // IntervalDateQueryFilter, LimitQueryFilter, OffsetQueryFilter
		})

		It("should handle empty filter", func() {
			filter := services.NewTransactionFilterWithOptions()
			queryFilters := filter.QueriesFn()
			Expect(queryFilters).To(HaveLen(3)) // Should still have the 3 basic filters
		})
	})

	Context("Edge Cases", func() {
		It("should handle zero amount transactions", func() {
			zeroTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				0.0,
				"Zero amount",
			)

			result, err := transactionService.CreateOrUpdate(context.TODO(), *zeroTransaction)
			Expect(err).To(BeNil())
			Expect(result.Amount).To(Equal(float32(0.0)))

			retrieved, err := transactionService.GetTransaction(context.TODO(), zeroTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Amount).To(Equal(float32(0.0)))
		})

		It("should handle very large amounts", func() {
			largeAmount := float32(999999.99)
			largeTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				largeAmount,
				"Large amount",
			)

			result, err := transactionService.CreateOrUpdate(context.TODO(), *largeTransaction)
			Expect(err).To(BeNil())
			Expect(result.Amount).To(Equal(largeAmount))
		})

		It("should handle empty content", func() {
			emptyContentTransaction := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-10.0,
				"", // Empty content
			)

			result, err := transactionService.CreateOrUpdate(context.TODO(), *emptyContentTransaction)
			Expect(err).To(BeNil())
			Expect(result.RawContent).To(Equal(""))

			retrieved, err := transactionService.GetTransaction(context.TODO(), emptyContentTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.RawContent).To(Equal(""))
		})
	})

	AfterEach(func() {
		// Clean up after each test
		_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
		Expect(err).To(BeNil())
	})
})
