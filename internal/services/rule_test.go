package services_test

import (
	"context"
	"database/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tupyy/ducket/internal/entity"
	"github.com/tupyy/ducket/internal/services"
	"github.com/tupyy/ducket/internal/store"
	srvErrors "github.com/tupyy/ducket/pkg/errors"
)

var _ = Describe("RuleService", func() {
	var (
		ctx context.Context
		svc *services.RuleService
		db  *sql.DB
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = store.NewDB(":memory:")
		Expect(err).NotTo(HaveOccurred())

		st := store.NewStore(db)
		err = st.Migrate(ctx)
		Expect(err).NotTo(HaveOccurred())

		svc = services.NewRuleService(st)
	})

	AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	Describe("Create", func() {
		It("should create a rule and return it with ID", func() {
			rule, err := svc.Create(ctx, entity.Rule{
				Name:   "grocery",
				Filter: "content ~ /KAUFLAND|LIDL/",
				Tags:   []string{"food", "groceries"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(rule.ID).To(BeNumerically(">", 0))
			Expect(rule.Name).To(Equal("grocery"))
		})

		It("should reject invalid filter", func() {
			_, err := svc.Create(ctx, entity.Rule{
				Name:   "bad",
				Filter: "invalid %%% syntax",
				Tags:   []string{"x"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid filter"))
		})
	})

	Describe("Get", func() {
		It("should return a rule by ID", func() {
			created, err := svc.Create(ctx, entity.Rule{
				Name:   "fuel",
				Filter: "content ~ /SHELL/",
				Tags:   []string{"transport"},
			})
			Expect(err).NotTo(HaveOccurred())

			rule, err := svc.Get(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).NotTo(BeNil())
			Expect(rule.Name).To(Equal("fuel"))
		})

		It("should return nil for non-existent ID", func() {
			rule, err := svc.Get(ctx, 999)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).To(BeNil())
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			_, err := svc.Create(ctx, entity.Rule{Name: "grocery", Filter: "content ~ /KAUFLAND/", Tags: []string{"food"}})
			Expect(err).NotTo(HaveOccurred())
			_, err = svc.Create(ctx, entity.Rule{Name: "fuel", Filter: "content ~ /SHELL/", Tags: []string{"transport"}})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should list all rules", func() {
			rules, err := svc.List(ctx, store.NoFilter, 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(2))
		})

		It("should filter rules", func() {
			rules, err := svc.List(ctx, "name = 'fuel'", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(1))
			Expect(rules[0].Name).To(Equal("fuel"))
		})

		It("should paginate", func() {
			rules, err := svc.List(ctx, store.NoFilter, 1, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(1))
		})
	})

	Describe("Update", func() {
		It("should update an existing rule", func() {
			created, err := svc.Create(ctx, entity.Rule{
				Name:   "grocery",
				Filter: "content ~ /KAUFLAND/",
				Tags:   []string{"food"},
			})
			Expect(err).NotTo(HaveOccurred())

			err = svc.Update(ctx, entity.Rule{
				ID:     created.ID,
				Name:   "grocery-v2",
				Filter: "content ~ /KAUFLAND|LIDL/",
				Tags:   []string{"food", "groceries"},
			})
			Expect(err).NotTo(HaveOccurred())

			updated, err := svc.Get(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Name).To(Equal("grocery-v2"))
			Expect(updated.Tags).To(ConsistOf("food", "groceries"))
		})

		It("should return not found for non-existent rule", func() {
			err := svc.Update(ctx, entity.Rule{ID: 999, Name: "x", Filter: "content ~ /X/", Tags: []string{"x"}})
			Expect(err).To(HaveOccurred())
			Expect(srvErrors.IsResourceNotFoundError(err)).To(BeTrue())
		})

		It("should reject invalid filter", func() {
			created, err := svc.Create(ctx, entity.Rule{
				Name:   "test",
				Filter: "content ~ /OK/",
				Tags:   []string{"ok"},
			})
			Expect(err).NotTo(HaveOccurred())

			err = svc.Update(ctx, entity.Rule{
				ID:     created.ID,
				Name:   "test",
				Filter: "bad %%% filter",
				Tags:   []string{"ok"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid filter"))
		})
	})

	Describe("Delete", func() {
		It("should delete a rule", func() {
			created, err := svc.Create(ctx, entity.Rule{
				Name:   "to-delete",
				Filter: "content ~ /DELETE/",
				Tags:   []string{"gone"},
			})
			Expect(err).NotTo(HaveOccurred())

			err = svc.Delete(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())

			rule, err := svc.Get(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).To(BeNil())
		})
	})
})
