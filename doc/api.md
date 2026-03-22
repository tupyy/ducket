# API Reference

All endpoints are served under `/api/v1`. Responses use `application/json`.

## Transactions

### List Transactions

```
GET /api/v1/transactions
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `filter` | string | Filter DSL expression (e.g. `kind = 'debit' and amount > 100`) |
| `tags` | string[] | Filter by tags — only transactions matching all specified tags are returned |
| `sort` | string[] | Sort fields with direction (e.g. `date:desc`, `amount:asc`). Valid fields: `date`, `amount`, `kind`, `account` |
| `limit` | int | Page size (default: 100, max: 1000) |
| `offset` | int | Number of items to skip |

**Response (200):**

```json
{
  "items": [
    {
      "id": 1,
      "hash": "a1b2c3...",
      "date": "2025-01-15T00:00:00Z",
      "account": 1234567890,
      "kind": "debit",
      "amount": 42.50,
      "content": "KAUFLAND TIMISOARA",
      "info": "Card payment",
      "recipient": "KAUFLAND SRL",
      "tags": ["food", "groceries"]
    }
  ],
  "total": 1523
}
```

**Examples:**

```bash
# All debit transactions in January 2025
curl 'localhost:8080/api/v1/transactions?filter=date%20%3E%3D%20%272025-01-01%27%20and%20date%20%3C%20%272025-02-01%27%20and%20kind%20%3D%20%27debit%27&limit=50'

# Transactions matching a regex, sorted by amount descending
curl 'localhost:8080/api/v1/transactions?filter=content%20~%20/KAUFLAND|LIDL/&sort=amount:desc'

# Filter by tags
curl 'localhost:8080/api/v1/transactions?tags=food&tags=groceries&limit=20'
```

### Get Transaction

```
GET /api/v1/transactions/{id}
```

**Response (200):**

```json
{
  "id": 1,
  "hash": "a1b2c3...",
  "date": "2025-01-15T00:00:00Z",
  "account": 1234567890,
  "kind": "debit",
  "amount": 42.50,
  "content": "KAUFLAND TIMISOARA",
  "info": "Card payment",
  "recipient": "KAUFLAND SRL",
  "tags": ["food", "groceries"]
}
```

**Response (404):** `{"error": "transaction 99 not found"}`

### Create Transaction

```
POST /api/v1/transactions
Content-Type: application/json
```

**Request Body:**

```json
{
  "kind": "debit",
  "account": 1234567890,
  "date": "2025-01-15T00:00:00Z",
  "amount": 42.50,
  "content": "KAUFLAND TIMISOARA",
  "info": "Card payment",
  "recipient": "KAUFLAND SRL"
}
```

All fields except `info` and `recipient` are required. The `kind` must be `debit` or `credit`. A SHA-256 hash is computed from `kind|account|date|amount|content` — duplicate hashes are rejected with 409.

**Response (201):** The created transaction (same shape as Get).

**Response (409):** `{"error": "transaction with hash \"a1b2c3...\" already exists"}`

### Update Transaction

```
PUT /api/v1/transactions/{id}
Content-Type: application/json
```

**Request Body:** Same as Create. The hash is recomputed from the new values.

**Response (204):** No content.

**Response (404):** `{"error": "transaction 99 not found"}`

### Delete Transaction

```
DELETE /api/v1/transactions/{id}
```

**Response (204):** No content.

**Response (404):** `{"error": "transaction 99 not found"}`

### Import Transactions

```
POST /api/v1/transactions/import
Content-Type: multipart/form-data
```

**Form Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `files` | file[] | CSV or Excel files (.csv, .xlsx, .xls) |
| `account` | int | Account number to assign to CSV imports (Excel files contain their own) |

Duplicate transactions (same hash) are skipped, not rejected.

**Response (200):**

```json
{
  "files": [
    {
      "filename": "january.xlsx",
      "created": 150,
      "skipped": 3,
      "errors": 0
    }
  ],
  "message": "imported 150 transactions from 1 file(s)"
}
```

## Rules

Rules define filters that tag transactions dynamically. A rule's `filter` field uses the same DSL as the transaction list endpoint. When transactions are queried, each rule's filter is evaluated and matching transactions receive the rule's tags.

### List Rules

```
GET /api/v1/rules
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `filter` | string | Filter DSL expression (e.g. `name ~ /grocery/`) |
| `limit` | int | Page size (default: 100, max: 1000) |
| `offset` | int | Number of items to skip |

**Response (200):**

```json
[
  {
    "id": 1,
    "name": "groceries",
    "filter": "content ~ /KAUFLAND|LIDL|PENNY/",
    "tags": ["food", "groceries"],
    "created_at": "2025-01-01T12:00:00Z"
  },
  {
    "id": 2,
    "name": "fuel",
    "filter": "content ~ /SHELL|OMV|PETROM/",
    "tags": ["transport", "fuel"],
    "created_at": "2025-01-01T12:05:00Z"
  }
]
```

