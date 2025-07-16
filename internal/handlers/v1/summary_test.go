package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	v1Impl "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SummaryHandlers", func() {
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
		router = engine.Group("/api/v1")
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

	Context("GET /api/v1/summary", func() {
		It("should handle requests without query parameters", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should return OK or accept the request
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle valid timestamp parameters", func() {
			// Using timestamps: 1704067200000 = 2024-01-01, 1706745600000 = 2024-02-01
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary?startDate=1704067200000&endDate=1706745600000", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should process valid timestamp parameters
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle invalid timestamp parameters gracefully", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary?startDate=invalid-timestamp&endDate=invalid-timestamp", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should handle invalid timestamps gracefully (logs warning but continues)
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle partial timestamp parameters", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary?startDate=1704067200000", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should handle partial timestamp parameters
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should handle empty query parameters", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary?startDate=&endDate=", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should handle empty query parameters
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return 400 when startDate is after endDate", func() {
			// startDate: 1706745600000 = 2024-02-01, endDate: 1704067200000 = 2024-01-01
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary?startDate=1706745600000&endDate=1704067200000", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should return 400 Bad Request for invalid date range
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			if resp.StatusCode == http.StatusOK {
				contentType := resp.Header.Get("Content-Type")
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
				req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary?startDate="+timestampFormat, nil)
				resp, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())

				// Should not crash regardless of timestamp format
				Expect(resp.StatusCode).ToNot(Equal(http.StatusInternalServerError))
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
				req, _ := http.NewRequest("GET", srv.URL+"/api/v1/summary?startDate="+edgeCase, nil)
				resp, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())

				// Should handle edge cases gracefully
				Expect(resp.StatusCode).ToNot(Equal(http.StatusInternalServerError))
			}
		})
	})
})

// TestSummaryHandlers is handled by the main handlers_suite_test.go
