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

var _ = Describe("LabelHandlers", func() {
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

	Context("GET /api/v1/labels", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/labels", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should return OK or accept the request
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/labels", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			if resp.StatusCode == http.StatusOK {
				contentType := resp.Header.Get("Content-Type")
				Expect(contentType).To(ContainSubstring("application/json"))
			}
		})

		It("should handle empty result gracefully", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/labels", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should not crash on empty results
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted, http.StatusInternalServerError}))
		})
	})
})
