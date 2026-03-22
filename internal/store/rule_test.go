package store_test

import (
	"context"
	"database/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tupyy/ducket/internal/entity"
	"github.com/tupyy/ducket/internal/store"
)

var _ = Describe("RuleStore", func() {
	var (
		ctx context.Context
		s   *store.Store
		db  *sql.DB
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = store.NewDB(":memory:")
		Expect(err).NotTo(HaveOccurred())

		s = store.NewStore(db)
		err = s.Migrate(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	createRule := func(name, filter string, tags []string) entity.Rule {
		r := entity.Rule{Name: name, Filter: filter, Tags: tags}
		id, err := s.CreateRule(ctx, r)
		Expect(err).NotTo(HaveOccurred())
		r.ID = id
		return r
	}

	Describe("CreateRule", func() {
		It("should create a rule and return its ID", func() {
			id, err := s.CreateRule(ctx, entity.Rule{
				Name:   "grocery",
				Filter: "content ~ /KAUFLAND|LIDL/",
				Tags:   []string{"food", "groceries"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(BeNumerically(">", 0))
		})

		It("should reject duplicate names", func() {
			createRule("grocery", "content ~ /KAUFLAND/", []string{"food"})

			_, err := s.CreateRule(ctx, entity.Rule{
				Name:   "grocery",
				Filter: "content ~ /LIDL/",
				Tags:   []string{"food"},
			})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetRule", func() {
		It("should return a rule by ID", func() {
			created := createRule("fuel", "content ~ /SHELL|OMV/", []string{"transport", "fuel"})

			rule, err := s.GetRule(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).NotTo(BeNil())
			Expect(rule.Name).To(Equal("fuel"))
			Expect(rule.Filter).To(Equal("content ~ /SHELL|OMV/"))
			Expect(rule.Tags).To(ConsistOf("transport", "fuel"))
		})

		It("should return nil for non-existent ID", func() {
			rule, err := s.GetRule(ctx, 999)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).To(BeNil())
		})
	})

	Describe("ListRules", func() {
		BeforeEach(func() {
			createRule("grocery", "content ~ /KAUFLAND/", []string{"food"})
			createRule("fuel", "content ~ /SHELL/", []string{"transport"})
			createRule("rent", "content ~ /CHIRIE/", []string{"housing"})
		})

		It("should list all rules without filter", func() {
			rules, err := s.ListRules(ctx, store.NoFilter, 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(3))
		})

		It("should filter by name", func() {
			rules, err := s.ListRules(ctx, "name = 'fuel'", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(1))
			Expect(rules[0].Name).To(Equal("fuel"))
		})

		It("should filter by name regex", func() {
			rules, err := s.ListRules(ctx, "name ~ /^(grocery|fuel)$/", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(2))
		})

		It("should paginate with limit", func() {
			rules, err := s.ListRules(ctx, store.NoFilter, 2, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(2))
		})

		It("should paginate with offset", func() {
			rules, err := s.ListRules(ctx, store.NoFilter, 1, 2)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(1))
			Expect(rules[0].Name).To(Equal("rent"))
		})

		It("should reject invalid filter", func() {
			_, err := s.ListRules(ctx, "invalid %%%", 0, 0)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("UpdateRule", func() {
		It("should update rule fields", func() {
			created := createRule("grocery", "content ~ /KAUFLAND/", []string{"food"})

			err := s.UpdateRule(ctx, entity.Rule{
				ID:     created.ID,
				Name:   "grocery-updated",
				Filter: "content ~ /KAUFLAND|LIDL|PENNY/",
				Tags:   []string{"food", "groceries", "essentials"},
			})
			Expect(err).NotTo(HaveOccurred())

			updated, err := s.GetRule(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Name).To(Equal("grocery-updated"))
			Expect(updated.Filter).To(Equal("content ~ /KAUFLAND|LIDL|PENNY/"))
			Expect(updated.Tags).To(ConsistOf("food", "groceries", "essentials"))
		})
	})

	Describe("DeleteRule", func() {
		It("should delete a rule", func() {
			created := createRule("fuel", "content ~ /SHELL/", []string{"transport"})

			err := s.DeleteRule(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())

			rule, err := s.GetRule(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).To(BeNil())
		})
	})

	Describe("Tags round-trip", func() {
		It("should preserve empty tags", func() {
			created := createRule("empty-tags", "content ~ /NOTHING/", []string{})

			rule, err := s.GetRule(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule.Tags).To(BeEmpty())
		})

		It("should preserve tags with special characters", func() {
			created := createRule("special", "content ~ /TEST/", []string{"it's", "a \"test\"", "semi;colon"})

			rule, err := s.GetRule(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule.Tags).To(ConsistOf("it's", "a \"test\"", "semi;colon"))
		})
	})
})
