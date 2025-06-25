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

		It("should parse query parameters correctly", func() {
			req, _ := http.NewRequest("GET", "/api/v1/transactions?start=01/01/2024&end=31/01/2024", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should process the query parameters without crashing
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle invalid query parameters gracefully", func() {
			req, _ := http.NewRequest("GET", "/api/v1/transactions?start=invalid-date&end=invalid-date", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should handle invalid dates gracefully (logs warning but continues)
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
	})

	Context("POST /api/v1/transactions", func() {
		It("should return error for invalid JSON", func() {
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate form data", func() {
			form := inbound.TransactionForm{
				Kind:    "invalid_kind", // Invalid kind
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  -100.50,
				Tags:    map[string]string{"category": "rule1"},
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

		It("should validate required fields", func() {
			form := inbound.TransactionForm{
				// Missing required fields
				Kind: "debit",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate date format", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "invalid-date-format",
				Content: "Test transaction",
				Amount:  -100.50,
				Tags:    map[string]string{"category": "rule1"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate tags", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  -100.50,
				Tags:    map[string]string{}, // Empty tags should fail validation
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("PUT /api/v1/transactions/:id", func() {
		It("should return error for invalid ID", func() {
			form := inbound.TransactionForm{
				Kind:    "credit",
				Date:    "15/01/2024",
				Content: "Updated transaction",
				Amount:  200.75,
				Tags:    map[string]string{"type": "rule2"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", "/api/v1/transactions/invalid_id", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).To(Equal("id must be an int"))
		})

		It("should validate form data for updates", func() {
			form := inbound.TransactionForm{
				Kind:    "invalid_kind", // Invalid kind
				Date:    "15/01/2024",
				Content: "Updated transaction",
				Amount:  200.75,
				Tags:    map[string]string{"type": "rule2"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", "/api/v1/transactions/1", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("DELETE /api/v1/transactions/:id", func() {
		It("should return error for invalid ID", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/transactions/invalid_id", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).To(Equal("id must be an int"))
		})

		It("should accept valid ID format", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/transactions/1", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not return bad request for valid ID format
			Expect(w.Code).ToNot(Equal(http.StatusBadRequest))
		})
	})

	Context("Form to Entity Mapping", func() {
		It("should correctly map TransactionForm to entity.Transaction", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test mapping",
				Amount:  -50.25,
				Tags:    map[string]string{"category": "rule1", "type": "rule2"},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Kind).To(Equal(entity.DebitTransaction))
			Expect(txEntity.Date.Day()).To(Equal(15))
			Expect(txEntity.Date.Month()).To(Equal(time.January))
			Expect(txEntity.Date.Year()).To(Equal(2024))
			Expect(txEntity.Amount).To(Equal(float32(-50.25)))
			Expect(txEntity.RawContent).To(Equal("Test mapping"))
			Expect(txEntity.Tags).To(HaveLen(2))
			Expect(txEntity.Tags["category"]).To(Equal("rule1"))
			Expect(txEntity.Tags["type"]).To(Equal("rule2"))
		})

		It("should correctly map credit transactions", func() {
			form := inbound.TransactionForm{
				Kind:    "credit",
				Date:    "15/01/2024",
				Content: "Credit transaction",
				Amount:  100.75,
				Tags:    map[string]string{"income": "salary"},
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Kind).To(Equal(entity.CreditTransaction))
			Expect(txEntity.Amount).To(Equal(float32(100.75)))
			Expect(txEntity.RawContent).To(Equal("Credit transaction"))
		})

		It("should return error for invalid date format", func() {
			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "invalid-date",
				Content: "Test mapping",
				Amount:  -50.25,
				Tags:    map[string]string{"category": "rule1"},
			}

			_, err := form.Entity()
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("unable to parse transaction date"))
		})

		It("should preserve tags in entity mapping", func() {
			tags := map[string]string{
				"category": "rule1",
				"type":     "rule2",
				"priority": "rule3",
			}

			form := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Multi-tag transaction",
				Amount:  -25.50,
				Tags:    tags,
			}

			txEntity, err := form.Entity()
			Expect(err).To(BeNil())
			Expect(txEntity.Tags).To(Equal(tags))
		})

		It("should generate correct hash for transactions", func() {
			form1 := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  -50.25,
				Tags:    map[string]string{"category": "rule1"},
			}

			form2 := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction",
				Amount:  -50.25,
				Tags:    map[string]string{"category": "rule1"},
			}

			entity1, err1 := form1.Entity()
			entity2, err2 := form2.Entity()

			Expect(err1).To(BeNil())
			Expect(err2).To(BeNil())
			Expect(entity1.Hash).To(Equal(entity2.Hash)) // Same data should produce same hash
			Expect(entity1.Hash).ToNot(BeEmpty())
		})

		It("should generate different hashes for different transactions", func() {
			form1 := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction 1",
				Amount:  -50.25,
				Tags:    map[string]string{"category": "rule1"},
			}

			form2 := inbound.TransactionForm{
				Kind:    "debit",
				Date:    "15/01/2024",
				Content: "Test transaction 2", // Different content
				Amount:  -50.25,
				Tags:    map[string]string{"category": "rule1"},
			}

			entity1, err1 := form1.Entity()
			entity2, err2 := form2.Entity()

			Expect(err1).To(BeNil())
			Expect(err2).To(BeNil())
			Expect(entity1.Hash).ToNot(Equal(entity2.Hash)) // Different data should produce different hashes
		})
	})
})

// TestTransactionHandlers is handled by the main handlers_suite_test.go
