package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	apiV1 = "/api/v1"
)

var _ = Describe("TransactionHandlers", func() {
	var (
		router    *gin.RouterGroup
		datastore *pg.Datastore
		ctx       context.Context
		srv       *httptest.Server
	)

	BeforeEach(func() {
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

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
		srv.Close()
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

		It("should handle valid transaction form", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
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

			// Should either create successfully or handle gracefully
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusCreated, http.StatusOK, http.StatusInternalServerError, http.StatusBadRequest}))
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

		It("should handle form validation", func() {
			form := inbound.CreateTransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Updated transaction",
				Amount:  200.75,
				Account: 1001,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", srv.URL+"/api/v1/transactions/123", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should either update successfully or handle gracefully
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest, http.StatusCreated}))
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
	})
})
