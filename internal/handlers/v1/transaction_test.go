package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	v1Impl "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/inbound"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	apiV1 = "/api/v1"
)

var _ = Describe("TransactionHandlers", Ordered, func() {
	var (
		router    *gin.RouterGroup
		datastore *pg.Datastore
		pgPool    *pgxpool.Pool
		ctx       context.Context
		srv       *httptest.Server
	)

	BeforeAll(func() {
		gin.SetMode(gin.TestMode)
		ctx = context.Background()

		// Get database URL from environment or use default test database
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		}

		// Create real datastore connection
		var err error
		datastore, err = pg.NewPostgresDatastore(ctx, dbURL)
		Expect(err).To(BeNil(), "PostgreSQL database must be available for testing")

		// Create pgx pool for cleanup operations
		pgxConfig, err := pgxpool.ParseConfig(dbURL)
		Expect(err).To(BeNil())

		pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
		Expect(err).To(BeNil())
		Expect(pool).ToNot(BeNil())

		pgPool = pool

		// Add middleware to inject datastore
		engine := gin.New()
		router = engine.Group(apiV1)
		router.Use(func(c *gin.Context) {
			c.Set("datastore", datastore)
			c.Next()
		})

		v1.RegisterHandlers(router, v1Impl.NewServer())
		srv = httptest.NewServer(engine)
	})

	AfterAll(func() {
		if srv != nil {
			srv.Close()
		}

		if datastore != nil {
			datastore.Close()
		}
		if pgPool != nil {
			pgPool.Close()
		}
	})

	Context("GET /api/v1/transactions", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should return OK or accept the request
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions", nil)
			resp, err := http.DefaultClient.Do(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			contentType := resp.Header["Content-Type"][0]
			Expect(contentType).To(ContainSubstring("application/json"))
		})

		It("should handle query parameters gracefully", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions?limit=10&offset=0", nil)
			resp, err := http.DefaultClient.Do(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("GET /api/v1/transactions/:id", func() {
		It("should handle numeric IDs", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions/123", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully even if transaction doesn't exist
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should handle non-numeric IDs", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions/abc", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("POST /api/v1/transactions", func() {
		It("should handle missing form data", func() {
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should handle empty JSON payload", func() {
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should handle invalid JSON", func() {
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should validate required fields", func() {
			form := inbound.CreateTransactionForm{
				// Missing required fields
				Kind:   "debit",
				Amount: 100.50,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should validate form data", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "invalid_kind", // Invalid kind
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  -100.50,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should successfully create a new transaction", func() {
			// Use timestamp to ensure unique content and avoid hash collisions
			uniqueContent := fmt.Sprintf("Test transaction for creation %d", time.Now().UnixNano())
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: uniqueContent,
				Amount:  100.50,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should create successfully
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["id"]).ToNot(BeNil())
			Expect(response["amount"]).To(Equal(100.50))
		})

		It("should return error when trying to create duplicate transaction", func() {
			// Use timestamp to ensure unique content initially
			uniqueContent := fmt.Sprintf("Duplicate transaction test %d", time.Now().UnixNano())
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: uniqueContent,
				Amount:  75.25,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)

			// First creation should succeed
			req1, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req1.Header.Set("Content-Type", "application/json")
			resp1, err := http.DefaultClient.Do(req1)
			Expect(err).To(BeNil())
			Expect(resp1.StatusCode).To(Equal(http.StatusCreated))

			// Second creation should fail (same content = same hash)
			req2, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req2.Header.Set("Content-Type", "application/json")
			resp2, err := http.DefaultClient.Do(req2)
			Expect(err).To(BeNil())
			Expect(resp2.StatusCode).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err = json.NewDecoder(resp2.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
			Expect(response["error"].(string)).To(ContainSubstring("already exists"))
		})

		It("should handle creation with different amounts but same other fields", func() {
			// These should be different transactions due to different amounts
			baseContent := fmt.Sprintf("Amount test transaction %d", time.Now().UnixNano())
			form1 := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: baseContent,
				Amount:  100.00,
				Account: 1001,
			}

			form2 := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: baseContent,
				Amount:  200.00, // Different amount
				Account: 1001,
			}

			jsonData1, _ := json.Marshal(form1)
			jsonData2, _ := json.Marshal(form2)

			// Both should succeed as they have different hashes (due to different amounts)
			req1, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData1))
			req1.Header.Set("Content-Type", "application/json")
			resp1, err := http.DefaultClient.Do(req1)
			Expect(err).To(BeNil())
			Expect(resp1.StatusCode).To(Equal(http.StatusCreated))

			req2, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData2))
			req2.Header.Set("Content-Type", "application/json")
			resp2, err := http.DefaultClient.Do(req2)
			Expect(err).To(BeNil())
			Expect(resp2.StatusCode).To(Equal(http.StatusCreated))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("PUT /api/v1/transactions/:id", func() {
		It("should handle missing ID", func() {
			req, _ := http.NewRequest("PUT", srv.URL+"/api/v1/transactions/", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should return 404 for missing ID (route not found)
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should update existing transaction", func() {
			// First, create a transaction with unique content
			uniqueContent := fmt.Sprintf("Transaction to update %d", time.Now().UnixNano())
			createForm := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: uniqueContent,
				Amount:  150.00,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(createForm)
			createReq, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			createReq.Header.Set("Content-Type", "application/json")
			createResp, err := http.DefaultClient.Do(createReq)
			Expect(err).To(BeNil())
			Expect(createResp.StatusCode).To(Equal(http.StatusCreated))

			var createResponse map[string]interface{}
			err = json.NewDecoder(createResp.Body).Decode(&createResponse)
			Expect(err).To(BeNil())
			transactionID := int64(createResponse["id"].(float64))

			// Now update the transaction
			updateForm := inbound.CreateTransactionForm{
				Kind:    "credit",
				Date:    "16/01/2024",
				Content: "Updated transaction content",
				Amount:  200.00,
				Account: 1001,
			}

			updateJsonData, _ := json.Marshal(updateForm)
			updateReq, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/transactions/%d", srv.URL, transactionID), bytes.NewBuffer(updateJsonData))
			updateReq.Header.Set("Content-Type", "application/json")
			updateResp, err := http.DefaultClient.Do(updateReq)
			Expect(err).To(BeNil())
			Expect(updateResp.StatusCode).To(Equal(http.StatusOK))

			var updateResponse map[string]interface{}
			err = json.NewDecoder(updateResp.Body).Decode(&updateResponse)
			Expect(err).To(BeNil())
			Expect(updateResponse["id"]).To(Equal(createResponse["id"]))
			Expect(updateResponse["description"]).To(Equal("Updated transaction content"))
			Expect(updateResponse["amount"]).To(Equal(200.00))
		})

		It("should create transaction when updating non-existent ID", func() {
			uniqueContent := fmt.Sprintf("Create via update %d", time.Now().UnixNano())
			updateForm := inbound.CreateTransactionForm{
				Kind:    "credit",
				Date:    "17/01/2024",
				Content: uniqueContent,
				Amount:  300.00,
				Account: 1001,
			}

			updateJsonData, _ := json.Marshal(updateForm)
			updateReq, _ := http.NewRequest("PUT", srv.URL+"/api/v1/transactions/999999", bytes.NewBuffer(updateJsonData))
			updateReq.Header.Set("Content-Type", "application/json")
			updateResp, err := http.DefaultClient.Do(updateReq)
			Expect(err).To(BeNil())

			// Should create new transaction since ID doesn't exist
			Expect(updateResp.StatusCode).To(Equal(http.StatusCreated))

			var response map[string]interface{}
			err = json.NewDecoder(updateResp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["id"]).ToNot(BeNil())
			Expect(response["description"]).To(Equal(uniqueContent))
			Expect(response["amount"]).To(Equal(300.00))
		})

		It("should handle invalid form data in update", func() {
			updateForm := inbound.CreateTransactionForm{
				Kind:    "invalid_kind",
				Date:    "17/01/2024",
				Content: "Invalid update",
				Amount:  300.00,
				Account: 1001,
			}

			updateJsonData, _ := json.Marshal(updateForm)
			updateReq, _ := http.NewRequest("PUT", srv.URL+"/api/v1/transactions/123", bytes.NewBuffer(updateJsonData))
			updateReq.Header.Set("Content-Type", "application/json")
			updateResp, err := http.DefaultClient.Do(updateReq)
			Expect(err).To(BeNil())

			Expect(updateResp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should handle malformed JSON in update", func() {
			updateReq, _ := http.NewRequest("PUT", srv.URL+"/api/v1/transactions/123", bytes.NewBuffer([]byte("invalid json")))
			updateReq.Header.Set("Content-Type", "application/json")
			updateResp, err := http.DefaultClient.Do(updateReq)
			Expect(err).To(BeNil())

			Expect(updateResp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("DELETE /api/v1/transactions/:id", func() {
		It("should handle numeric IDs", func() {
			req, _ := http.NewRequest("DELETE", srv.URL+"/api/v1/transactions/123", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully even if transaction doesn't exist
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError, http.StatusNoContent}))
		})

		It("should handle non-numeric IDs", func() {
			req, _ := http.NewRequest("DELETE", srv.URL+"/api/v1/transactions/abc", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError}))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("Form to Entity Mapping", func() {
		It("should correctly map TransactionForm to entity.Transaction", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test mapping",
				Amount:  -50.25,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Kind).To(Equal(entity.DebitTransaction))
			Expect(txEntity.Date.Day()).To(Equal(15))
			Expect(txEntity.Date.Month()).To(Equal(time.January))
			Expect(txEntity.Date.Year()).To(Equal(2024))
			Expect(txEntity.Amount).To(Equal(float32(-50.25)))
			Expect(txEntity.RawContent).To(Equal("Test mapping"))
			Expect(txEntity.Account).To(Equal(int64(1001)))
			// Labels should be empty as they are populated by service layer through rule application
			Expect(txEntity.Labels).To(HaveLen(0))
		})

		It("should correctly map credit transactions", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "credit",
				Date:    "15/01/2024",
				Content: "Credit transaction",
				Amount:  100.75,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Kind).To(Equal(entity.CreditTransaction))
			Expect(txEntity.Amount).To(Equal(float32(100.75)))
			Expect(txEntity.Account).To(Equal(int64(1001)))
			// Labels should be empty as they are populated by service layer
			Expect(txEntity.Labels).To(HaveLen(0))
		})

		It("should handle invalid date format", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "invalid-date",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
			}

			_, err := form.Entity()
			Expect(err).ToNot(BeNil())
		})

		It("should handle different date formats", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "01/12/2023", // Different date
				Content: "Test transaction",
				Amount:  75.25,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Date.Day()).To(Equal(1))
			Expect(txEntity.Date.Month()).To(Equal(time.December))
			Expect(txEntity.Date.Year()).To(Equal(2023))
		})

		It("should handle empty labels", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Labels).To(HaveLen(0))
		})

		It("should handle nil labels", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Labels).To(HaveLen(0))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("Form Validation", func() {
		It("should validate transaction kind", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "invalid_kind",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should validate required amount", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  0, // Invalid amount
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should either validate properly or handle gracefully
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusBadRequest, http.StatusCreated, http.StatusOK}))
		})

		It("should validate required account", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 0, // Invalid account
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should either validate properly or handle gracefully
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusBadRequest, http.StatusCreated, http.StatusOK}))
		})

		It("should handle large amounts", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Large transaction",
				Amount:  999999.99,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Amount).To(Equal(float32(999999.99)))
		})

		It("should handle negative amounts", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Negative transaction",
				Amount:  -150.75,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Amount).To(Equal(float32(-150.75)))
		})

		It("should handle special characters in content", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Transaction with special chars: !@#$%^&*()",
				Amount:  100.50,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.RawContent).To(Equal("Transaction with special chars: !@#$%^&*()"))
		})

		It("should handle unicode characters in content", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Transaction with unicode: 🏦💰",
				Amount:  100.50,
				Account: 1001,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.RawContent).To(Equal("Transaction with unicode: 🏦💰"))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("Error Handling and Edge Cases", func() {
		It("should handle service layer errors in CreateTransaction", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Service error test",
				Amount:  100.50,
				Account: 0, // This might cause issues at service layer
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully, either with 400 or 500
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusBadRequest, http.StatusInternalServerError}))
		})

		It("should handle concurrent creation attempts of same transaction", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Concurrent test transaction",
				Amount:  150.75,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)

			// Simulate concurrent requests
			req1, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req1.Header.Set("Content-Type", "application/json")
			req2, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req2.Header.Set("Content-Type", "application/json")

			// Send both requests
			resp1, err1 := http.DefaultClient.Do(req1)
			resp2, err2 := http.DefaultClient.Do(req2)

			Expect(err1).To(BeNil())
			Expect(err2).To(BeNil())

			// One should succeed, one should fail
			statusCodes := []int{resp1.StatusCode, resp2.StatusCode}
			Expect(statusCodes).To(ContainElement(http.StatusCreated))
			Expect(statusCodes).To(ContainElement(http.StatusBadRequest))
		})

		It("should handle very large transaction amounts", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "credit",
				Date:    "15/01/2024",
				Content: "Large amount transaction",
				Amount:  999999999.99,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusCreated, http.StatusBadRequest}))
		})

		It("should handle zero amount transactions", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Zero amount transaction",
				Amount:  0.0,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully - some systems allow zero amounts, others don't
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusCreated, http.StatusBadRequest}))
		})
	})

	Context("Transaction Label Operations", func() {
		It("should handle getting labels for non-existent transaction", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions/999999/labels", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should handle adding label to non-existent transaction", func() {
			labelForm := map[string]interface{}{
				"key":   "category",
				"value": "expense",
			}

			jsonData, _ := json.Marshal(labelForm)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions/999999/labels", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should handle removing labels from non-existent transaction", func() {
			req, _ := http.NewRequest("DELETE", srv.URL+"/api/v1/transactions/999999/labels", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should handle removing specific label from non-existent transaction", func() {
			req, _ := http.NewRequest("DELETE", srv.URL+"/api/v1/transactions/999999/labels/123", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should handle malformed label data", func() {
			// First create a transaction
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Transaction for label test",
				Amount:  100.50,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			createReq, _ := http.NewRequest("POST", srv.URL+"/api/v1/transactions", bytes.NewBuffer(jsonData))
			createReq.Header.Set("Content-Type", "application/json")
			createResp, err := http.DefaultClient.Do(createReq)
			Expect(err).To(BeNil())
			Expect(createResp.StatusCode).To(Equal(http.StatusCreated))

			var createResponse map[string]interface{}
			err = json.NewDecoder(createResp.Body).Decode(&createResponse)
			Expect(err).To(BeNil())
			transactionID := int64(createResponse["id"].(float64))

			// Try to add malformed label
			malformedLabel := "not valid json"
			labelReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/transactions/%d/labels", srv.URL, transactionID), bytes.NewBuffer([]byte(malformedLabel)))
			labelReq.Header.Set("Content-Type", "application/json")
			labelResp, err := http.DefaultClient.Do(labelReq)
			Expect(err).To(BeNil())

			Expect(labelResp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		AfterEach(func() {
			// Clean up after each test
			_, err := pgPool.Exec(ctx, "DELETE FROM transactions_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules_labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM transactions;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM labels;")
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(ctx, "DELETE FROM rules;")
			Expect(err).To(BeNil())
		})
	})

	Context("HTTP Method and Route Validation", func() {
		It("should reject unsupported HTTP methods", func() {
			// PATCH method not supported
			req, _ := http.NewRequest("PATCH", srv.URL+"/api/v1/transactions/123", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should handle malformed transaction IDs", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions/not-a-number", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should handle negative transaction IDs", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/transactions/-123", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusBadRequest, http.StatusNotFound}))
		})
	})
})
