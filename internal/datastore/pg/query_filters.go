package pg

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

func TransactionHashQueryFilter(hash string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(o sq.SelectBuilder) sq.SelectBuilder {
		if hash == "" {
			return o
		}
		return o.Where(sq.Eq{colHash: hash})
	}
}

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

func LimitQueryFilter(limit int) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if limit == 0 {
			return orig
		}
		return orig.Limit(uint64(limit))
	}
}

func OffsetQueryFilter(offset int) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		return orig.Offset(uint64(offset))
	}
}

func RuleNameQueryFilter(name string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if name == "" {
			return orig
		}
		return orig.Where(sq.Eq{colID: name})
	}
}

func RuleTagNameQueryFilter(ruleName string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if ruleName == "" {
			return orig
		}
		return orig.Where(sq.Eq{colRuleID: ruleName})
	}
}

func TagQueryFilter(tags []string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		or := sq.Or{}
		for _, t := range tags {
			or = append(or, sq.Eq{colValue: t})
		}
		return orig.Where(or)
	}
}
