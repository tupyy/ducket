# Store Layer

The store package implements the DuckDB persistence layer. It uses the `squirrel` query builder with `?` placeholder format and routes all queries through a `queryInterceptor` that handles context-based transactions and auto-checkpointing.

## Why DuckDB

DuckDB is an embedded columnar database — no server process, no network, single file. It handles analytical queries (aggregations, window functions) efficiently and supports `VARCHAR[]` arrays natively, which is used for rule tags. The tradeoff is single-writer concurrency, managed by `SetMaxOpenConns(1)` and a mutex in the interceptor.

## Why Dynamic Tags

Tags are not stored on transactions. Instead, rules define filter expressions and tag arrays. When transactions are queried, every rule's filter is evaluated against the transactions table via a CTE, and matching transactions receive the rule's tags. This means:

- Adding or changing a rule retroactively re-tags all existing transactions
- No data migration needed when rules change
- No junction tables or denormalized tag columns
- The cost is paid at query time, not write time

## Query Interceptor

All queries go through `queryInterceptor`, which:

1. Checks context for a `*sql.Tx` — if present, routes the query through it
2. For `ExecContext` outside a transaction: acquires a mutex, executes the write, then runs `FORCE CHECKPOINT`
3. Logs all queries at debug level

This ensures writes are serialized and the DuckDB WAL is flushed after every standalone write.

## Transactor

`WithTx` begins a transaction, stores it in the context, and passes the context to the callback. If the callback returns an error, the transaction is rolled back; otherwise it's committed. Nested `WithTx` calls detect the existing transaction in the context and reuse it.

## Rule Queries

### ListRules

Returns rules matching an optional filter with pagination.

```sql
SELECT id, name, filter, tags, created_at
FROM rules
WHERE <filter>                  -- omitted when filter is empty
ORDER BY id ASC
LIMIT ? OFFSET ?                -- omitted when 0
```

**Example:** `ListRules(ctx, "name ~ /grocery/", 10, 0)`

```sql
SELECT id, name, filter, tags, created_at
FROM rules
WHERE regexp_matches(name, ?)
ORDER BY id ASC
LIMIT 10
-- args: ['grocery']
```

### GetRule

Returns a single rule by ID.

```sql
SELECT id, name, filter, tags, created_at
FROM rules
WHERE id = ?
```

Returns `ResourceNotFoundError` if no rows.

### CreateRule

Inserts a rule and returns its generated ID. Tags are serialized as a DuckDB array literal.

```sql
INSERT INTO rules (name, filter, tags)
VALUES (?, ?, ['food', 'groceries'])
RETURNING id
```

The service layer validates the filter expression before calling `CreateRule`.

### UpdateRule

Updates a rule by ID. Returns `ResourceNotFoundError` if `RowsAffected() == 0`.

```sql
UPDATE rules
SET name = ?, filter = ?, tags = ['food', 'groceries']
WHERE id = ?
```

### DeleteRule

Deletes a rule by ID. Returns `ResourceNotFoundError` if `RowsAffected() == 0`.

```sql
DELETE FROM rules WHERE id = ?
```

## Transaction Queries

### Tag CTE Construction

Before any transaction read, the store loads all rules and builds a CTE. Given rules:

```
Rule 1: filter="content ~ /KAUFLAND|LIDL/"  tags=['food', 'groceries']
Rule 2: filter="content ~ /SHELL|OMV/"      tags=['transport', 'fuel']
```

The CTE becomes:

```sql
WITH rule_matches AS (
    SELECT t.id AS transaction_id, ['food', 'groceries'] AS tags
    FROM transactions t
    WHERE regexp_matches(t.content, ?)            -- args: 'KAUFLAND|LIDL'
    UNION ALL
    SELECT t.id AS transaction_id, ['transport', 'fuel'] AS tags
    FROM transactions t
    WHERE regexp_matches(t.content, ?)            -- args: 'SHELL|OMV'
),
rule_tags AS (
    SELECT transaction_id,
           list_distinct(flatten(list(tags))) AS tags
    FROM rule_matches
    GROUP BY transaction_id
)
```

When there are no rules, the CTE is omitted entirely and tags default to `[]::VARCHAR[]`.

### Tag Filtering

When the `tags` query parameter is provided (e.g. `tags=food&tags=groceries`):

1. Only rules whose tags contain **all** requested tags are included in the CTE
2. An `INNER JOIN` replaces the `LEFT JOIN`, so only transactions with matching tags appear

### ListTransactions

Returns paginated transactions with computed tags and a total count. Two queries are executed:

**Count query** (before pagination):

```sql
WITH <tag CTE>
SELECT COUNT(*)
FROM transactions t
LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
WHERE <filter>
```

**Data query:**

```sql
WITH <tag CTE>
SELECT t.id, t.hash, t.date, t.account, t.kind,
       t.amount, t.content, t.info, t.recipient,
       COALESCE(rt.tags, []::VARCHAR[]) AS tags
FROM transactions t
LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
WHERE <filter>
ORDER BY t.date DESC, t.id DESC        -- default sort
LIMIT ? OFFSET ?
```

