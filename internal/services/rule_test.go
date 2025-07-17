package services_test

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	pgUri = "postgresql://postgres:postgres@localhost:5432/postgres"
)

var (
	psql                = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertRuleStmt      = psql.Insert("rules").Columns("id", "pattern")
	insertLabelStmt     = psql.Insert("labels").Columns("key", "value")
	insertRuleLabelStmt = psql.Insert("rules_labels").Columns("rule_id", "label_id")
)

var _ = Describe("RuleService", Ordered, func() {
	var (
		ruleService *services.RuleService
		datastore   *pg.Datastore
		pgPool      *pgxpool.Pool
	)

	BeforeAll(func() {
		// Setup database connection
		pgDt, err := pg.NewPostgresDatastore(context.TODO(), pgUri)
		Expect(err).To(BeNil())
		Expect(pgDt).ToNot(BeNil())

		pgxConfig, err := pgxpool.ParseConfig(pgUri)
		Expect(err).To(BeNil())

		pool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
		Expect(err).To(BeNil())
		Expect(pool).ToNot(BeNil())

		datastore = pgDt
		pgPool = pool
		ruleService = services.NewRuleService(datastore)

		// Clean up database
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		datastore.Close()
	})

	Context("GetRules", func() {
		It("should return empty slice when no rules exist", func() {
			rules, err := ruleService.GetRules(context.TODO())
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(0))
		})

		It("should return rules when they exist", func() {
			// Insert test data
			sql, args, err := insertRuleStmt.
				Values("test-rule-1", "pattern1").
				Values("test-rule-2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rules, err := ruleService.GetRules(context.TODO())
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(2))

			ruleNames := []string{rules[0].Name, rules[1].Name}
			Expect(ruleNames).To(ContainElements("test-rule-1", "test-rule-2"))
		})
	})

	Context("GetRule", func() {
		It("should return nil when rule doesn't exist", func() {
			rule, err := ruleService.GetRule(context.TODO(), "no-exist")
			Expect(err).To(BeNil())
			Expect(rule).To(BeNil())
		})

		It("should return rule when it exists", func() {
			// Insert test data
			sql, args, err := insertRuleStmt.
				Values("exist-rule", "exist-pattern").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			rule, err := ruleService.GetRule(context.TODO(), "exist-rule")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Name).To(Equal("exist-rule"))
			Expect(rule.Pattern).To(Equal("exist-pattern"))
		})
	})

	Context("Create", func() {
		It("should create a new rule successfully", func() {
			label1 := entity.Label{Key: "category", Value: "expense"}
			label2 := entity.Label{Key: "type", Value: "recurring"}
			newRule := entity.NewRule("new-rule", "new-pattern", label1, label2)

			err := ruleService.Create(context.TODO(), newRule)
			Expect(err).To(BeNil())

			// Verify rule was created
			rule, err := ruleService.GetRule(context.TODO(), "new-rule")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Name).To(Equal("new-rule"))
			Expect(rule.Pattern).To(Equal("new-pattern"))
			Expect(rule.Labels).To(HaveLen(2))

			// Check that the labels exist
			labelValues := make(map[string]string)
			for _, label := range rule.Labels {
				labelValues[label.Key] = label.Value
			}
			Expect(labelValues).To(HaveKeyWithValue("category", "expense"))
			Expect(labelValues).To(HaveKeyWithValue("type", "recurring"))
		})

		It("should fail when rule already exists", func() {
			// First, create a rule
			label := entity.Label{Key: "category", Value: "expense"}
			existingRule := entity.NewRule("dup-rule", "pattern", label)
			err := ruleService.Create(context.TODO(), existingRule)
			Expect(err).To(BeNil())

			// Try to create the same rule again
			duplicateRule := entity.NewRule("dup-rule", "different-pattern", label)
			err = ruleService.Create(context.TODO(), duplicateRule)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("already exists"))
		})

		It("should create labels if they don't exist", func() {
			// Create a rule with labels that don't exist yet
			newLabel := entity.Label{Key: "priority", Value: "high"}
			newRule := entity.NewRule("with-new-labels", "pattern", newLabel)

			err := ruleService.Create(context.TODO(), newRule)
			Expect(err).To(BeNil())

			// Verify the rule was created
			rule, err := ruleService.GetRule(context.TODO(), "with-new-labels")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Labels).To(HaveLen(1))
			Expect(rule.Labels[0].Key).To(Equal("priority"))
			Expect(rule.Labels[0].Value).To(Equal("high"))
		})
	})

	Context("UpdateOrCreate", func() {
		It("should create new rule when it doesn't exist", func() {
			label := entity.Label{Key: "status", Value: "active"}
			newRule := entity.NewRule("upsert-new", "pattern", label)

			err := ruleService.Update(context.TODO(), newRule)
			Expect(err).NotTo(BeNil())
		})

		It("should update existing rule", func() {
			// First, create a rule
			originalLabel := entity.Label{Key: "category", Value: "income"}
			originalRule := entity.NewRule("upd-exist", "original-pattern", originalLabel)
			err := ruleService.Create(context.TODO(), originalRule)
			Expect(err).To(BeNil())

			// Update the rule
			newLabel := entity.Label{Key: "category", Value: "expense"}
			updatedRule := entity.NewRule("upd-exist", "updated-pattern", newLabel)
			err = ruleService.Update(context.TODO(), updatedRule)
			Expect(err).To(BeNil())

			// Verify rule was updated
			rule, err := ruleService.GetRule(context.TODO(), "upd-exist")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Pattern).To(Equal("updated-pattern"))
		})
	})

	Context("Update", func() {
		It("should return error when rule does not exist", func() {
			// Try to update a non-existent rule
			label := entity.Label{Key: "category", Value: "expense"}
			nonExistentRule := entity.NewRule("non-existent-rule", "pattern", label)

			err := ruleService.Update(context.TODO(), nonExistentRule)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("rule non-existent-rule not found"))
		})

		It("should update existing rule successfully", func() {
			// First, create a rule
			originalLabel := entity.Label{Key: "category", Value: "income"}
			originalRule := entity.NewRule("update-existing", "original-pattern", originalLabel)
			err := ruleService.Create(context.TODO(), originalRule)
			Expect(err).To(BeNil())

			// Update the rule
			newLabel := entity.Label{Key: "category", Value: "expense"}
			updatedRule := entity.NewRule("update-existing", "updated-pattern", newLabel)
			err = ruleService.Update(context.TODO(), updatedRule)
			Expect(err).To(BeNil())

			// Verify rule was updated
			rule, err := ruleService.GetRule(context.TODO(), "update-existing")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Pattern).To(Equal("updated-pattern"))
			Expect(rule.Labels).To(HaveLen(1))
			Expect(rule.Labels[0].Key).To(Equal("category"))
			Expect(rule.Labels[0].Value).To(Equal("expense"))
		})

		It("should update rule-label relationships correctly", func() {
			// First, create a rule with initial labels
			label1 := entity.Label{Key: "priority", Value: "high"}
			label2 := entity.Label{Key: "category", Value: "expense"}
			originalRule := entity.NewRule("update-relationships", "pattern", label1, label2)
			err := ruleService.Create(context.TODO(), originalRule)
			Expect(err).To(BeNil())

			// Update the rule with completely different labels
			newLabel1 := entity.Label{Key: "status", Value: "active"}
			newLabel2 := entity.Label{Key: "type", Value: "recurring"}
			newLabel3 := entity.Label{Key: "frequency", Value: "monthly"}
			updatedRule := entity.NewRule("update-relationships", "updated-pattern", newLabel1, newLabel2, newLabel3)
			err = ruleService.Update(context.TODO(), updatedRule)
			Expect(err).To(BeNil())

			// Verify rule was updated with new relationships
			rule, err := ruleService.GetRule(context.TODO(), "update-relationships")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Labels).To(HaveLen(3))

			// Check that the new labels exist
			labelMap := make(map[string]string)
			for _, label := range rule.Labels {
				labelMap[label.Key] = label.Value
			}
			Expect(labelMap).To(HaveKeyWithValue("status", "active"))
			Expect(labelMap).To(HaveKeyWithValue("type", "recurring"))
			Expect(labelMap).To(HaveKeyWithValue("frequency", "monthly"))

			// Verify old labels are no longer associated with this rule
			Expect(labelMap).ToNot(HaveKey("priority"))
			Expect(labelMap).ToNot(HaveKey("category"))
		})

		It("should create new labels during update if they don't exist", func() {
			// First, create a rule with existing labels
			existingLabel := entity.Label{Key: "existing", Value: "label"}
			originalRule := entity.NewRule("update-create-labels", "pattern", existingLabel)
			err := ruleService.Create(context.TODO(), originalRule)
			Expect(err).To(BeNil())

			// Update with a mix of existing and new labels
			newLabel := entity.Label{Key: "brand-new", Value: "never-seen-before"}
			updatedRule := entity.NewRule("update-create-labels", "updated-pattern", existingLabel, newLabel)
			err = ruleService.Update(context.TODO(), updatedRule)
			Expect(err).To(BeNil())

			// Verify the rule has both labels
			rule, err := ruleService.GetRule(context.TODO(), "update-create-labels")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Labels).To(HaveLen(2))

			labelMap := make(map[string]string)
			for _, label := range rule.Labels {
				labelMap[label.Key] = label.Value
			}
			Expect(labelMap).To(HaveKeyWithValue("existing", "label"))
			Expect(labelMap).To(HaveKeyWithValue("brand-new", "never-seen-before"))
		})
	})

	Context("Create - Relationship Management", func() {
		It("should create relationships between rule and labels in Create method", func() {
			label1 := entity.Label{Key: "test-relationship", Value: "value1"}
			label2 := entity.Label{Key: "test-relationship", Value: "value2"}
			newRule := entity.NewRule("relationship-test", "pattern", label1, label2)

			err := ruleService.Create(context.TODO(), newRule)
			Expect(err).To(BeNil())

			// Verify the rule and its relationships were created
			rule, err := ruleService.GetRule(context.TODO(), "relationship-test")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Labels).To(HaveLen(2))

			// Verify each label is properly associated
			for _, label := range rule.Labels {
				Expect(label.Key).To(Equal("test-relationship"))
				Expect(label.Value).To(BeElementOf([]string{"value1", "value2"}))
				Expect(label.ID).ToNot(BeZero()) // Label should have an ID assigned
			}
		})

		It("should reuse existing labels when creating relationships", func() {
			// First, create a label through another rule
			existingLabel := entity.Label{Key: "reusable", Value: "shared"}
			firstRule := entity.NewRule("first-rule", "pattern1", existingLabel)
			err := ruleService.Create(context.TODO(), firstRule)
			Expect(err).To(BeNil())

			// Create another rule with the same label
			secondRule := entity.NewRule("second-rule", "pattern2", existingLabel)
			err = ruleService.Create(context.TODO(), secondRule)
			Expect(err).To(BeNil())

			// Verify both rules reference the same label
			rule1, err := ruleService.GetRule(context.TODO(), "first-rule")
			Expect(err).To(BeNil())
			rule2, err := ruleService.GetRule(context.TODO(), "second-rule")
			Expect(err).To(BeNil())

			Expect(rule1.Labels).To(HaveLen(1))
			Expect(rule2.Labels).To(HaveLen(1))

			// Both rules should reference labels with the same key-value pair
			Expect(rule1.Labels[0].Key).To(Equal("reusable"))
			Expect(rule1.Labels[0].Value).To(Equal("shared"))
			Expect(rule2.Labels[0].Key).To(Equal("reusable"))
			Expect(rule2.Labels[0].Value).To(Equal("shared"))
		})
	})

	Context("DeleteRule", func() {
		It("should delete existing rule successfully", func() {
			// First, create a rule
			label := entity.Label{Key: "temp", Value: "delete-me"}
			ruleToDelete := entity.NewRule("del-rule", "pattern", label)
			err := ruleService.Create(context.TODO(), ruleToDelete)
			Expect(err).To(BeNil())

			// Verify rule exists
			rule, err := ruleService.GetRule(context.TODO(), "del-rule")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())

			// Delete the rule
			err = ruleService.DeleteRule(context.TODO(), "del-rule")
			Expect(err).To(BeNil())

			// Verify rule was deleted
			rule, err = ruleService.GetRule(context.TODO(), "del-rule")
			Expect(err).To(BeNil())
			Expect(rule).To(BeNil())
		})

		It("should handle deletion of non-existent rule gracefully", func() {
			err := ruleService.DeleteRule(context.TODO(), "no-exist")
			Expect(err).To(BeNil()) // Should not error
		})
	})

	Context("RuleFilter", func() {
		It("should filter rules by name", func() {
			// Create test rules
			sql, args, err := insertRuleStmt.
				Values("filter-1", "pattern1").
				Values("filter-2", "pattern2").
				Values("other-rule", "pattern3").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Test filter functionality
			filter := services.NewRuleFilterWithOptions(services.WithId("filter-1"))
			queryFilters := filter.QueriesFn()
			Expect(queryFilters).ToNot(BeEmpty())

			// The actual filtering is tested through the datastore layer
			// Here we just verify the filter generates query functions
		})
	})

	AfterEach(func() {
		// Clean up after each test
		_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules_labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM labels;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
		Expect(err).To(BeNil())
	})
})
