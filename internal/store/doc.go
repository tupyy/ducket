// Package store implements the DuckDB persistence layer for the finance application.
//
// # Architecture
//
// The store uses a single-writer DuckDB connection (SetMaxOpenConns(1)) with a
// query interceptor that routes reads/writes through context-based transactions
// and auto-checkpoints after every write outside a transaction.
//
// There are only two tables: rules and transactions. Tags are not stored on
// transactions — they are computed dynamically at query time by evaluating each
// rule's filter expression against the transactions table.
//
// # Schema
//
//	CREATE TABLE rules (
//	    id          INTEGER PRIMARY KEY DEFAULT nextval('rules_id_seq'),
//	    name        VARCHAR NOT NULL UNIQUE,
//	    filter      VARCHAR NOT NULL,
//	    tags        VARCHAR[] NOT NULL,
//	    created_at  TIMESTAMP DEFAULT current_timestamp
//	);
//
//	CREATE TABLE transactions (
//	    id          INTEGER PRIMARY KEY DEFAULT nextval('transactions_id_seq'),
//	    hash        VARCHAR NOT NULL UNIQUE,
//	    date        DATE NOT NULL,
//	    account     BIGINT NOT NULL,
//	    kind        VARCHAR NOT NULL,
//	    amount      DECIMAL(12,2) NOT NULL,
//	    content     VARCHAR NOT NULL,
//	    info        VARCHAR,
//	    recipient   VARCHAR,
//	    created_at  TIMESTAMP DEFAULT current_timestamp
//	);
//
// # Filter DSL
//
// Both ListRules and ListTransactions accept a filter string parsed by the
// filter package (github.com/kubev2v/assisted-migration-agent/pkg/filter).
// The DSL supports comparison, regex, IN/NOT IN, and boolean operators.
//
// Transaction filter fields: id, hash, date, account, kind, amount, content, info, recipient
// Rule filter fields: id, name, filter
//
// # DuckDB Array Handling
//
// Tags are stored as VARCHAR[] in DuckDB. The duckdb-go driver returns these as
// []interface{} on scan. The StringArray type implements sql.Scanner to handle
// this conversion. For writes, tags are serialized as DuckDB array literals
// (e.g. ['food', 'groceries']) embedded directly in the SQL via sq.Expr.
//
// # Transaction Support
//
// The Store.WithTx method wraps a function in a DuckDB transaction. The *sql.Tx
// is stored in the context and automatically picked up by the query interceptor,
// so all store methods called within the closure participate in the same transaction.
//
//	store.WithTx(ctx, func(ctx context.Context) error {
//	    id, err := store.CreateTransaction(ctx, txn)
//	    if err != nil {
//	        return err // triggers rollback
//	    }
//	    return nil // triggers commit
//	})
//
// # Dynamic Tag CTE
//
// When querying transactions, the store loads all rules and builds a CTE that
// evaluates each rule's filter against the transactions table. Each rule's filter
// is parsed into a parameterized WHERE clause via the filter DSL. The results are
// aggregated so each transaction gets a merged, deduplicated tag array.
//
// Given two rules:
//
//	Rule 1: filter="content ~ /KAUFLAND|LIDL/"  tags=['food', 'groceries']
//	Rule 2: filter="content ~ /SHELL|OMV/"      tags=['transport', 'fuel']
//
// The generated CTE is:
//
//	WITH rule_matches AS (
//	    SELECT t.id AS transaction_id, ['food', 'groceries'] AS tags
//	    FROM transactions t
//	    WHERE regexp_matches(t.content, ?)
//	    UNION ALL
//	    SELECT t.id AS transaction_id, ['transport', 'fuel'] AS tags
//	    FROM transactions t
//	    WHERE regexp_matches(t.content, ?)
//	),
//	rule_tags AS (
//	    SELECT transaction_id,
//	           list_distinct(flatten(list(tags))) AS tags
//	    FROM rule_matches
//	    GROUP BY transaction_id
//	)
//
// When there are no rules, the CTE is omitted and tags default to an empty array.
// When filterTags is provided, only rules whose tags contain all requested tags
// are included and an INNER JOIN replaces the LEFT JOIN, so only matching
// transactions are returned.
//
// # Rule Queries
//
// ListRules returns all rules matching an optional filter, with pagination:
//
//	SELECT id, name, filter, tags, created_at
//	FROM rules
//	WHERE <filter>
//	ORDER BY id ASC
//	LIMIT ? OFFSET ?
//
// GetRule returns a single rule by ID:
//
//	SELECT id, name, filter, tags, created_at
//	FROM rules
//	WHERE id = ?
//
// CreateRule inserts a rule and returns its ID:
//
//	INSERT INTO rules (name, filter, tags)
//	VALUES (?, ?, ['food', 'groceries'])
//	RETURNING id
//
// UpdateRule modifies a rule by ID:
//
//	UPDATE rules SET name = ?, filter = ?, tags = ['food', 'groceries']
//	WHERE id = ?
//
// DeleteRule removes a rule by ID:
//
//	DELETE FROM rules WHERE id = ?
//
// # Transaction Queries
//
// ListTransactions returns transactions with dynamically computed tags,
// server-side pagination, sorting, and a total count. The tag CTE is prepended
// when rules exist:
//
//	WITH <tag CTE>
//	SELECT t.id, t.hash, t.date, t.account, t.kind,
//	       t.amount, t.content, t.info, t.recipient,
//	       COALESCE(rt.tags, []::VARCHAR[]) AS tags
//	FROM transactions t
//	LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
//	WHERE <filter>
//	ORDER BY t.date DESC, t.id DESC
//	LIMIT ? OFFSET ?
//
// The total count is computed from the same base query before pagination:
//
//	WITH <tag CTE>
//	SELECT COUNT(*)
//	FROM transactions t
//	LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
//	WHERE <filter>
//
// GetTransaction returns a single transaction by ID with tags:
//
//	WITH <tag CTE>
//	SELECT t.id, t.hash, t.date, t.account, t.kind,
//	       t.amount, t.content, t.info, t.recipient,
//	       COALESCE(rt.tags, []::VARCHAR[]) AS tags
//	FROM transactions t
//	LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
//	WHERE t.id = ?
//
// GetTransactionByHash returns a single transaction by hash with tags:
//
//	WITH <tag CTE>
//	SELECT t.id, t.hash, t.date, t.account, t.kind,
//	       t.amount, t.content, t.info, t.recipient,
//	       COALESCE(rt.tags, []::VARCHAR[]) AS tags
//	FROM transactions t
//	LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
//	WHERE t.hash = ?
//
// CreateTransaction inserts a transaction and returns its ID:
//
//	INSERT INTO transactions (hash, date, account, kind, amount, content, info, recipient)
//	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
//	RETURNING id
//
// UpdateTransaction modifies a transaction by ID:
//
//	UPDATE transactions
//	SET hash = ?, date = ?, account = ?, kind = ?, amount = ?,
//	    content = ?, info = ?, recipient = ?
//	WHERE id = ?
//
// DeleteTransaction removes a transaction by ID:
//
//	DELETE FROM transactions WHERE id = ?
//
// # Summary Queries
//
// GetSummaryOverview returns aggregate stats for all transactions matching a filter:
//
//	SELECT COUNT(*),
//	       COALESCE(SUM(CASE WHEN kind='debit' THEN amount ELSE 0 END), 0),
//	       COALESCE(SUM(CASE WHEN kind='credit' THEN amount ELSE 0 END), 0),
//	       COUNT(DISTINCT account)
//	FROM transactions t
//	WHERE <filter>
//
// It also calls countUniqueTags to compute the number of distinct tags across
// all matching transactions.
//
// countUniqueTags counts distinct tags for transactions matching a filter:
//
//	WITH <tag CTE>
//	SELECT COUNT(DISTINCT tag)
//	FROM rule_tags, UNNEST(tags) AS t(tag)
//	WHERE transaction_id IN (SELECT t.id FROM transactions t WHERE <filter>)
//
// GetTagSummary returns per-tag debit/credit totals and transaction counts.
// It extends the tag CTE with an unnested layer and joins back to transactions:
//
//	WITH <tag CTE>,
//	unnested AS (
//	    SELECT transaction_id, UNNEST(tags) AS tag
//	    FROM rule_tags
//	)
//	SELECT u.tag,
//	       COALESCE(SUM(CASE WHEN t.kind='debit' THEN t.amount ELSE 0 END), 0) AS total_debit,
//	       COALESCE(SUM(CASE WHEN t.kind='credit' THEN t.amount ELSE 0 END), 0) AS total_credit,
//	       COUNT(*) AS count
//	FROM unnested u
//	JOIN transactions t ON t.id = u.transaction_id
//	WHERE <filter>
//	GROUP BY u.tag
//	ORDER BY total_debit DESC
//
// GetBalanceTrend returns monthly debit/credit totals with a cumulative balance:
//
//	SELECT strftime(date, '%Y-%m') AS month,
//	       COALESCE(SUM(CASE WHEN kind='debit' THEN amount ELSE 0 END), 0) AS debit,
//	       COALESCE(SUM(CASE WHEN kind='credit' THEN amount ELSE 0 END), 0) AS credit
//	FROM transactions t
//	WHERE <filter>
//	GROUP BY strftime(date, '%Y-%m')
//	ORDER BY month ASC
//
// The cumulative balance is computed in Go by summing (credit - debit) across
// rows in month order.
package store