**Sorting:** Configurable via `SortParam` on fields `date`, `amount`, `kind`, `account`. Each can be `ASC` or `DESC`. `t.id DESC` is always appended as a tiebreaker.

**Example:** `ListTransactions(ctx, "kind = 'debit' and amount > 100", nil, []SortParam{{Field: "amount", Desc: true}}, 20, 0)`

```sql
WITH rule_matches AS (...), rule_tags AS (...)
SELECT t.id, t.hash, t.date, t.account, t.kind,
       t.amount, t.content, t.info, t.recipient,
       COALESCE(rt.tags, []::VARCHAR[]) AS tags
FROM transactions t
LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
WHERE t.kind = ? AND t.amount > ?
ORDER BY t.amount DESC, t.id DESC
LIMIT 20
-- args: [<cte_args...>, 'debit', 100]
```

### GetTransaction

Returns a single transaction by ID with computed tags.

```sql
WITH <tag CTE>
SELECT t.id, t.hash, t.date, t.account, t.kind,
       t.amount, t.content, t.info, t.recipient,
       COALESCE(rt.tags, []::VARCHAR[]) AS tags
FROM transactions t
LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
WHERE t.id = ?
```

Returns `ResourceNotFoundError` if no rows.

### GetTransactionByHash

Returns a single transaction by hash with computed tags. Returns `nil, nil` if not found (used internally for deduplication).

```sql
WITH <tag CTE>
SELECT t.id, t.hash, t.date, t.account, t.kind,
       t.amount, t.content, t.info, t.recipient,
       COALESCE(rt.tags, []::VARCHAR[]) AS tags
FROM transactions t
LEFT JOIN rule_tags rt ON rt.transaction_id = t.id
WHERE t.hash = ?
```

### CreateTransaction

Inserts a transaction and returns its generated ID. The hash is computed by the entity layer from `kind|account|date|amount|content` using SHA-256. The service layer checks for hash uniqueness within a transaction before inserting.

```sql
INSERT INTO transactions (hash, date, account, kind, amount, content, info, recipient)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id
```

### UpdateTransaction

Updates a transaction by ID. The hash is recomputed by the service layer before calling the store. Returns `ResourceNotFoundError` if `RowsAffected() == 0`.

```sql
UPDATE transactions
SET hash = ?, date = ?, account = ?, kind = ?, amount = ?,
    content = ?, info = ?, recipient = ?
WHERE id = ?
```

### DeleteTransaction

Deletes a transaction by ID. Returns `ResourceNotFoundError` if `RowsAffected() == 0`.

```sql
DELETE FROM transactions WHERE id = ?
```

## Summary Queries

### GetSummaryOverview

Returns aggregate statistics for transactions matching a filter. Executes two queries: one for transaction aggregates and one for unique tag count.

**Transaction aggregates:**

```sql
SELECT COUNT(*),
       COALESCE(SUM(CASE WHEN kind='debit' THEN amount ELSE 0 END), 0),
       COALESCE(SUM(CASE WHEN kind='credit' THEN amount ELSE 0 END), 0),
       COUNT(DISTINCT account)
FROM transactions t
WHERE <filter>
```

**Unique tag count** (uses the tag CTE):

```sql
WITH <tag CTE>
SELECT COUNT(DISTINCT tag)
FROM rule_tags, UNNEST(tags) AS t(tag)
WHERE transaction_id IN (SELECT t.id FROM transactions t WHERE <filter>)
```

`balance` is computed in Go as `total_credit - total_debit`.

### GetTagSummary

Returns per-tag debit/credit totals. Extends the tag CTE with an `unnested` layer that explodes the tag arrays, then joins back to transactions for amount aggregation.

```sql
WITH <tag CTE>,
unnested AS (
    SELECT transaction_id, UNNEST(tags) AS tag
    FROM rule_tags
)
SELECT u.tag,
       COALESCE(SUM(CASE WHEN t.kind='debit' THEN t.amount ELSE 0 END), 0) AS total_debit,
       COALESCE(SUM(CASE WHEN t.kind='credit' THEN t.amount ELSE 0 END), 0) AS total_credit,
       COUNT(*) AS count
FROM unnested u
JOIN transactions t ON t.id = u.transaction_id
WHERE <filter>
GROUP BY u.tag
ORDER BY total_debit DESC
```

Returns an empty slice when no rules exist.

### GetBalanceTrend

Returns monthly debit/credit totals. The cumulative balance is computed in Go by summing `credit - debit` across months in order.

```sql
SELECT strftime(date, '%Y-%m') AS month,
       COALESCE(SUM(CASE WHEN kind='debit' THEN amount ELSE 0 END), 0) AS debit,
       COALESCE(SUM(CASE WHEN kind='credit' THEN amount ELSE 0 END), 0) AS credit
FROM transactions t
WHERE <filter>
GROUP BY strftime(date, '%Y-%m')
ORDER BY month ASC
```

## Argument Ordering

All queries using the tag CTE follow the same argument ordering pattern:

1. CTE arguments (from rule filter parsing) come first
2. SELECT/WHERE arguments (from the user's filter) come second

This is because the CTE is prepended as a raw string before the squirrel-generated query:

```go
query = "WITH " + tagCTE + " " + query
allArgs = append(tagArgs, selectArgs...)
```
