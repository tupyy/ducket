// Package pg provides PostgreSQL datastore implementation with query filtering capabilities.
//
// This file contains query filter functions that can be applied to SQL SELECT statements
// to add WHERE clauses, LIMIT/OFFSET clauses, and other query modifications. These filters
// follow the functional composition pattern, allowing multiple filters to be chained together.
package pg

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

// TransactionHashQueryFilter creates a filter that adds a WHERE clause to match a specific transaction hash.
// If the provided hash is empty, the filter is a no-op and returns the original query unchanged.
//
// Parameters:
//   - hash: The transaction hash to filter by. If empty, no filtering is applied.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
func TransactionHashQueryFilter(hash string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(o sq.SelectBuilder) sq.SelectBuilder {
		if hash == "" {
			return o
		}
		return o.Where(sq.Eq{colHash: hash})
	}
}

// IntervalDateQueryFilter creates a filter that adds date range constraints to the query.
// This filter supports three modes of operation:
// 1. Both start and end provided: Filters for dates between start and end (inclusive)
// 2. Only start provided: Filters for dates greater than or equal to start
// 3. Only end provided: Filters for dates less than or equal to end
// 4. Neither provided: No filtering is applied
//
// Parameters:
//   - start: The earliest date to include (inclusive). Can be nil.
//   - end: The latest date to include (inclusive). Can be nil.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
func IntervalDateQueryFilter(start, end *time.Time) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if start != nil && end != nil {
			return orig.Where(
				sq.And{
					sq.GtOrEq{"date": start},
					sq.LtOrEq{"date": end},
				},
			)
		}

		if start != nil {
			return orig.Where(sq.GtOrEq{"date": start})
		}

		if end != nil {
			return orig.Where(sq.LtOrEq{"date": end})
		}

		return orig
	}
}

// LimitQueryFilter creates a filter that adds a LIMIT clause to restrict the number of results.
// If the limit is 0, no LIMIT clause is added to the query.
//
// Parameters:
//   - limit: The maximum number of rows to return. If 0, no limit is applied.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
func LimitQueryFilter(limit int) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if limit == 0 {
			return orig
		}
		return orig.Limit(uint64(limit))
	}
}

// OffsetQueryFilter creates a filter that adds an OFFSET clause to skip a number of rows.
// This is commonly used in conjunction with LimitQueryFilter for pagination.
//
// Parameters:
//   - offset: The number of rows to skip from the beginning of the result set.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
func OffsetQueryFilter(offset int) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		return orig.Offset(uint64(offset))
	}
}

// RuleIDQueryFilter creates a filter that adds a WHERE clause to match a specific rule by name.
// If the provided name is empty, the filter is a no-op and returns the original query unchanged.
//
// Parameters:
//   - name: The rule id (ID) to filter by. If empty, no filtering is applied.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
func RuleIDQueryFilter(id string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if id == "" {
			return orig
		}
		return orig.Where(sq.Eq{colFilterRuleID: id})
	}
}

// RuleLabelNameQueryFilter creates a filter that adds a WHERE clause to match records
// associated with a specific rule name. This is typically used when querying labels
// or transactions that are associated with a particular rule.
//
// Parameters:
//   - ruleName: The rule name to filter by. If empty, no filtering is applied.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
func RuleLabelNameQueryFilter(ruleName string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if ruleName == "" {
			return orig
		}
		return orig.Where(sq.Eq{colRuleID: ruleName})
	}
}

// LabelQueryFilter creates a filter that adds a WHERE clause to match any of the provided label values.
// This filter uses an OR condition, so records matching any of the provided label values will be included.
// If the labels slice is empty, the filter will still add an empty OR condition to the query.
//
// Parameters:
//   - labels: A slice of label values to match against. Each value is checked using equality.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
//
// Note: This function creates an OR condition even if the labels slice is empty, which may
// affect query performance. Consider checking for empty slices before applying this filter.
func LabelQueryFilter(labels []string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		or := sq.Or{}
		for _, l := range labels {
			or = append(or, sq.Eq{colLabelValue: l})
		}
		return orig.Where(or)
	}
}

// LabelKeyValueQueryFilter creates a filter that adds a WHERE clause to match labels by key and optionally by value.
// This filter requires a key and optionally filters by value as well.
// If value is empty, only the key constraint is applied.
//
// Parameters:
//   - key: The label key to filter by. This parameter is required and cannot be empty.
//   - value: The label value to filter by. If empty, no value filtering is applied.
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
func LabelKeyValueQueryFilter(key, value string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if key == "" {
			return orig
		}
		if value == "" {
			return orig.Where(sq.Eq{colLabelKey: key})
		}

		conditions := sq.And{sq.Eq{colLabelKey: key}}
		conditions = append(conditions, sq.Eq{colLabelValue: value})

		return orig.Where(conditions)
	}
}

// LabelKeyValuesQueryFilter creates a filter that adds a WHERE clause to match labels by key and any of the provided values.
// This filter requires a key and uses the SQL IN clause to match any of the provided values with OR logic.
// If the values slice is empty, only the key constraint is applied.
//
// Parameters:
//   - key: The label key to filter by. This parameter is required and cannot be empty.
//   - values: A slice of label values to match against. Uses SQL IN clause (equivalent to OR logic between values).
//
// Returns: A QueryFilter function that can be applied to a SelectBuilder.
//
// Example: LabelKeyValuesQueryFilter("category", []string{"expense", "income"})
// Generates: WHERE key = 'category' AND value IN ('expense', 'income')
// Which is equivalent to: WHERE key = 'category' AND (value = 'expense' OR value = 'income')
func LabelKeyValuesQueryFilter(key string, values []string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if key == "" {
			return orig
		}
		if len(values) == 0 {
			return orig.Where(sq.Eq{colLabelKey: key})
		}

		conditions := sq.And{sq.Eq{colLabelKey: key}}
		// sq.Eq with a slice generates an IN clause, providing OR logic between values
		conditions = append(conditions, sq.Eq{colLabelValue: values})

		return orig.Where(conditions)
	}
}
