package pg_test

import (
	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Filters", func() {
	var (
		psql      = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		baseQuery = psql.Select("*").From("labels")
	)

	Describe("LabelKeyValueQueryFilter", func() {
		It("should filter by key and value", func() {
			filter := pg.LabelKeyValueQueryFilter("category", "expense")
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE (key = $1 AND value = $2)"))
			Expect(args).To(Equal([]interface{}{"category", "expense"}))
		})

		It("should filter by key only when value is empty", func() {
			filter := pg.LabelKeyValueQueryFilter("category", "")
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE key = $1"))
			Expect(args).To(Equal([]interface{}{"category"}))
		})

		It("should return original query when key is empty", func() {
			filter := pg.LabelKeyValueQueryFilter("", "expense")
			result := filter(baseQuery)

			sql, _, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels"))
		})

		It("should return original query when both key and value are empty", func() {
			filter := pg.LabelKeyValueQueryFilter("", "")
			result := filter(baseQuery)

			sql, _, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels"))
		})
	})

	Describe("LabelKeyValuesQueryFilter", func() {
		It("should filter by key and single value", func() {
			filter := pg.LabelKeyValuesQueryFilter("category", []string{"expense"})
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE (key = $1 AND value IN ($2))"))
			Expect(args).To(Equal([]interface{}{"category", "expense"}))
		})

		It("should filter by key and multiple values using IN clause", func() {
			filter := pg.LabelKeyValuesQueryFilter("category", []string{"expense", "income", "transfer"})
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE (key = $1 AND value IN ($2,$3,$4))"))
			Expect(args).To(Equal([]interface{}{"category", "expense", "income", "transfer"}))
		})

		It("should filter by key only when values slice is empty", func() {
			filter := pg.LabelKeyValuesQueryFilter("category", []string{})
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE key = $1"))
			Expect(args).To(Equal([]interface{}{"category"}))
		})

		It("should filter by key only when values slice is nil", func() {
			filter := pg.LabelKeyValuesQueryFilter("category", nil)
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE key = $1"))
			Expect(args).To(Equal([]interface{}{"category"}))
		})

		It("should return original query when key is empty", func() {
			filter := pg.LabelKeyValuesQueryFilter("", []string{"expense", "income"})
			result := filter(baseQuery)

			sql, _, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels"))
		})

		It("should return original query when both key and values are empty", func() {
			filter := pg.LabelKeyValuesQueryFilter("", []string{})
			result := filter(baseQuery)

			sql, _, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels"))
		})
	})

	Describe("LabelQueryFilter", func() {
		It("should filter by single label value", func() {
			filter := pg.LabelQueryFilter([]string{"expense"})
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE (value = $1)"))
			Expect(args).To(Equal([]interface{}{"expense"}))
		})

		It("should filter by multiple label values using OR logic", func() {
			filter := pg.LabelQueryFilter([]string{"expense", "income", "transfer"})
			result := filter(baseQuery)

			sql, args, err := result.ToSql()
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal("SELECT * FROM labels WHERE (value = $1 OR value = $2 OR value = $3)"))
			Expect(args).To(Equal([]interface{}{"expense", "income", "transfer"}))
		})
	})
})
