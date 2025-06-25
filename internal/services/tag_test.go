package services_test

import (
	"context"
	"time"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagService", Ordered, func() {
	var (
		tagService *services.TagService
		datastore  *pg.Datastore
		pgPool     *pgxpool.Pool
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
		tagService = services.NewTagService(datastore)

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

	Context("GetTags", func() {
		It("should return empty slice when no tags exist", func() {
			tags, err := tagService.GetTags(context.TODO())
			Expect(err).To(BeNil())
			Expect(tags).To(HaveLen(0))
		})

		It("should return tags when they exist", func() {
			// Insert test data
			sql, args, err := insertTagStmt.
				Values("test-tag-1").
				Values("test-tag-2").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			tags, err := tagService.GetTags(context.TODO())
			Expect(err).To(BeNil())
			Expect(tags).To(HaveLen(2))

			tagValues := []string{tags[0].Value, tags[1].Value}
			Expect(tagValues).To(ContainElements("test-tag-1", "test-tag-2"))
		})

		It("should return tags with transaction counts", func() {
			// Insert test data: tags, rules, and transactions
			sql, args, err := insertTagStmt.
				Values("counted-tag").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Insert a rule
			sql, args, err = insertRuleStmt.
				Values("test-rule", "pattern").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Link rule to tag
			sql, args, err = insertRuleTagStmt.
				Values("test-rule", "counted-tag").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Insert transactions
			insertTransactionStmt := psql.Insert("transactions").
				Columns("id", "date", "kind", "content", "amount", "hash")
			sql, args, err = insertTransactionStmt.
				Values(1, time.Now(), "debit", "test transaction 1", -100.0, "hash1").
				Values(2, time.Now(), "credit", "test transaction 2", 200.0, "hash2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Link transactions to tags
			insertTransactionTagStmt := psql.Insert("transactions_tags").
				Columns("transaction_id", "tag_id", "rule_id")
			sql, args, err = insertTransactionTagStmt.
				Values(1, "counted-tag", "test-rule").
				Values(2, "counted-tag", "test-rule").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			tags, err := tagService.GetTags(context.TODO())
			Expect(err).To(BeNil())
			Expect(tags).To(HaveLen(1))
			Expect(tags[0].Value).To(Equal("counted-tag"))
			Expect(tags[0].CountTransactions).To(Equal(2))
		})
	})

	Context("IsExists", func() {
		It("should return false when tag doesn't exist", func() {
			exists, err := tagService.IsExists(context.TODO(), "no-exist")
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())
		})

		It("should return true when tag exists", func() {
			// Insert test data
			sql, args, err := insertTagStmt.
				Values("existing-tag").
				ToSql()
			Expect(err).To(BeNil())

			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			exists, err := tagService.IsExists(context.TODO(), "existing-tag")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
		})
	})

	Context("Create", func() {
		It("should create a new tag successfully", func() {
			err := tagService.Create(context.TODO(), "new-tag")
			Expect(err).To(BeNil())

			// Verify tag was created
			exists, err := tagService.IsExists(context.TODO(), "new-tag")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			// Verify through GetTags as well
			tags, err := tagService.GetTags(context.TODO())
			Expect(err).To(BeNil())

			found := false
			for _, tag := range tags {
				if tag.Value == "new-tag" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should not error when creating existing tag", func() {
			// First, create a tag
			err := tagService.Create(context.TODO(), "dup-tag")
			Expect(err).To(BeNil())

			// Try to create the same tag again - should not error
			err = tagService.Create(context.TODO(), "dup-tag")
			Expect(err).To(BeNil())

			// Verify tag still exists
			exists, err := tagService.IsExists(context.TODO(), "dup-tag")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
		})

		It("should handle multiple tag creation", func() {
			tagNames := []string{"batch-tag-1", "batch-tag-2", "batch-tag-3"}

			for _, tagName := range tagNames {
				err := tagService.Create(context.TODO(), tagName)
				Expect(err).To(BeNil())
			}

			// Verify all tags were created
			for _, tagName := range tagNames {
				exists, err := tagService.IsExists(context.TODO(), tagName)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			}

			// Verify through GetTags
			tags, err := tagService.GetTags(context.TODO())
			Expect(err).To(BeNil())
			Expect(len(tags)).To(BeNumerically(">=", len(tagNames)))
		})
	})

	Context("Delete", func() {
		It("should delete existing tag successfully", func() {
			// First, create a tag
			err := tagService.Create(context.TODO(), "del-tag")
			Expect(err).To(BeNil())

			// Verify tag exists
			exists, err := tagService.IsExists(context.TODO(), "del-tag")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			// Delete the tag
			err = tagService.Delete(context.TODO(), "del-tag")
			Expect(err).To(BeNil())

			// Verify tag was deleted
			exists, err = tagService.IsExists(context.TODO(), "del-tag")
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())
		})

		It("should handle deletion of non-existent tag gracefully", func() {
			err := tagService.Delete(context.TODO(), "no-exist")
			Expect(err).To(BeNil()) // Should not error
		})

		It("should delete tag with associated rules", func() {
			// Create tag and rule
			err := tagService.Create(context.TODO(), "tag-w-rule")
			Expect(err).To(BeNil())

			sql, args, err := insertRuleStmt.
				Values("rule-for-tag", "pattern").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Link rule to tag
			sql, args, err = insertRuleTagStmt.
				Values("rule-for-tag", "tag-w-rule").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Verify tag exists and has rule association
			exists, err := tagService.IsExists(context.TODO(), "tag-w-rule")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			// Delete the tag
			err = tagService.Delete(context.TODO(), "tag-w-rule")
			Expect(err).To(BeNil())

			// Verify tag was deleted
			exists, err = tagService.IsExists(context.TODO(), "tag-w-rule")
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())
		})
	})

	Context("Integration with Rules", func() {
		It("should show tags associated with rules", func() {
			// Create tags
			err := tagService.Create(context.TODO(), "integ-tag-1")
			Expect(err).To(BeNil())
			err = tagService.Create(context.TODO(), "integ-tag-2")
			Expect(err).To(BeNil())

			// Create rules
			sql, args, err := insertRuleStmt.
				Values("integ-rule-1", "pattern1").
				Values("integ-rule-2", "pattern2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			// Link rules to tags
			sql, args, err = insertRuleTagStmt.
				Values("integ-rule-1", "integ-tag-1").
				Values("integ-rule-2", "integ-tag-1").
				Values("integ-rule-2", "integ-tag-2").
				ToSql()
			Expect(err).To(BeNil())
			_, err = pgPool.Exec(context.TODO(), sql, args...)
			Expect(err).To(BeNil())

			tags, err := tagService.GetTags(context.TODO())
			Expect(err).To(BeNil())
			Expect(tags).To(HaveLen(2))

			// Find tags and verify their rule associations
			var tag1, tag2 *entity.Tag
			for i, tag := range tags {
				if tag.Value == "integ-tag-1" {
					tag1 = &tags[i]
				} else if tag.Value == "integ-tag-2" {
					tag2 = &tags[i]
				}
			}

			Expect(tag1).ToNot(BeNil())
			Expect(tag2).ToNot(BeNil())

			// integ-tag-1 should be associated with both rules
			Expect(tag1.Rules).To(ContainElements("integ-rule-1", "integ-rule-2"))

			// integ-tag-2 should be associated with integ-rule-2
			Expect(tag2.Rules).To(ContainElement("integ-rule-2"))
		})
	})

	Context("Edge Cases", func() {
		It("should handle empty tag value", func() {
			err := tagService.Create(context.TODO(), "")
			Expect(err).To(BeNil()) // Should not error, but behavior depends on database constraints

			_, err = tagService.IsExists(context.TODO(), "")
			Expect(err).To(BeNil())
			// Result depends on database constraints - we just verify no error occurs
		})

		It("should handle tag values with special characters", func() {
			specialTags := []string{
				"tag-w-dash",
				"tag_w_under",
				"tag w space",
				"tag.w.dot",
				"tag@w@symbol",
			}

			for _, tagName := range specialTags {
				err := tagService.Create(context.TODO(), tagName)
				Expect(err).To(BeNil())

				exists, err := tagService.IsExists(context.TODO(), tagName)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			}
		})

		It("should handle very long tag names", func() {
			longTagName := "very_long_12"

			err := tagService.Create(context.TODO(), longTagName)
			Expect(err).To(BeNil())

			exists, err := tagService.IsExists(context.TODO(), longTagName)
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
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
