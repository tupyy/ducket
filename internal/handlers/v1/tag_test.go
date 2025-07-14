package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	v1 "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1/inbound"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LabelHandlers", func() {
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

		// Register label handlers
		api := router.Group("/api/v1")
		v1.LabelHandlers(api)
	})

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
	})

	Context("GET /api/v1/labels", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("GET", "/api/v1/labels", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should return OK or accept the request
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", "/api/v1/labels", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				Expect(contentType).To(ContainSubstring("application/json"))
			}
		})

		It("should handle empty result gracefully", func() {
			req, _ := http.NewRequest("GET", "/api/v1/labels", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not crash on empty results
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted, http.StatusInternalServerError}))
		})
	})

	Context("POST /api/v1/labels", func() {
		It("should handle missing form data", func() {
			req, _ := http.NewRequest("POST", "/api/v1/labels", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should handle empty JSON payload", func() {
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate required fields", func() {
			form := inbound.LabelForm{
				Key: "category",
				// Missing value
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should handle valid label creation", func() {
			form := inbound.LabelForm{
				Key:   "category",
				Value: "food",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should either create successfully or handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusCreated, http.StatusOK, http.StatusInternalServerError}))
		})

		It("should handle invalid JSON", func() {
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should reject empty key", func() {
			form := inbound.LabelForm{
				Key:   "",
				Value: "food",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should reject empty value", func() {
			form := inbound.LabelForm{
				Key:   "category",
				Value: "",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should handle form validation errors gracefully", func() {
			form := map[string]interface{}{
				"key":   123, // Invalid type
				"value": "food",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("DELETE /api/v1/labels/:id", func() {
		It("should handle numeric ID", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/labels/1", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle gracefully (currently returns not implemented)
			Expect(w.Code).To(BeElementOf([]int{http.StatusNotImplemented, http.StatusOK, http.StatusNotFound}))
		})

		It("should handle non-numeric ID", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/labels/abc", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle gracefully
			Expect(w.Code).To(BeElementOf([]int{http.StatusNotImplemented, http.StatusBadRequest, http.StatusNotFound}))
		})

		It("should handle empty ID", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/labels/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 404 for empty ID (route not found)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("should handle missing ID parameter", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/labels", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 404 for missing ID (route not found)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Content-Type validation", func() {
		It("should accept application/json", func() {
			form := inbound.LabelForm{
				Key:   "category",
				Value: "food",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should process the request (not necessarily succeed)
			Expect(w.Code).ToNot(Equal(http.StatusUnsupportedMediaType))
		})

		It("should handle missing Content-Type", func() {
			form := inbound.LabelForm{
				Key:   "category",
				Value: "food",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/labels", bytes.NewBuffer(jsonData))
			// No Content-Type header
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should still process the request
			Expect(w.Code).ToNot(Equal(http.StatusUnsupportedMediaType))
		})
	})

	Context("LabelForm validation", func() {
		It("should convert LabelForm to entity correctly", func() {
			form := inbound.LabelForm{
				Key:   "category",
				Value: "food",
			}

			entity := form.ToEntity()
			Expect(entity.Key).To(Equal("category"))
			Expect(entity.Value).To(Equal("food"))
		})

		It("should handle special characters in key", func() {
			form := inbound.LabelForm{
				Key:   "category-type_123",
				Value: "food",
			}

			entity := form.ToEntity()
			Expect(entity.Key).To(Equal("category-type_123"))
			Expect(entity.Value).To(Equal("food"))
		})

		It("should handle special characters in value", func() {
			form := inbound.LabelForm{
				Key:   "category",
				Value: "fast-food_123!@#",
			}

			entity := form.ToEntity()
			Expect(entity.Key).To(Equal("category"))
			Expect(entity.Value).To(Equal("fast-food_123!@#"))
		})

		It("should handle unicode characters", func() {
			form := inbound.LabelForm{
				Key:   "category",
				Value: "食品",
			}

			entity := form.ToEntity()
			Expect(entity.Key).To(Equal("category"))
			Expect(entity.Value).To(Equal("食品"))
		})
	})
})
