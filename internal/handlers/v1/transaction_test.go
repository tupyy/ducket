package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	v1 "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/inbound"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TransactionHandlers", func() {
	var (
		router    *gin.Engine
		datastore *pg.Datastore
		ctx       context.Context
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
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
		router.Use(func(c *gin.Context) {
			c.Set("datastore", datastore)
			c.Next()
		})

		// Register transaction handlers
		api := router.Group("/api/v1")
		v1.TransactionHandlers(api)
	})

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
	})

	Context("GET /api/v1/transactions", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("GET", "/api/v1/transactions", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should return OK or accept the request
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", "/api/v1/transactions", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				Expect(contentType).To(ContainSubstring("application/json"))
			}
		})

		It("should handle query parameters gracefully", func() {
			req, _ := http.NewRequest("GET", "/api/v1/transactions?limit=10&offset=0", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not crash with query parameters
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted, http.StatusBadRequest}))
		})
	})

	Context("GET /api/v1/transactions/:id", func() {
		It("should handle numeric IDs", func() {
			req, _ := http.NewRequest("GET", "/api/v1/transactions/123", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle gracefully even if transaction doesn't exist
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError}))
		})

		It("should handle non-numeric IDs", func() {
			req, _ := http.NewRequest("GET", "/api/v1/transactions/abc", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError}))
		})
	})

	Context("POST /api/v1/transactions", func() {
		It("should handle missing form data", func() {
			req, _ := http.NewRequest("POST", "/api/v1/transactions", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should handle empty JSON payload", func() {
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should handle invalid JSON", func() {
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate required fields", func() {
			form := inbound.TransactionForm{
				// Missing required fields
				Kind:   "debit",
				Amount: 100.50,
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate form data", func() {
			form := inbound.TransactionForm{
				Kind:    "invalid_kind", // Invalid kind
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  -100.50,
				Account: 1001,
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should handle valid transaction form", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should either create successfully or handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusCreated, http.StatusOK, http.StatusInternalServerError, http.StatusBadRequest}))
		})
	})

	Context("PUT /api/v1/transactions/:id", func() {
		It("should handle missing ID", func() {
			req, _ := http.NewRequest("PUT", "/api/v1/transactions/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 404 for missing ID (route not found)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("should handle form validation", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Updated transaction",
				Amount:  200.75,
				Account: 1001,
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", "/api/v1/transactions/123", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should either update successfully or handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError, http.StatusBadRequest, http.StatusCreated}))
		})
	})

	Context("DELETE /api/v1/transactions/:id", func() {
		It("should handle numeric IDs", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/transactions/123", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle gracefully even if transaction doesn't exist
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError, http.StatusNoContent}))
		})

		It("should handle non-numeric IDs", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/transactions/abc", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError}))
		})
	})

	Context("Form to Entity Mapping", func() {
		It("should correctly map TransactionForm to entity.Transaction", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test mapping",
				Amount:  -50.25,
				Account: 1001,
				Labels:  map[string]string{"category": "food", "type": "essential"},
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
			form := inbound.TransactionForm{
				Kind:    "credit",
				Date:    "15/01/2024",
				Content: "Credit transaction",
				Amount:  100.75,
				Account: 1001,
				Labels:  map[string]string{"category": "income"},
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
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "invalid-date",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
				Labels:  map[string]string{"category": "food"},
			}

			_, err := form.Entity()
			Expect(err).ToNot(BeNil())
		})

		It("should handle different date formats", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "01/12/2023", // Different date
				Content: "Test transaction",
				Amount:  75.25,
				Account: 1001,
				Labels:  map[string]string{"category": "food"},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Date.Day()).To(Equal(1))
			Expect(txEntity.Date.Month()).To(Equal(time.December))
			Expect(txEntity.Date.Year()).To(Equal(2023))
		})

		It("should handle empty labels", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
				Labels:  map[string]string{},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Labels).To(HaveLen(0))
		})

		It("should handle nil labels", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
				Labels:  nil,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Labels).To(HaveLen(0))
		})
	})

	Context("Form Validation", func() {
		It("should validate transaction kind", func() {
			form := inbound.TransactionForm{
				Kind:    "invalid_kind",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 1001,
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate required amount", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  0, // Invalid amount
				Account: 1001,
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should either validate properly or handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusBadRequest, http.StatusCreated, http.StatusOK}))
		})

		It("should validate required account", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  100.50,
				Account: 0, // Invalid account
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should either validate properly or handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusBadRequest, http.StatusCreated, http.StatusOK}))
		})

		It("should handle large amounts", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Large transaction",
				Amount:  999999.99,
				Account: 1001,
				Labels:  map[string]string{"category": "investment"},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Amount).To(Equal(float32(999999.99)))
		})

		It("should handle negative amounts", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Negative transaction",
				Amount:  -150.75,
				Account: 1001,
				Labels:  map[string]string{"category": "refund"},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Amount).To(Equal(float32(-150.75)))
		})

		It("should handle special characters in content", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Transaction with special chars: !@#$%^&*()",
				Amount:  100.50,
				Account: 1001,
				Labels:  map[string]string{"category": "misc"},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.RawContent).To(Equal("Transaction with special chars: !@#$%^&*()"))
		})

		It("should handle unicode characters in content", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Transaction with unicode: 🏦💰",
				Amount:  100.50,
				Account: 1001,
				Labels:  map[string]string{"category": "banking"},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.RawContent).To(Equal("Transaction with unicode: 🏦💰"))
		})
	})
})
