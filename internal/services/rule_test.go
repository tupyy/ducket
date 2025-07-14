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

			updated, err := ruleService.UpdateOrCreate(context.TODO(), newRule)
			Expect(err).To(BeNil())
			Expect(updated).To(BeFalse()) // Should be false for creation

			// Verify rule was created
			rule, err := ruleService.GetRule(context.TODO(), "upsert-new")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
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
			updated, err := ruleService.UpdateOrCreate(context.TODO(), updatedRule)
			Expect(err).To(BeNil())
			Expect(updated).To(BeTrue()) // Should be true for update

			// Verify rule was updated
			rule, err := ruleService.GetRule(context.TODO(), "upd-exist")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Pattern).To(Equal("updated-pattern"))
		})

		It("should handle updating rule with new labels", func() {
			// First, create a rule
			originalLabel := entity.Label{Key: "old-key", Value: "old-value"}
			originalRule := entity.NewRule("upd-with-labels", "pattern", originalLabel)
			err := ruleService.Create(context.TODO(), originalRule)
			Expect(err).To(BeNil())

			// Update with new labels
			newLabel1 := entity.Label{Key: "new-key1", Value: "new-value1"}
			newLabel2 := entity.Label{Key: "new-key2", Value: "new-value2"}
			updatedRule := entity.NewRule("upd-with-labels", "updated-pattern", newLabel1, newLabel2)
			updated, err := ruleService.UpdateOrCreate(context.TODO(), updatedRule)
			Expect(err).To(BeNil())
			Expect(updated).To(BeTrue())

			// Verify rule was updated with new labels
			rule, err := ruleService.GetRule(context.TODO(), "upd-with-labels")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Labels).To(HaveLen(2))
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
