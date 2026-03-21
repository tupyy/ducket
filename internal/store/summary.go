package store

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type SummaryOverview struct {
	TotalTransactions int
	TotalDebit        float64
	TotalCredit       float64
	Balance           float64
	UniqueAccounts    int
	UniqueTags        int
}

type TagSummary struct {
	Tag        string
	TotalDebit float64
	TotalCredit float64
	Count      int
}

type BalanceTrendPoint struct {
	Month   string
	Debit   float64
	Credit  float64
	Balance float64
}

func (s *Store) GetSummaryOverview(ctx context.Context, filter string) (*SummaryOverview, error) {
	q := sq.Select(
		"COUNT(*)",
		"COALESCE(SUM(CASE WHEN kind='debit' THEN amount ELSE 0 END), 0)",
		"COALESCE(SUM(CASE WHEN kind='credit' THEN amount ELSE 0 END), 0)",
		"COUNT(DISTINCT account)",
	).From("transactions t").PlaceholderFormat(sq.Question)

	if filter != "" {
		sqlizer, err := ParseTransactionFilter(filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter: %w", err)
		}
		q = q.Where(sqlizer)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var overview SummaryOverview
	if err := s.qi.QueryRowContext(ctx, query, args...).Scan(
		&overview.TotalTransactions,
		&overview.TotalDebit,
		&overview.TotalCredit,
		&overview.UniqueAccounts,
	); err != nil {
		return nil, err
	}
	overview.Balance = overview.TotalCredit - overview.TotalDebit

	uniqueTags, err := s.countUniqueTags(ctx, filter)
	if err != nil {
		return nil, err
	}
	overview.UniqueTags = uniqueTags

	return &overview, nil
}

func (s *Store) countUniqueTags(ctx context.Context, filter string) (int, error) {
	rules, err := s.ListRules(ctx, NoFilter, 0, 0)
	if err != nil {
		return 0, err
	}
	tagCTE, tagArgs := buildTagCTE(rules)
	if tagCTE == "" {
		return 0, nil
	}

	countQ := sq.Select("COUNT(DISTINCT tag)").
		From("rule_tags, UNNEST(tags) AS t(tag)").
		PlaceholderFormat(sq.Question)

	if filter != "" {
		countQ = countQ.Where("transaction_id IN (SELECT t.id FROM transactions t WHERE " + filter + ")")
	}

	countSQL, countSelectArgs, err := countQ.ToSql()
	if err != nil {
		return 0, err
	}

	fullSQL := "WITH " + tagCTE + " " + countSQL
	allArgs := make([]any, 0, len(tagArgs)+len(countSelectArgs))
	allArgs = append(allArgs, tagArgs...)
	allArgs = append(allArgs, countSelectArgs...)

	var count int
	if err := s.qi.QueryRowContext(ctx, fullSQL, allArgs...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) GetTagSummary(ctx context.Context, filter string) ([]TagSummary, error) {
	rules, err := s.ListRules(ctx, NoFilter, 0, 0)
	if err != nil {
		return nil, err
	}
	tagCTE, tagArgs := buildTagCTE(rules)
	if tagCTE == "" {
		return []TagSummary{}, nil
	}

	var filterClause string
	if filter != "" {
		sqlizer, err := ParseTransactionFilter(filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter: %w", err)
		}
		filterSQL, filterArgs, err := sqlizer.ToSql()
		if err != nil {
			return nil, err
		}
		filterClause = " AND " + filterSQL
		tagArgs = append(tagArgs, filterArgs...)
	}

	query := "WITH " + tagCTE + `, unnested AS (SELECT transaction_id, UNNEST(tags) AS tag FROM rule_tags)
SELECT u.tag,
       COALESCE(SUM(CASE WHEN t.kind='debit' THEN t.amount ELSE 0 END), 0) AS total_debit,
       COALESCE(SUM(CASE WHEN t.kind='credit' THEN t.amount ELSE 0 END), 0) AS total_credit,
       COUNT(*) AS count
FROM unnested u
JOIN transactions t ON t.id = u.transaction_id
WHERE 1=1` + filterClause + `
GROUP BY u.tag
ORDER BY total_debit DESC`

	rows, err := s.qi.QueryContext(ctx, query, tagArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTagSummaries(rows)
}

func scanTagSummaries(rows *sql.Rows) ([]TagSummary, error) {
	var results []TagSummary
	for rows.Next() {
		var ts TagSummary
		if err := rows.Scan(&ts.Tag, &ts.TotalDebit, &ts.TotalCredit, &ts.Count); err != nil {
			return nil, err
		}
		results = append(results, ts)
	}
	if results == nil {
		results = []TagSummary{}
	}
	return results, rows.Err()
}

func (s *Store) GetBalanceTrend(ctx context.Context, filter string) ([]BalanceTrendPoint, error) {
	q := sq.Select(
		"strftime(date, '%Y-%m') AS month",
		"COALESCE(SUM(CASE WHEN kind='debit' THEN amount ELSE 0 END), 0) AS debit",
		"COALESCE(SUM(CASE WHEN kind='credit' THEN amount ELSE 0 END), 0) AS credit",
	).From("transactions t").
		GroupBy("strftime(date, '%Y-%m')").
		OrderBy("month ASC").
		PlaceholderFormat(sq.Question)

	if filter != "" {
		sqlizer, err := ParseTransactionFilter(filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter: %w", err)
		}
		q = q.Where(sqlizer)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.qi.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []BalanceTrendPoint
	var cumulative float64
	for rows.Next() {
		var p BalanceTrendPoint
		if err := rows.Scan(&p.Month, &p.Debit, &p.Credit); err != nil {
			return nil, err
		}
		cumulative += p.Credit - p.Debit
		p.Balance = cumulative
		points = append(points, p)
	}
	if points == nil {
		points = []BalanceTrendPoint{}
	}
	return points, rows.Err()
}

