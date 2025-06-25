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
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuleHandlers", func() {
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

		// Register rule handlers
		api := router.Group("/api/v1")
		v1.RulesHandlers(api)
	})

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
	})

	Context("GET /api/v1/rules", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("GET", "/api/v1/rules", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// The handler should return OK or accept the request
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", "/api/v1/rules", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				Expect(contentType).To(ContainSubstring("application/json"))
			}
		})
	})

	Context("POST /api/v1/rules", func() {
		It("should return error for invalid JSON", func() {
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate rule name length", func() {
			form := inbound.RuleForm{
				Name:    "this_is_a_very_long_rule_name_that_exceeds_the_twenty_character_limit",
				Pattern: ".*test.*",
				Tags:    []string{"category"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate regex pattern", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: "[invalid_regex", // Invalid regex pattern
				Tags:    []string{"category"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate required fields", func() {
			form := inbound.RuleForm{
				// Missing required fields
				Name: "test_rule",
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate tag name", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Tags:    []string{}, // Empty tags should fail validation
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("PUT /api/v1/rules/:id", func() {
		It("should create rule even if the rule does not exists", func() {
			form := inbound.UpdateRuleForm{
				Pattern: ".*updated.*",
				Tags:    []string{"category"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", "/api/v1/rules/missing_id", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should validate form data for updates", func() {
			form := inbound.UpdateRuleForm{
				Pattern: ".*updated.*",
				Tags:    []string{}, // Empty tags should fail validation
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", "/api/v1/rules/test_rule", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("DELETE /api/v1/rules/:id", func() {
		It("should accept string IDs", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/rules/test_rule", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not return bad request for valid string ID format
			Expect(w.Code).ToNot(Equal(http.StatusBadRequest))
		})

		It("should handle empty ID", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/rules/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 404 for empty ID (route not found)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Form to Entity Mapping", func() {
		It("should correctly map RuleForm to entity.Rule", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Tags:    []string{"category"},
			}

			ruleEntity := inbound.FormToEntity(form)
			Expect(ruleEntity.Name).To(Equal("test_rule"))
			Expect(ruleEntity.Pattern).To(Equal(".*test.*"))
			Expect(ruleEntity.Tags).To(HaveLen(1))
			Expect(ruleEntity.Tags).To(ContainElement("category"))
		})

		It("should validate regex during form validation", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: "[invalid_regex", // Invalid regex
				Tags:    []string{"category"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.RuleFormValidation, inbound.RuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})

		It("should validate rule name length during form validation", func() {
			form := inbound.RuleForm{
				Name:    "this_is_a_very_long_rule_name_that_exceeds_the_twenty_character_limit",
				Pattern: ".*test.*",
				Tags:    []string{"category"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.RuleFormValidation, inbound.RuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})

		It("should preserve all fields in entity mapping", func() {
			form := inbound.RuleForm{
				Name:    "complex_rule",
				Pattern: "^[A-Z]+\\s+\\d+$",
				Tags:    []string{"transaction_type", "payment"},
			}

			ruleEntity := inbound.FormToEntity(form)
			Expect(ruleEntity.Name).To(Equal("complex_rule"))
			Expect(ruleEntity.Pattern).To(Equal("^[A-Z]+\\s+\\d+$"))
			Expect(ruleEntity.Tags).To(HaveLen(2))
			Expect(ruleEntity.Tags).To(ContainElements("transaction_type", "payment"))
		})

		It("should handle special characters in rule values", func() {
			form := inbound.RuleForm{
				Name:    "special_rule",
				Pattern: ".*special.*",
				Tags:    []string{"special-value_123!@#"},
			}

			ruleEntity := inbound.FormToEntity(form)
			Expect(ruleEntity.Tags).To(ContainElement("special-value_123!@#"))
		})

		It("should validate empty required fields", func() {
			form := inbound.RuleForm{
				Name:    "", // Empty name
				Pattern: ".*test.*",
				Tags:    []string{"category"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.RuleFormValidation, inbound.RuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})
	})
})

// TestRuleHandlers is handled by the main handlers_suite_test.go
