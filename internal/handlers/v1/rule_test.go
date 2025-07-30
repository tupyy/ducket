package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	v1Impl "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuleHandlers", func() {
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

		// Register rule handlers
		v1.RegisterHandlers(router, v1Impl.NewServer())
		srv = httptest.NewServer(engine)
	})

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
		srv.Close()
	})

	Context("GET /api/v1/rules", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/rules", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// The handler should return OK or accept the request
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusAccepted}))
		})

		It("should return JSON response", func() {
			req, _ := http.NewRequest("GET", srv.URL+"/api/v1/rules", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			if resp.StatusCode == http.StatusOK {
				contentType := resp.Header.Get("Content-Type")
				Expect(contentType).To(ContainSubstring("application/json"))
			}
		})
	})

	Context("POST /api/v1/rules", func() {
		It("should handle missing form data", func() {
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/rules", nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should handle empty JSON payload", func() {
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/rules", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should validate rule name length", func() {
			form := v1.RuleForm{
				Name:    "this_is_a_very_long_rule_name_that_exceeds_the_twenty_character_limit",
				Pattern: ".*test.*",
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate regex pattern", func() {
			form := v1.RuleForm{
				Name:    "test_rule",
				Pattern: "[invalid_regex", // Invalid regex pattern
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})

		It("should validate empty labels", func() {
			form := v1.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Labels:  map[string]string{},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("POST", srv.URL+"/api/v1/rules", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})
	})

	Context("PUT /api/v1/rules/:id", func() {
		It("should create rule even if the rule does not exists", func() {
			form := v1.RuleForm{
				Name:    "test-rule",
				Pattern: ".*updated.*",
				Labels:  map[string]string{"category": "food"},
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", srv.URL+"/api/v1/rules/missing_id", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusCreated, http.StatusOK}))
		})

		It("should validate form data for updates", func() {
			form := v1.RuleForm{
				Pattern: ".*updated.*",
				Labels:  map[string]string{}, // Empty labels should fail validation
			}

			jsonData, _ := json.Marshal(form)
			req, _ := http.NewRequest("PUT", srv.URL+"/api/v1/rules/test_rule", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response["error"]).ToNot(BeNil())
		})
	})

	Context("DELETE /api/v1/rules/:id", func() {
		It("should handle requests without crashing", func() {
			req, _ := http.NewRequest("DELETE", srv.URL+"/api/v1/rules/nonexistent", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(BeNil())

			// Should handle gracefully even if rule doesn't exist
			Expect(resp.StatusCode).To(BeElementOf([]int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError, http.StatusNoContent}))
		})
	})

	Context("FormToEntity conversion", func() {
		It("should convert RuleForm to entity.Rule correctly", func() {
			form := v1.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Labels:  map[string]string{"category": "food"},
			}

			ruleEntity := form.Entity()
			Expect(ruleEntity.Name).To(Equal("test_rule"))
			Expect(ruleEntity.Pattern).To(Equal(".*test.*"))
			Expect(ruleEntity.Labels).To(HaveLen(1))
			Expect(ruleEntity.Labels[0].Key).To(Equal("category"))
			Expect(ruleEntity.Labels[0].Value).To(Equal("food"))
		})

		It("should validate regex during form validation", func() {
			form := v1.RuleForm{
				Name:    "test_rule",
				Pattern: "[invalid_regex", // Invalid regex
				Labels:  map[string]string{"category": "food"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(v1.RuleFormValidation, v1.RuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})

		It("should handle multiple labels correctly", func() {
			form := v1.RuleForm{
				Name:    "multi_label_rule",
				Pattern: ".*multi.*",
				Labels:  map[string]string{"category": "food", "type": "essential"},
			}

			ruleEntity := form.Entity()
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
			form := v1.RuleForm{
				Name:    "test_rule",
				Pattern: ".*test.*",
				Labels:  map[string]string{"category": "this_is_a_very_long_label_value_that_exceeds_the_twenty_character_limit"},
			}

			validator := validator.New()
			validator.RegisterStructValidation(v1.RuleFormValidation, v1.RuleForm{})

			err := validator.Struct(form)
			Expect(err).ToNot(BeNil())
		})
	})
})

// TestRuleHandlers is handled by the main handlers_suite_test.go
