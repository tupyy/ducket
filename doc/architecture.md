# Finante Architecture

Personal finance application built with Go, DuckDB, and Gin.

## Overview

Finante tracks bank transactions and applies user-defined tagging rules to categorize spending. Tags are never stored on transactions — they are computed dynamically at query time by evaluating each rule's filter expression against the transactions table.

```
HTTP request
  |
  v
Gin router (/api/v1)
  |
  v
Handler (handlers/)         -- parse HTTP, map to entities, call services
  |
  v
Service (services/)         -- validation, deduplication, transactions
  |
  v
Store (store/)              -- squirrel query builder, DuckDB
  |
  v
DuckDB (single file)
```

## Database

DuckDB is used as an embedded single-file database. The connection is opened with `SetMaxOpenConns(1)` because DuckDB is single-writer. A query interceptor serializes writes and runs `FORCE CHECKPOINT` after every write outside a transaction.

### Schema

Two tables, no junction tables, no stored tags:

```sql
CREATE SEQUENCE IF NOT EXISTS rules_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS transactions_id_seq START 1;

CREATE TABLE rules (
    id          INTEGER PRIMARY KEY DEFAULT nextval('rules_id_seq'),
    name        VARCHAR NOT NULL UNIQUE,
    filter      VARCHAR NOT NULL,
    tags        VARCHAR[] NOT NULL,
    created_at  TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE transactions (
    id          INTEGER PRIMARY KEY DEFAULT nextval('transactions_id_seq'),
    hash        VARCHAR NOT NULL UNIQUE,
    date        DATE NOT NULL,
    account     BIGINT NOT NULL,
    kind        VARCHAR NOT NULL,       -- 'debit' or 'credit'
    amount      DECIMAL(12,2) NOT NULL,
    content     VARCHAR NOT NULL,
    info        VARCHAR,
    recipient   VARCHAR,
    created_at  TIMESTAMP DEFAULT current_timestamp
);
```

Migrations are managed by a custom runner (not goose) using embedded SQL files in `internal/store/migrations/sql/`. Applied versions are tracked in a `schema_migrations` table.

### DuckDB Array Handling

Tags are stored as `VARCHAR[]` in DuckDB. The driver returns arrays as `[]interface{}` on scan. The `StringArray` type implements `sql.Scanner` for this conversion. For writes, tags are serialized as DuckDB array literals (e.g. `['food', 'groceries']`) embedded in the SQL via `sq.Expr`.

## Transaction Support

The `Store.WithTx` method wraps a function in a DuckDB transaction. The `*sql.Tx` is stored in the context and automatically picked up by the query interceptor, so all store methods called within the closure participate in the same transaction. Nested `WithTx` calls reuse the existing transaction.

```go
store.WithTx(ctx, func(ctx context.Context) error {
    id, err := store.CreateTransaction(ctx, txn)
    if err != nil {
        return err // triggers rollback
    }
    return nil // triggers commit
})
```

## Filter DSL

All list and summary endpoints accept a `filter` query parameter parsed by `github.com/kubev2v/assisted-migration-agent/pkg/filter`. The DSL supports comparison, regex, IN/NOT IN, and boolean operators. Filters are parsed into parameterized `sq.Sqlizer` clauses — no raw string interpolation.

**Transaction filter fields:** `id`, `hash`, `date`, `account`, `kind`, `amount`, `content`, `info`, `recipient`

**Rule filter fields:** `id`, `name`, `filter`

**Examples:**
- `date >= '2025-01-01' and date < '2025-02-01' and kind = 'debit'`
- `content ~ /KAUFLAND|LIDL/`
- `amount > 1000 and kind = 'debit'`
- `name ~ /grocery/`

## Dynamic Tag CTE

When querying transactions (list, get, summary), the store loads all rules and builds a CTE that evaluates each rule's filter against the transactions table. Each matching rule contributes its tag array. The results are aggregated per transaction into a merged, deduplicated tag array.

Given two rules:

```
Rule 1: filter="content ~ /KAUFLAND|LIDL/"  tags=['food', 'groceries']
Rule 2: filter="content ~ /SHELL|OMV/"      tags=['transport', 'fuel']
```

The generated CTE is:

```sql
WITH rule_matches AS (
    SELECT t.id AS transaction_id, ['food', 'groceries'] AS tags
    FROM transactions t
    WHERE regexp_matches(t.content, ?)
    UNION ALL
    SELECT t.id AS transaction_id, ['transport', 'fuel'] AS tags
    FROM transactions t
    WHERE regexp_matches(t.content, ?)
),
rule_tags AS (
    SELECT transaction_id,
           list_distinct(flatten(list(tags))) AS tags
    FROM rule_matches
    GROUP BY transaction_id
)
```

When there are no rules, the CTE is omitted and tags default to `[]::VARCHAR[]`.

When filtering by tags (via the `tags` query parameter), only rules whose tags contain all requested tags are included, and an `INNER JOIN` replaces the `LEFT JOIN` so only matching transactions are returned.

## Error Handling

Three typed errors in `pkg/errors`:

| Error | HTTP Status | When |
|-------|-------------|------|
| `ResourceNotFoundError` | 404 | Get/Update/Delete with nonexistent ID |
| `DuplicateResourceError` | 409 | Create transaction with existing hash |
| `ValidationError` | 400 | Invalid transaction kind, invalid filter expression |

## Configuration

Environment variables prefixed with `FINANTE_` or CLI flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--db-uri` | | Path to DuckDB file (e.g. `./finante.db` or `:memory:`) |
| `--server-port` | `8080` | HTTP listen port |
| `--server-gin-mode` | | `release` or `debug` |
| `--server-mode` | `dev` | `dev` or `prod` (prod requires `--statics-folder`) |
| `--statics-folder` | | Path to static frontend files |

## File Import

Transactions can be bulk-imported from CSV or Excel files via `POST /api/v1/transactions/import`. The `pkg/reader` package detects format by extension (`.csv`, `.xlsx`, `.xls`) and parses transactions. Duplicates (by hash) are skipped, not rejected.

## Project Structure

```
cmd/
  serve.go              -- cobra command: wire DB, services, handlers, HTTP server
  migrate.go            -- cobra command: run migrations

api/v1/
  openapi.yaml          -- OpenAPI 3.0 spec
  types.gen.go          -- generated types (oapi-codegen)
  spec.gen.go           -- generated server interface + router
  extension.go          -- entity-to-API type converters

internal/
  config/               -- configuration structs (optgen)
  entity/               -- domain types: Transaction, Rule
  handlers/             -- HTTP handlers (Gin)
  server/               -- HTTP server setup
  services/             -- business logic, validation
  store/                -- DuckDB persistence layer
    migrations/         -- custom migration runner + SQL files

pkg/
  errors/               -- typed errors (NotFound, Duplicate, Validation)
  logger/               -- zap logger setup
  reader/               -- CSV and Excel transaction parsers
```
