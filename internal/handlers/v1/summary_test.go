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

		It("should handle valid date parameters", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary?start=01/01/2024&end=31/01/2024", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should process valid date parameters
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle invalid date parameters gracefully", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary?start=invalid-date&end=invalid-date", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should handle invalid dates gracefully (logs warning but continues)
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle partial date parameters", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary?start=01/01/2024", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should handle partial date parameters
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle malformed query parameters", func() {
			req, _ := http.NewRequest("GET", "/api/v1/summary?start=&end=", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should handle empty query parameters
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
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
		It("should handle various date formats gracefully", func() {
			dateFormats := []string{
				"01/01/2024",
				"1/1/2024",
				"2024-01-01",
				"01-01-2024",
				"2024/01/01",
			}

			for _, dateFormat := range dateFormats {
				req, _ := http.NewRequest("GET", "/api/v1/summary?start="+dateFormat, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Should not crash regardless of date format
				Expect(w.Code).ToNot(Equal(http.StatusInternalServerError))
			}
		})

		It("should handle edge case dates", func() {
			edgeCases := []string{
				"29/02/2024", // Leap year
				"31/12/2023", // End of year
				"01/01/2000", // Y2K
				"",           // Empty string
				"not-a-date", // Invalid format
			}

			for _, edgeCase := range edgeCases {
				req, _ := http.NewRequest("GET", "/api/v1/summary?start="+edgeCase, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Should handle edge cases gracefully
				Expect(w.Code).ToNot(Equal(http.StatusInternalServerError))
			}
		})
	})
})

// TestSummaryHandlers is handled by the main handlers_suite_test.go
