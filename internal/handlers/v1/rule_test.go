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
		It("should handle missing form data", func() {
			req, _ := http.NewRequest("POST", "/api/v1/rules", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should handle empty JSON payload", func() {
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should validate rule name length", func() {
			form := inbound.RuleForm{
				Name:    "this_is_a_very_long_rule_name_that_exceeds_the_twenty_character_limit",
				Pattern: ".*test.*",
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate regex pattern", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: "[invalid_regex", // Invalid regex pattern
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate empty labels", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Labels:  map[string]string{},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["error"]).ToNot(BeNil())
		})
	})

	Context("PUT /api/v1/rules/:id", func() {
		It("should create rule even if the rule does not exists", func() {
			form := inbound.UpdateRuleForm{
				Pattern: ".*updated.*",
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", "/api/v1/rules/missing_id", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(BeElementOf([]int{http.StatusCreated, http.StatusOK}))
		})

		It("should validate form data for updates", func() {
			form := inbound.UpdateRuleForm{
				Pattern: ".*updated.*",
				Labels:  map[string]string{}, // Empty labels should fail validation
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", "/api/v1/rules/test_rule", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["error"]).ToNot(BeNil())
		})
	})

	Context("DELETE /api/v1/rules/:id", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("DELETE", "/api/v1/rules/nonexistent", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle gracefully even if rule doesn't exist
			Expect(w.Code).To(BeElementOf([]int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError, http.StatusNoContent}))
		})
	})

	Context("FormToEntity conversion", func() {
		It("should convert RuleForm to entity.Rule correctly", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Labels:  map[string]string{"category": "food"},
			}

			ruleEntity := inbound.FormToEntity(form)
			Expect(ruleEntity.Name).To(Equal("test_rule"))
			Expect(ruleEntity.Pattern).To(Equal(".*test.*"))
			Expect(ruleEntity.Labels).To(HaveLen(1))
			Expect(ruleEntity.Labels[0].Key).To(Equal("category"))
			Expect(ruleEntity.Labels[0].Value).To(Equal("food"))
		})

		It("should validate regex during form validation", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: "[invalid_regex", // Invalid regex
				Labels:  map[string]string{"category": "food"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.RuleFormValidation, inbound.RuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})

		It("should handle multiple labels correctly", func() {
			form := inbound.RuleForm{
				Name:    "multi_label_rule",
				Pattern: ".*multi.*",
				Labels:  map[string]string{"category": "food", "type": "essential"},
			}

			ruleEntity := inbound.FormToEntity(form)
			Expect(ruleEntity.Name).To(Equal("multi_label_rule"))
			Expect(ruleEntity.Pattern).To(Equal(".*multi.*"))
			Expect(ruleEntity.Labels).To(HaveLen(2))

			// Check that both labels are present
			labelKeys := make(map[string]string)
			for _, label := range ruleEntity.Labels {
				labelKeys[label.Key] = label.Value
			}
			Expect(labelKeys["category"]).To(Equal("food"))
			Expect(labelKeys["type"]).To(Equal("essential"))
		})

		It("should validate label length", func() {
			form := inbound.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Labels:  map[string]string{"category": "this_is_a_very_long_label_value_that_exceeds_the_twenty_character_limit"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.RuleFormValidation, inbound.RuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})
	})

	Context("UpdateRuleForm validation", func() {
		It("should validate pattern correctly", func() {
			form := inbound.UpdateRuleForm{
				Pattern: ".*valid.*",
				Labels:  map[string]string{"category": "food"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.UpdateRuleFormValidation, inbound.UpdateRuleForm{})

			err := validator.Struct(form)
			Expect(err).To(BeNil())
		})

		It("should reject invalid regex pattern", func() {
			form := inbound.UpdateRuleForm{
				Pattern: "[invalid",
				Labels:  map[string]string{"category": "food"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.UpdateRuleFormValidation, inbound.UpdateRuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})

		It("should reject empty labels", func() {
			form := inbound.UpdateRuleForm{
				Pattern: ".*valid.*",
				Labels:  map[string]string{},
			}

			validator := validator.New()
			validator.RegisterStructValidation(inbound.UpdateRuleFormValidation, inbound.UpdateRuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})
	})
})

// TestRuleHandlers is handled by the main handlers_suite_test.go