### Get Rule

```
GET /api/v1/rules/{id}
```

**Response (200):** Single rule object.

**Response (404):** `{"error": "rule 99 not found"}`

### Create Rule

```
POST /api/v1/rules
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "groceries",
  "filter": "content ~ /KAUFLAND|LIDL|PENNY/",
  "tags": ["food", "groceries"]
}
```

The `filter` field is validated by parsing it as a transaction filter expression. Invalid filters are rejected with 400.

**Response (201):** The created rule.

**Response (400):** `{"error": "invalid filter: unknown transaction field: foo"}`

### Update Rule

```
PUT /api/v1/rules/{id}
Content-Type: application/json
```

**Request Body:** Same as Create.

**Response (204):** No content.

**Response (404):** `{"error": "rule 99 not found"}`

### Delete Rule

```
DELETE /api/v1/rules/{id}
```

**Response (204):** No content.

**Response (404):** `{"error": "rule 99 not found"}`

## Summary

Summary endpoints provide aggregate views over transactions. All accept an optional `filter` query parameter.

### Overview

```
GET /api/v1/summary/overview
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `filter` | string | Filter DSL expression |

**Response (200):**

```json
{
  "total_transactions": 1523,
  "total_debit": 45230.50,
  "total_credit": 52100.00,
  "balance": 6869.50,
  "unique_accounts": 3,
  "unique_tags": 12
}
```

`balance` = `total_credit` - `total_debit`.

**Example:**

```bash
# Overview for 2025 only
curl 'localhost:8080/api/v1/summary/overview?filter=date%20%3E%3D%20%272025-01-01%27'
```

### By Tag

```
GET /api/v1/summary/by-tag
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `filter` | string | Filter DSL expression |

Returns per-tag spending breakdown. Each tag appears once with its total debit, credit, and transaction count. A transaction with multiple tags contributes to each tag independently.

**Response (200):**

```json
[
  {
    "tag": "food",
    "total_debit": 8500.00,
    "total_credit": 0.00,
    "count": 230
  },
  {
    "tag": "transport",
    "total_debit": 3200.00,
    "total_credit": 0.00,
    "count": 85
  }
]
```

Results are ordered by `total_debit` descending.

**Example:**

```bash
# Tag breakdown for debit transactions only
curl 'localhost:8080/api/v1/summary/by-tag?filter=kind%20%3D%20%27debit%27'
```

### Balance Trend

```
GET /api/v1/summary/balance-trend
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `filter` | string | Filter DSL expression |

Returns monthly debit/credit totals with a running cumulative balance.

**Response (200):**

```json
[
  {
    "month": "2025-01",
    "debit": 4200.00,
    "credit": 5100.00,
    "balance": 900.00
  },
  {
    "month": "2025-02",
    "debit": 3800.00,
    "credit": 4900.00,
    "balance": 2000.00
  }
]
```

`balance` is cumulative: each month's balance = previous balance + (credit - debit).

**Example:**

```bash
# Trend for a specific account
curl 'localhost:8080/api/v1/summary/balance-trend?filter=account%20%3D%201234567890'
```

## Filter DSL Reference

The filter DSL is used in the `filter` query parameter across all list and summary endpoints.

### Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Equals | `kind = 'debit'` |
| `!=` | Not equals | `kind != 'credit'` |
| `>`, `>=`, `<`, `<=` | Comparison | `amount > 1000` |
| `~` | Regex match | `content ~ /KAUFLAND\|LIDL/` |
| `and` | Logical AND | `kind = 'debit' and amount > 100` |
| `or` | Logical OR | `account = 111 or account = 222` |
| `in` | Value in set | `kind in ('debit', 'credit')` |
| `not in` | Value not in set | `account not in (111, 222)` |

### Transaction Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Transaction ID |
| `hash` | string | SHA-256 deduplication hash |
| `date` | date | Transaction date |
| `account` | int | Bank account number |
| `kind` | string | `debit` or `credit` |
| `amount` | decimal | Transaction amount |
| `content` | string | Transaction description |
| `info` | string | Additional info |
| `recipient` | string | Payment recipient |

### Rule Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Rule ID |
| `name` | string | Rule name |
| `filter` | string | Rule's filter expression |

### Examples

```
# All debit transactions over 500 in January 2025
date >= '2025-01-01' and date < '2025-02-01' and kind = 'debit' and amount > 500

# Grocery stores by regex
content ~ /KAUFLAND|LIDL|PENNY|AUCHAN/

# Specific accounts
account in (1234567890, 9876543210)

# Rules containing "food" in name
name ~ /food/
```
