package pg

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

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

func TagNameQueryFilter(ruleName string) func(orig sq.SelectBuilder) sq.SelectBuilder {
	return func(orig sq.SelectBuilder) sq.SelectBuilder {
		if ruleName == "" {
			return orig
		}
		return orig.Where(sq.Eq{colRuleID: ruleName})
	}
}
