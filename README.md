# Ducket

Personal finance tracker built with Go and DuckDB.

Ducket imports bank transactions from CSV and Excel files, stores them in an embedded DuckDB database, and lets you categorize spending through dynamic tagging rules. Tags are never stored on transactions — they are computed at query time by evaluating each rule's filter expression, so changing a rule retroactively re-tags all existing transactions without any data migration.

## Why

Most personal finance tools either require manual tagging or store tags as static data. When you realize your categorization was wrong, you have to re-tag everything. Ducket takes a different approach: rules are filter expressions evaluated at query time via SQL CTEs. Add a rule today and every past transaction that matches it is instantly tagged.

DuckDB was chosen because it's an embedded columnar database — no server, no network, single file — with native array types (`VARCHAR[]`) and analytical query performance that makes the dynamic tag CTE fast even with many rules.

## Quick Start

```bash
go build -o ducket .

# Run migrations and start the server
./ducket run --db-uri ./ducket.db --server-port 8080

# Import transactions
curl -F 'files=@january.xlsx' localhost:8080/api/v1/transactions/import

# Create a tagging rule
curl -X POST localhost:8080/api/v1/rules \
  -H 'Content-Type: application/json' \
  -d '{"name": "groceries", "filter": "content ~ /KAUFLAND|LIDL|PENNY/", "tags": ["food", "groceries"]}'

# List debit transactions (tags are computed on the fly)
curl 'localhost:8080/api/v1/transactions?filter=kind%20%3D%20%27debit%27&limit=20'

# Get spending breakdown by tag
curl localhost:8080/api/v1/summary/by-tag
```

## Features

- **Transaction CRUD** with SHA-256 hash deduplication
- **Bulk import** from CSV and Excel (.xlsx) files
- **Dynamic tagging** — rules define filter expressions + tag arrays, evaluated at query time
- **Filter DSL** — comparison, regex, IN/NOT IN, boolean operators on any field
- **Sorting and pagination** — server-side, with total count
- **Dashboard summaries** — overview, per-tag breakdown, monthly balance trend
- **Tag filtering** — query only transactions matching specific tags

## API

All endpoints under `/api/v1`:

| Method | Path | Description |
|--------|------|-------------|
| GET | /transactions | List transactions (filter, sort, paginate, tag filter) |
| POST | /transactions | Create transaction |
| GET | /transactions/{id} | Get transaction |
| PUT | /transactions/{id} | Update transaction |
| DELETE | /transactions/{id} | Delete transaction |
| POST | /transactions/import | Import from CSV/Excel |
| GET | /rules | List rules |
| POST | /rules | Create rule |
| GET | /rules/{id} | Get rule |
| PUT | /rules/{id} | Update rule |
| DELETE | /rules/{id} | Delete rule |
| GET | /summary/overview | Aggregate totals |
| GET | /summary/by-tag | Per-tag spending |
| GET | /summary/balance-trend | Monthly trend |

See [doc/api.md](doc/api.md) for request/response examples.

## Stack

| | |
|---|---|
| Go | Single binary, fast builds |
| DuckDB | Embedded columnar DB, native arrays, analytical queries |
| Gin | HTTP router with generated handlers from OpenAPI |
| Squirrel | SQL query builder (no ORM) |

## Documentation

- [doc/architecture.md](doc/architecture.md) — design overview, schema, configuration
- [doc/api.md](doc/api.md) — full API reference with examples
- [doc/store.md](doc/store.md) — store internals, every SQL query explained

## License

MIT
