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
	psql              = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertRuleStmt    = psql.Insert("rules").Columns("id", "pattern")
	insertTagStmt     = psql.Insert("tags").Columns("value")
	insertRuleTagStmt = psql.Insert("rules_tags").Columns("rule_id", "tag")
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
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions_tags;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules_tags;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM tags;")
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
			newRule := entity.NewRule("new-rule", "new-pattern", "tag1", "tag2")

			err := ruleService.Create(context.TODO(), newRule)
			Expect(err).To(BeNil())

			// Verify rule was created
			rule, err := ruleService.GetRule(context.TODO(), "new-rule")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Name).To(Equal("new-rule"))
			Expect(rule.Pattern).To(Equal("new-pattern"))
			Expect(rule.Tags).To(ContainElements("tag1", "tag2"))
		})

		It("should fail when rule already exists", func() {
			// First, create a rule
			existingRule := entity.NewRule("dup-rule", "pattern")
			err := ruleService.Create(context.TODO(), existingRule)
			Expect(err).To(BeNil())

			// Try to create the same rule again
			duplicateRule := entity.NewRule("dup-rule", "different-pattern")
			err = ruleService.Create(context.TODO(), duplicateRule)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("already exists"))
		})
	})

	Context("UpdateOrCreate", func() {
		It("should create new rule when it doesn't exist", func() {
			newRule := entity.NewRule("upsert-new", "pattern")

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
			originalRule := entity.NewRule("upd-exist", "original-pattern")
			err := ruleService.Create(context.TODO(), originalRule)
			Expect(err).To(BeNil())

			// Update the rule
			updatedRule := entity.NewRule("upd-exist", "updated-pattern", "new-tag")
			updated, err := ruleService.UpdateOrCreate(context.TODO(), updatedRule)
			Expect(err).To(BeNil())
			Expect(updated).To(BeTrue()) // Should be true for update

			// Verify rule was updated
			rule, err := ruleService.GetRule(context.TODO(), "upd-exist")
			Expect(err).To(BeNil())
			Expect(rule).ToNot(BeNil())
			Expect(rule.Pattern).To(Equal("updated-pattern"))
		})
	})

	Context("DeleteRule", func() {
		It("should delete existing rule successfully", func() {
			// First, create a rule
			ruleToDelete := entity.NewRule("del-rule", "pattern")
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
			filter := services.NewRuleFilterWithOptions(services.WithName("filter-1"))
			queryFilters := filter.QueriesFn()
			Expect(queryFilters).ToNot(BeEmpty())

			// The actual filtering is tested through the datastore layer
			// Here we just verify the filter generates query functions
		})
	})

	AfterEach(func() {
		// Clean up after each test
		_, err := pgPool.Exec(context.TODO(), "DELETE FROM transactions_tags;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules_tags;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM transactions;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM tags;")
		Expect(err).To(BeNil())
		_, err = pgPool.Exec(context.TODO(), "DELETE FROM rules;")
		Expect(err).To(BeNil())
	})
})
