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

var _ = Describe("TagHandlers", func() {
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

		// Register tag handlers
		api := router.Group("/api/v1")
		v1.TagHandlers(api)
	})

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
	})

	Context("GET /api/v1/tags", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("GET", "/api/v1/tags", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should return OK or accept the request
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", "/api/v1/tags", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				Expect(contentType).To(ContainSubstring("application/json"))
			}
		})
	})

	Context("POST /api/v1/tags", func() {
		It("should return error for invalid JSON", func() {
			req, _ := http.NewRequest("POST", "/api/v1/tags", bytes.NewBuffer([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate tag form", func() {
			form := inbound.TagForm{
				Value: "", // Empty value should fail validation
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/tags", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should accept valid tag form", func() {
			form := inbound.TagForm{
				Value: "test_tag",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/tags", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not return bad request for valid form
			Expect(w.Code).ToNot(Equal(http.StatusBadRequest))
		})
	})

	Context("DELETE /api/v1/tags/:id", func() {
		It("should accept string tag values", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/tags/test_tag", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not return bad request for valid tag value
			Expect(w.Code).ToNot(Equal(http.StatusBadRequest))
		})

		It("should handle special characters in tag values", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/tags/special-tag_123", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not return bad request for valid tag value with special characters
			Expect(w.Code).ToNot(Equal(http.StatusBadRequest))
		})

		It("should handle empty tag value", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/tags/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 404 for empty tag value (route not found)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Form Validation and JSON Handling", func() {
		It("should correctly marshal and unmarshal TagForm", func() {
			originalForm := inbound.TagForm{
				Value: "test_tag_value",
			}

			jsonData, err := json.Marshal(originalForm)
			Expect(err).To(BeNil())
			Expect(jsonData).ToNot(BeEmpty())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("test_tag_value"))
		})

		It("should handle special characters in tag values", func() {
			form := inbound.TagForm{
				Value: "special-tag_123!@#$%^&*()",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("special-tag_123!@#$%^&*()"))
		})

		It("should handle Unicode characters in tag values", func() {
			form := inbound.TagForm{
				Value: "tag_with_unicode_🏷️_characters",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("tag_with_unicode_🏷️_characters"))
		})

		It("should handle empty string values", func() {
			form := inbound.TagForm{
				Value: "",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal(""))
		})

		It("should handle whitespace-only values", func() {
			form := inbound.TagForm{
				Value: "   \t\n   ",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("   \t\n   "))
		})

		It("should handle very long tag values", func() {
			longValue := ""
			for i := 0; i < 1000; i++ {
				longValue += "a"
			}

			form := inbound.TagForm{
				Value: longValue,
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal(longValue))
		})

		It("should handle numeric-like string values", func() {
			form := inbound.TagForm{
				Value: "12345.67",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("12345.67"))
		})

		It("should handle boolean-like string values", func() {
			form := inbound.TagForm{
				Value: "true",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("true"))
		})

		It("should handle null-like string values", func() {
			form := inbound.TagForm{
				Value: "null",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("null"))
		})

		It("should handle path-like string values", func() {
			form := inbound.TagForm{
				Value: "/path/to/some/resource",
			}

			jsonData, err := json.Marshal(form)
			Expect(err).To(BeNil())

			var unmarshaledForm inbound.TagForm
			err = json.Unmarshal(jsonData, &unmarshaledForm)
			Expect(err).To(BeNil())
			Expect(unmarshaledForm.Value).To(Equal("/path/to/some/resource"))
		})
	})
})

// TestTagHandlers is handled by the main handlers_suite_test.go
