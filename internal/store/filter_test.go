package store_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tupyy/ducket/internal/store"
)

var _ = Describe("Filter", func() {
	Describe("ParseTransactionFilter", func() {
		It("should parse a simple equality", func() {
			sqlizer, err := store.ParseTransactionFilter("kind = 'debit'")
			Expect(err).NotTo(HaveOccurred())

			sql, args, err := sqlizer.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(ContainSubstring("t.kind"))
			Expect(args).To(ContainElement("debit"))
		})

		It("should parse a regex filter", func() {
			sqlizer, err := store.ParseTransactionFilter("content ~ /KAUFLAND/")
			Expect(err).NotTo(HaveOccurred())

			sql, _, err := sqlizer.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(ContainSubstring("regexp_matches"))
			Expect(sql).To(ContainSubstring("t.content"))
		})

		It("should parse combined filters", func() {
			sqlizer, err := store.ParseTransactionFilter("kind = 'debit' and amount > 100")
			Expect(err).NotTo(HaveOccurred())

			sql, _, err := sqlizer.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(ContainSubstring("t.kind"))
			Expect(sql).To(ContainSubstring("t.amount"))
		})

		It("should reject unknown fields", func() {
			_, err := store.ParseTransactionFilter("unknown_field = 'x'")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown transaction field"))
		})

		It("should reject invalid syntax", func() {
			_, err := store.ParseTransactionFilter("kind ===")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ParseRuleFilter", func() {
		It("should parse a name filter", func() {
			sqlizer, err := store.ParseRuleFilter("name = 'grocery'")
			Expect(err).NotTo(HaveOccurred())

			sql, args, err := sqlizer.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(ContainSubstring("name"))
			Expect(args).To(ContainElement("grocery"))
		})

		It("should reject unknown fields", func() {
			_, err := store.ParseRuleFilter("tags = 'food'")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown rule field"))
		})
	})
})
