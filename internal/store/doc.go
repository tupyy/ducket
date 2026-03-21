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
// Examples:
//
//	// All debit transactions in January 2025
//	store.ListTransactions(ctx, "date >= '2025-01-01' and date < '2025-02-01' and kind = 'debit'", 50, 0)
//
//	// Transactions matching a regex on content
//	store.ListTransactions(ctx, "content ~ /KAUFLAND|LIDL/", 0, 0)
//
//	// High-value transactions
//	store.ListTransactions(ctx, "amount > 1000 and kind = 'debit'", 20, 0)
//
//	// Rules by name pattern
//	store.ListRules(ctx, "name ~ /grocery/", 0, 0)
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
// The generated SQL for ListTransactions(ctx, "kind = 'debit'", 20, 0) is:
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
//	SELECT t.id, t.hash, t.date, t.account, t.kind,
//	       t.amount, t.content, t.info, t.recipient,
//	       COALESCE(rt.tags, []::VARCHAR[]) AS tags
//	FROM transactions t
//	LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
//	WHERE t.kind = ?
//	ORDER BY t.date DESC, t.id DESC
//	LIMIT 20
//	-- args: ['KAUFLAND|LIDL', 'SHELL|OMV', 'debit']
//
// When there are no rules, the CTE is omitted and tags default to an empty array.
//
// # CRUD Queries
//
// Rule insert:
//
//	INSERT INTO rules (name, filter, tags)
//	VALUES (?, ?, ['food', 'groceries'])
//	RETURNING id
//
// Rule update:
//
//	UPDATE rules SET name = ?, filter = ?, tags = ['food', 'groceries']
//	WHERE id = ?
//
// Rule delete:
//
//	DELETE FROM rules WHERE id = ?
//
// Transaction insert:
//
//	INSERT INTO transactions (hash, date, account, kind, amount, content, info, recipient)
//	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
//	RETURNING id
//
// Transaction update:
//
//	UPDATE transactions
//	SET date = ?, account = ?, kind = ?, amount = ?, content = ?, info = ?, recipient = ?
//	WHERE id = ?
//
// Transaction delete:
//
//	DELETE FROM transactions WHERE id = ?
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
package store
