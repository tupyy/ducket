package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tupyy/ducket/internal/entity"
	pkgErrors "github.com/tupyy/ducket/pkg/errors"
	sq "github.com/Masterminds/squirrel"
)

func (s *Store) ListRules(ctx context.Context, filter string, limit, offset int) ([]entity.Rule, error) {
	q := sq.Select("id", "name", "filter", "tags", "created_at").
		From("rules").
		OrderBy("id ASC").
		PlaceholderFormat(sq.Question)

	if filter != "" {
		sqlizer, err := ParseRuleFilter(filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter: %w", err)
		}
		q = q.Where(sqlizer)
	}
	if limit > 0 {
		q = q.Limit(uint64(limit))
	}
	if offset > 0 {
		q = q.Offset(uint64(offset))
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.qi.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return scanRules(rows)
}

func (s *Store) GetRule(ctx context.Context, id int) (*entity.Rule, error) {
	query, args, err := sq.Select("id", "name", "filter", "tags", "created_at").
		From("rules").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.qi.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	rules, err := scanRules(rows)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return nil, pkgErrors.NewResourceNotFoundError("rule", fmt.Sprintf("%d", id))
	}
	return &rules[0], nil
}

func (s *Store) CreateRule(ctx context.Context, r entity.Rule) (int, error) {
	tagsLiteral := duckdbArray(r.Tags)

	query, args, err := sq.Insert("rules").
		Columns("name", "filter", "tags").
		Values(r.Name, r.Filter, sq.Expr(tagsLiteral)).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int
	err = s.qi.QueryRowContext(ctx, query, args...).Scan(&id)
	return id, err
}

func (s *Store) UpdateRule(ctx context.Context, r entity.Rule) error {
	tagsLiteral := duckdbArray(r.Tags)

	query, args, err := sq.Update("rules").
		Set("name", r.Name).
		Set("filter", r.Filter).
		Set("tags", sq.Expr(tagsLiteral)).
		Where(sq.Eq{"id": r.ID}).
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return err
	}

	result, err := s.qi.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return pkgErrors.NewResourceNotFoundError("rule", fmt.Sprintf("%d", r.ID))
	}
	return nil
}

func (s *Store) DeleteRule(ctx context.Context, id int) error {
	query, args, err := sq.Delete("rules").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return err
	}

	result, err := s.qi.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return pkgErrors.NewResourceNotFoundError("rule", fmt.Sprintf("%d", id))
	}
	return nil
}

func scanRules(rows *sql.Rows) ([]entity.Rule, error) {
	var rules []entity.Rule
	for rows.Next() {
		var r entity.Rule
		var tags StringArray
		if err := rows.Scan(&r.ID, &r.Name, &r.Filter, &tags, &r.CreatedAt); err != nil {
			return nil, err
		}
		r.Tags = []string(tags)
		rules = append(rules, r)
	}
	return rules, rows.Err()
}

// duckdbArray converts a Go string slice to a DuckDB array literal: ['a', 'b']
func duckdbArray(tags []string) string {
	if len(tags) == 0 {
		return "[]::VARCHAR[]"
	}
	escaped := make([]string, len(tags))
	for i, t := range tags {
		escaped[i] = "'" + strings.ReplaceAll(t, "'", "''") + "'"
	}
	return "[" + strings.Join(escaped, ", ") + "]"
}
