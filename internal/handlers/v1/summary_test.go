package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	v1 "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SummaryHandlers", func() {
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

		// Register summary handlers
		api := router.Group("/api/v1")
		v1.SummaryHandlers(api)
	})

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
	})

	Context("GET /api/v1/summary", func() {
		It("should handle requests without query parameters", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should return OK or accept the request
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle valid timestamp parameters", func() {
			// Using timestamps: 1704067200000 = 2024-01-01, 1706745600000 = 2024-02-01
			req, _ := http.NewRequest("GET", "/api/v1/summary?startDate=1704067200000&endDate=1706745600000", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should process valid timestamp parameters
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle invalid timestamp parameters gracefully", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary?startDate=invalid-timestamp&endDate=invalid-timestamp", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should handle invalid timestamps gracefully (logs warning but continues)
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle partial timestamp parameters", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary?startDate=1704067200000", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should handle partial timestamp parameters
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle empty query parameters", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary?startDate=&endDate=", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should handle empty query parameters
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return 400 when startDate is after endDate", func() {
			// startDate: 1706745600000 = 2024-02-01, endDate: 1704067200000 = 2024-01-01
			req, _ := http.NewRequest("GET", "/api/v1/summary?startDate=1706745600000&endDate=1704067200000", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should return 400 Bad Request for invalid date range
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				Expect(contentType).To(ContainSubstring("application/json"))
			}
		})
	})

	Context("Date Parameter Parsing", func() {
		It("should handle various timestamp formats gracefully", func() {
			timestampFormats := []string{
				"1704067200000", // Valid timestamp
				"1706745600000", // Valid timestamp
				"0",             // Zero timestamp (epoch)
				"-1000000000",   // Negative timestamp
				"9999999999999", // Large timestamp
			}

			for _, timestampFormat := range timestampFormats {
				req, _ := http.NewRequest("GET", "/api/v1/summary?startDate="+timestampFormat, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Should not crash regardless of timestamp format
				Expect(w.Code).ToNot(Equal(http.StatusInternalServerError))
			}
		})

		It("should handle edge case timestamps", func() {
			edgeCases := []string{
				"1704067200000",   // Valid timestamp
				"",                // Empty string
				"not-a-timestamp", // Invalid format
				"abc123",          // Non-numeric
				"123.456",         // Decimal number
			}

			for _, edgeCase := range edgeCases {
				req, _ := http.NewRequest("GET", "/api/v1/summary?startDate="+edgeCase, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Should handle edge cases gracefully
				Expect(w.Code).ToNot(Equal(http.StatusInternalServerError))
			}
		})
	})
})

// TestSummaryHandlers is handled by the main handlers_suite_test.go
