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

	Context("Create", func() {
		It("should create a new transaction successfully", func() {
			newTransaction := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-75.50,
				"New tx content",
			)

			result, err := transactionService.Create(context.TODO(), *newTransaction)
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

		It("should return error when transaction already exists", func() {
			// First, create a transaction
			originalTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				100.0,
				"Original",
			)

			_, err := transactionService.Create(context.TODO(), *originalTransaction)
			Expect(err).To(BeNil())

			// Try to create the same transaction again (same hash)
			duplicateTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				100.0,
				"Original", // Same content will generate same hash
			)

			_, err = transactionService.Create(context.TODO(), *duplicateTransaction)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("already exists"))
		})

		It("should create relationships between transaction and labels", func() {
			// First, create some labels and rules
			insertLabelStmt := psql.Insert("labels").Columns("id", "key", "value")
			sql, args, err := insertLabelStmt.
				Values(101, "category", "expense").
				Values(102, "type", "recurring").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

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
				Values("rule1", 101).
				Values("rule2", 102).
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction with label assignments that include rule references
			transactionWithLabels := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-50.0,
				"Tx with relationship test",
			)

			// Simulate labels with rule associations
			ruleID1 := "rule1"
			ruleID2 := "rule2"
			transactionWithLabels.Labels = []entity.LabelAssociation{
				{
					Label:  entity.Label{ID: 101, Key: "category", Value: "expense"},
					RuleID: &ruleID1,
				},
				{
					Label:  entity.Label{ID: 102, Key: "type", Value: "recurring"},
					RuleID: &ruleID2,
				},
			}

			result, err := transactionService.Create(context.TODO(), *transactionWithLabels)
			Expect(err).To(BeNil())
			Expect(result.ID).To(BeNumerically(">", 0))

			count := 0
			err = pgPool.QueryRow(context.TODO(), "select count(*) from transactions_labels;").Scan(&count)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(2))

			// Verify transaction was created with relationships
			retrieved, err := transactionService.GetTransaction(context.TODO(), transactionWithLabels.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Labels).To(HaveLen(2))

			// Verify the relationships include rule references
			for _, assignment := range retrieved.Labels {
				Expect(assignment.RuleID).ToNot(BeNil())
				Expect(*assignment.RuleID).To(BeElementOf([]string{"rule1", "rule2"}))
			}
		})

		It("should create new labels during transaction creation if they don't exist", func() {
			transactionWithNewLabels := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				150.0,
				"Tx with new labels",
			)

			// Add labels that don't exist yet
			transactionWithNewLabels.Labels = []entity.LabelAssociation{
				{
					Label: entity.Label{Key: "new-category", Value: "investment"},
				},
				{
					Label: entity.Label{Key: "new-priority", Value: "high"},
				},
			}

			result, err := transactionService.Create(context.TODO(), *transactionWithNewLabels)
			Expect(err).To(BeNil())
			Expect(result.ID).To(BeNumerically(">", 0))

			// Verify transaction was created and labels were created
			retrieved, err := transactionService.GetTransaction(context.TODO(), transactionWithNewLabels.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Labels).To(HaveLen(2))

			// Verify the new labels were created with proper IDs
			for _, assignment := range retrieved.Labels {
				Expect(assignment.Label.ID).To(BeNumerically(">", 0))
				Expect(assignment.Label.Key).To(BeElementOf([]string{"new-category", "new-priority"}))
			}
		})
	})

	Context("Update", func() {
		It("should return error when transaction does not exist", func() {
			nonExistentTransaction := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-25.0,
				"Non-existent transaction",
			)

			_, err := transactionService.Update(context.TODO(), *nonExistentTransaction)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})

		It("should update existing transaction successfully", func() {
			// First, create a transaction
			originalTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				100.0,
				"Original content",
			)

			created, err := transactionService.Create(context.TODO(), *originalTransaction)
			Expect(err).To(BeNil())
			originalID := created.ID

			// Update the transaction
			updatedTransaction := &entity.Transaction{
				ID:         originalID,
				Hash:       originalTransaction.Hash, // Same hash
				Kind:       entity.CreditTransaction,
				Date:       testTime,
				Account:    1001,
				Amount:     200.0,             // Different amount
				RawContent: "Updated content", // Different content
			}

			result, err := transactionService.Update(context.TODO(), *updatedTransaction)
			Expect(err).To(BeNil())
			Expect(result.ID).To(Equal(originalID)) // Should keep same ID
			Expect(result.Amount).To(Equal(float32(200.0)))

			// Verify transaction was updated
			retrieved, err := transactionService.GetTransaction(context.TODO(), originalTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.ID).To(Equal(originalID))
			Expect(retrieved.Amount).To(Equal(float32(200.0)))
			Expect(retrieved.RawContent).To(Equal("Updated content"))
		})

		It("should update transaction-label relationships correctly", func() {
			// First, create labels and rules
			insertLabelStmt := psql.Insert("labels").Columns("id", "key", "value")
			sql, args, err := insertLabelStmt.
				Values(201, "old-category", "expense").
				Values(202, "old-type", "one-time").
				Values(203, "new-category", "income").
				Values(204, "new-priority", "high").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			insertRuleStmt := psql.Insert("rules").Columns("id", "pattern")
			sql, args, err = insertRuleStmt.
				Values("update-rule-1", "old-pattern").
				Values("update-rule-2", "new-pattern").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Link rules to labels
			insertRuleLabelStmt := psql.Insert("rules_labels").Columns("rule_id", "label_id")
			sql, args, err = insertRuleLabelStmt.
				Values("update-rule-1", 201).
				Values("update-rule-1", 202).
				Values("update-rule-2", 203).
				Values("update-rule-2", 204).
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction with initial labels
			originalTransaction := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-100.0,
				"Transaction for relationship update",
			)

			oldRuleID := "update-rule-1"
			originalTransaction.Labels = []entity.LabelAssociation{
				{
					Label:  entity.Label{ID: 201, Key: "old-category", Value: "expense"},
					RuleID: &oldRuleID,
				},
				{
					Label:  entity.Label{ID: 202, Key: "old-type", Value: "one-time"},
					RuleID: &oldRuleID,
				},
			}

			created, err := transactionService.Create(context.TODO(), *originalTransaction)
			Expect(err).To(BeNil())

			// Update with completely different labels and rules
			newRuleID := "update-rule-2"
			updatedTransaction := &entity.Transaction{
				ID:         created.ID,
				Hash:       originalTransaction.Hash,
				Kind:       entity.DebitTransaction,
				Date:       testTime,
				Account:    1001,
				Amount:     -150.0,
				RawContent: "Updated transaction content",
				Labels: []entity.LabelAssociation{
					{
						Label:  entity.Label{ID: 203, Key: "new-category", Value: "income"},
						RuleID: &newRuleID,
					},
					{
						Label:  entity.Label{ID: 204, Key: "new-priority", Value: "high"},
						RuleID: &newRuleID,
					},
				},
			}

			result, err := transactionService.Update(context.TODO(), *updatedTransaction)
			Expect(err).To(BeNil())
			Expect(result.ID).To(Equal(created.ID))

			// Verify the transaction was updated with new relationships
			retrieved, err := transactionService.GetTransaction(context.TODO(), originalTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Labels).To(HaveLen(2))

			// Verify the new relationships
			labelMap := make(map[string]string)
			for _, assignment := range retrieved.Labels {
				labelMap[assignment.Label.Key] = assignment.Label.Value
				Expect(assignment.RuleID).ToNot(BeNil())
				Expect(*assignment.RuleID).To(Equal("update-rule-2"))
			}
			Expect(labelMap).To(HaveKeyWithValue("new-category", "income"))
			Expect(labelMap).To(HaveKeyWithValue("new-priority", "high"))

			// Verify old relationships are gone
			Expect(labelMap).ToNot(HaveKey("old-category"))
			Expect(labelMap).ToNot(HaveKey("old-type"))
		})

		It("should create new labels during update if they don't exist", func() {
			// First, create a transaction with existing labels
			existingTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				75.0,
				"Transaction to update with new labels",
			)

			// Use existing label
			sql, args, err := psql.Insert("labels").Columns("id", "key", "value").
				Values(301, "existing", "label").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			existingTransaction.Labels = []entity.LabelAssociation{
				{
					Label: entity.Label{ID: 301, Key: "existing", Value: "label"},
				},
			}

			created, err := transactionService.Create(context.TODO(), *existingTransaction)
			Expect(err).To(BeNil())

			// Update with a mix of existing and brand new labels
			updatedTransaction := &entity.Transaction{
				ID:         created.ID,
				Hash:       existingTransaction.Hash,
				Kind:       entity.CreditTransaction,
				Date:       testTime,
				Account:    1001,
				Amount:     100.0,
				RawContent: "Updated with new labels",
				Labels: []entity.LabelAssociation{
					{
						Label: entity.Label{ID: 301, Key: "existing", Value: "label"}, // Keep existing
					},
					{
						Label: entity.Label{Key: "brand-new", Value: "created-on-update"}, // New label
					},
					{
						Label: entity.Label{Key: "another-new", Value: "also-new"}, // Another new label
					},
				},
			}

			result, err := transactionService.Update(context.TODO(), *updatedTransaction)
			Expect(err).To(BeNil())
			Expect(result.ID).To(Equal(created.ID))

			// Verify the transaction has all labels
			retrieved, err := transactionService.GetTransaction(context.TODO(), existingTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Labels).To(HaveLen(3))

			// Verify all labels have proper IDs and values
			labelMap := make(map[string]string)
			for _, assignment := range retrieved.Labels {
				Expect(assignment.Label.ID).To(BeNumerically(">", 0))
				labelMap[assignment.Label.Key] = assignment.Label.Value
			}
			Expect(labelMap).To(HaveKeyWithValue("existing", "label"))
			Expect(labelMap).To(HaveKeyWithValue("brand-new", "created-on-update"))
			Expect(labelMap).To(HaveKeyWithValue("another-new", "also-new"))
		})
	})

	Context("Relationship Management", func() {
		It("should handle complex transaction-rule-label relationships in Create", func() {
			// Setup complex scenario with multiple rules and labels
			insertLabelStmt := psql.Insert("labels").Columns("id", "key", "value")
			sql, args, err := insertLabelStmt.
				Values(401, "category", "shopping").
				Values(402, "merchant", "supermarket").
				Values(403, "payment", "card").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			insertRuleStmt := psql.Insert("rules").Columns("id", "pattern")
			sql, args, err = insertRuleStmt.
				Values("shopping-rule", "supermarket.*shopping").
				Values("payment-rule", "card.*payment").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Link rules to labels
			insertRuleLabelStmt := psql.Insert("rules_labels").Columns("rule_id", "label_id")
			sql, args, err = insertRuleLabelStmt.
				Values("shopping-rule", 401).
				Values("shopping-rule", 402).
				Values("payment-rule", 403).
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Create transaction with multiple rule-label relationships
			complexTransaction := entity.NewTransaction(
				entity.DebitTransaction,
				1001, // account
				testTime,
				-85.50,
				"Complex transaction with multiple relationships",
			)

			shoppingRule := "shopping-rule"
			paymentRule := "payment-rule"
			complexTransaction.Labels = []entity.LabelAssociation{
				{
					Label:  entity.Label{ID: 401, Key: "category", Value: "shopping"},
					RuleID: &shoppingRule,
				},
				{
					Label:  entity.Label{ID: 402, Key: "merchant", Value: "supermarket"},
					RuleID: &shoppingRule,
				},
				{
					Label:  entity.Label{ID: 403, Key: "payment", Value: "card"},
					RuleID: &paymentRule,
				},
			}

			result, err := transactionService.Create(context.TODO(), *complexTransaction)
			Expect(err).To(BeNil())
			Expect(result.ID).To(BeNumerically(">", 0))

			// Verify all relationships were created correctly
			retrieved, err := transactionService.GetTransaction(context.TODO(), complexTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Labels).To(HaveLen(3))

			// Verify each label assignment has the correct rule
			ruleMap := make(map[string]string)
			for _, assignment := range retrieved.Labels {
				Expect(assignment.RuleID).ToNot(BeNil())
				ruleMap[assignment.Label.Key] = *assignment.RuleID
			}
			Expect(ruleMap).To(HaveKeyWithValue("category", "shopping-rule"))
			Expect(ruleMap).To(HaveKeyWithValue("merchant", "shopping-rule"))
			Expect(ruleMap).To(HaveKeyWithValue("payment", "payment-rule"))
		})

		It("should handle transactions with labels but no rules", func() {
			// Create transaction with labels that don't have rule associations
			simpleTransaction := entity.NewTransaction(
				entity.CreditTransaction,
				1001, // account
				testTime,
				45.75,
				"Simple transaction without rules",
			)

			simpleTransaction.Labels = []entity.LabelAssociation{
				{
					Label:  entity.Label{Key: "manual-tag", Value: "user-added"},
					RuleID: nil, // No rule association
				},
				{
					Label:  entity.Label{Key: "custom", Value: "annotation"},
					RuleID: nil, // No rule association
				},
			}

			result, err := transactionService.Create(context.TODO(), *simpleTransaction)
			Expect(err).To(BeNil())
			Expect(result.ID).To(BeNumerically(">", 0))

			// Verify transaction was created with labels but no rule relationships
			retrieved, err := transactionService.GetTransaction(context.TODO(), simpleTransaction.Hash)
			Expect(err).To(BeNil())
			Expect(retrieved).ToNot(BeNil())
			Expect(retrieved.Labels).To(HaveLen(2))

			// Verify no rule associations
			for _, assignment := range retrieved.Labels {
				Expect(assignment.RuleID).To(BeNil())
				Expect(assignment.Label.Key).To(BeElementOf([]string{"manual-tag", "custom"}))
			}
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

			created, err := transactionService.Create(context.TODO(), *transactionToDelete)
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

			result, err := transactionService.Create(context.TODO(), *zeroTransaction)
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

			result, err := transactionService.Create(context.TODO(), *largeTransaction)
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

			result, err := transactionService.Create(context.TODO(), *emptyContentTransaction)
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
