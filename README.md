# Finante - Personal Financial Management System

A comprehensive personal financial management system built with Go and React, designed to help you track, categorize, and analyze your financial transactions through automated rule-based labeling.

## 🏗️ Architecture Overview

Finante follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React UI      │    │   REST API      │    │   PostgreSQL    │
│   (Frontend)    │◄──►│   (Go Backend)  │◄──►│   (Database)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Backend Components

- **HTTP Server**: Gin-based REST API with middleware for logging, CORS, and database injection
- **Service Layer**: Business logic for transactions, rules, and labels management
- **Repository Layer**: Database operations with PostgreSQL using pgx driver
- **Rule Engine**: Pattern-based transaction categorization system with key-value labels
- **File Parser**: Excel file import functionality for bulk transaction loading

## 🚀 Features

### Core Functionality
- **Transaction Management**: Create, read, update, and delete financial transactions
- **Rule-Based Categorization**: Automatic transaction labeling based on configurable rules
- **Label Management**: Organize transactions with structured key-value labels
- **Excel Import**: Bulk import transactions from Excel spreadsheets
- **Financial Analytics**: Transaction statistics and summaries

### Technical Features
- **RESTful API**: Well-structured API endpoints with JSON responses
- **Service Architecture**: Clean separation between Create and Update operations with proper error handling
- **Resource Management**: Comprehensive error handling for existing/non-existing resources
- **Database Migrations**: Automated schema management with goose
- **Relationship Integrity**: Automatic management of complex relationships between entities
- **Graceful Shutdown**: Proper server lifecycle management
- **Structured Logging**: Comprehensive logging with zap
- **Configuration Management**: Environment-based configuration
- **Connection Pooling**: Optimized database connection management
- **Test Coverage**: Comprehensive unit and integration tests with proper database cleanup

## 📋 Prerequisites

- **Go**: Version 1.19 or higher
- **Node.js**: Version 16 or higher (for frontend)
- **PostgreSQL**: Version 12 or higher
- **Make**: For build automation

## 🛠️ Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd finante
```

### 2. Database Setup
```bash
# Create PostgreSQL database
createdb finante

# Set database connection string
export DATABASE_URI="postgres://username:password@localhost:5432/finante?sslmode=disable"
```

### 3. Backend Setup
```bash
# Install Go dependencies
go mod download

# Run database migrations
make migrate

# Build the application
make build

# Start the server
./finante serve --db-conn-uri="$DATABASE_URI" --server-port=8080
```

### 4. Frontend Setup
```bash
cd ui
npm install
npm run build  # For production
npm start      # For development
```

## 🔧 Configuration

The application can be configured via command-line flags or environment variables:

### Database Configuration
- `--db-conn-uri`: PostgreSQL connection string
- `--db-ssl-mode`: Enable/disable SSL mode (default: false)

### Server Configuration  
- `--server-port`: HTTP server port (default: 8080)
- `--server-gin-mode`: Gin mode - `release` or `debug` (default: release)

### Logging Configuration
- `--log-level`: Log level - `debug`, `info`, `warn`, `error` (default: info)
- `--log-format`: Log format - `json` or `text` (default: json)

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Transactions
- `GET /transactions` - List all transactions with optional filtering
- `GET /transactions/:id` - Get specific transaction
- `POST /transactions` - Create new transaction (returns error if transaction already exists)
- `PUT /transactions/:id` - Update existing transaction (creates new if doesn't exist)
- `DELETE /transactions/:id` - Delete transaction

#### Transaction Labels
- `GET /transactions/:id/labels` - Get all labels for a transaction
- `POST /transactions/:id/labels` - Add a label to a transaction
- `DELETE /transactions/:id/labels` - Remove all labels from a transaction
- `DELETE /transactions/:id/labels/:labelId` - Remove a specific label from a transaction

#### Rules
- `GET /rules` - List all rules
- `GET /rules/:name` - Get specific rule
- `POST /rules` - Create new rule (returns error if rule already exists)
- `PUT /rules/:name` - Update existing rule (returns error if rule doesn't exist)
- `DELETE /rules/:name` - Delete rule

#### Labels
- `GET /labels` - List all labels with transaction counts
- `POST /labels` - Create new label
- `DELETE /labels/:id` - Delete label

#### Analytics
- `GET /summary` - Get transaction statistics and summaries

### Request/Response Examples

#### Transactions

##### List All Transactions
```bash
# Get all transactions
curl -X GET http://localhost:8080/api/v1/transactions

# Get transactions with filters
curl -X GET "http://localhost:8080/api/v1/transactions?start=2024-01-01&end=2024-01-31&labels=category:food&limit=10&offset=0"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "date": "2024-01-15T00:00:00Z",
      "kind": "debit",
      "amount": 45.50,
      "content": "grocery store purchase",
      "hash": "abc123def456",
      "labels": [
        {
          "key": "category",
          "value": "food",
          "href": "/api/v1/labels/1",
          "rule_href": "/api/v1/rules/grocery-rule"
        }
      ]
    }
  ],
  "total": 1,
  "start": "2024-01-01T00:00:00Z",
  "end": "2024-01-31T23:59:59Z"
}
```

##### Get Specific Transaction
```bash
curl -X GET http://localhost:8080/api/v1/transactions/abc123def456
```

**Response:**
```json
{
  "id": 1,
  "date": "2024-01-15T00:00:00Z",
  "kind": "debit",
  "amount": 45.50,
  "content": "grocery store purchase",
  "hash": "abc123def456",
  "labels": [
    {
      "key": "category",
      "value": "food",
      "href": "/api/v1/labels/1",
      "rule_href": "/api/v1/rules/grocery-rule"
    },
    {
      "key": "type",
      "value": "essential",
      "href": "/api/v1/labels/2",
      "rule_href": "/api/v1/rules/grocery-rule"
    }
  ]
}
```

##### Create Transaction
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-15",
    "kind": "debit",
    "amount": 45.50,
    "content": "grocery store purchase",
    "account": 1001,
    "labels": {
      "category": "food",
      "type": "essential"
    }
  }'
```

**Response:**
```json
{
  "id": 1,
  "date": "2024-01-15T00:00:00Z",
  "kind": "debit",
  "amount": 45.50,
  "content": "grocery store purchase",
  "hash": "abc123def456",
  "labels": []
}
```

**Attempting to create duplicate transaction:**
```bash
# Same request as above
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-15",
    "kind": "debit",
    "amount": 45.50,
    "content": "grocery store purchase",
    "account": 1001
  }'
```

**Error Response (400 Bad Request):**
```json
{
  "error": "transaction 1 already exists"
}
```

##### Update Transaction
```bash
curl -X PUT http://localhost:8080/api/v1/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-15",
    "kind": "debit",
    "amount": 47.50,
    "content": "grocery store purchase updated",
    "account": 1001,
    "labels": {
      "category": "food",
      "type": "essential"
    }
  }'
```

**Response (200 OK - existing transaction updated):**
```json
{
  "id": 1,
  "date": "2024-01-15T00:00:00Z",
  "kind": "debit",
  "amount": 47.50,
  "content": "grocery store purchase updated",
  "hash": "def456ghi789",
  "labels": []
}
```

**Updating non-existent transaction (creates new):**
```bash
curl -X PUT http://localhost:8080/api/v1/transactions/999 \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-16",
    "kind": "credit",
    "amount": 100.00,
    "content": "new transaction via update",
    "account": 1001
  }'
```

**Response (201 Created - new transaction created):**
```json
{
  "id": 150,
  "date": "2024-01-16T00:00:00Z",
  "kind": "credit",
  "amount": 100.00,
  "content": "new transaction via update",
  "hash": "ghi789jkl012",
  "labels": []
}
```

##### Delete Transaction
```bash
curl -X DELETE http://localhost:8080/api/v1/transactions/1
```

**Response:**
```json
{
  "message": "Transaction deleted successfully"
}
```

#### Transaction Labels

##### Get Transaction Labels
```bash
curl -X GET http://localhost:8080/api/v1/transactions/1/labels
```

**Response:**
```json
{
  "total": 2,
  "labels": [
    {
      "href": "/api/v1/labels/5",
      "key": "category",
      "value": "food",
      "rules": [
        {
          "name": "grocery-rule",
          "href": "/api/v1/rules/grocery-rule"
        }
      ]
    },
    {
      "href": "/api/v1/labels/8",
      "key": "type",
      "value": "essential",
      "rules": []
    }
  ]
}
```

##### Add Label to Transaction
```bash
curl -X POST http://localhost:8080/api/v1/transactions/1/labels \
  -H "Content-Type: application/json" \
  -d '{
    "key": "category",
    "value": "groceries"
  }'
```

**Response:**
```json
{
  "href": "/api/v1/labels/12",
  "key": "category",
  "value": "groceries",
  "rules": []
}
```

##### Remove All Labels from Transaction
```bash
curl -X DELETE http://localhost:8080/api/v1/transactions/1/labels
```

**Response:**
```
HTTP/1.1 204 No Content
```

##### Remove Specific Label from Transaction
```bash
curl -X DELETE http://localhost:8080/api/v1/transactions/1/labels/12
```

**Response:**
```
HTTP/1.1 204 No Content
```

#### Rules

##### List All Rules
```bash
curl -X GET http://localhost:8080/api/v1/rules
```

**Response:**
```json
{
  "data": [
    {
      "name": "grocery-rule",
      "pattern": "grocery|supermarket|food",
      "labels": [
        {
          "key": "category",
          "value": "food",
          "href": "/api/v1/labels/1"
        },
        {
          "key": "type",
          "value": "essential",
          "href": "/api/v1/labels/2"
        }
      ]
    },
    {
      "name": "salary-rule",
      "pattern": "salary|wage|payroll",
      "labels": [
        {
          "key": "category",
          "value": "income",
          "href": "/api/v1/labels/3"
        },
        {
          "key": "source",
          "value": "work",
          "href": "/api/v1/labels/4"
        }
      ]
    }
  ]
}
```

##### Get Specific Rule
```bash
curl -X GET http://localhost:8080/api/v1/rules/grocery-rule
```

**Response:**
```json
{
  "name": "grocery-rule",
  "pattern": "grocery|supermarket|food",
  "labels": [
    {
      "key": "category",
      "value": "food",
      "href": "/api/v1/labels/1"
    },
    {
      "key": "type",
      "value": "essential",
      "href": "/api/v1/labels/2"
    }
  ]
}
```

##### Create Rule
```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "grocery-rule",
    "pattern": "grocery|supermarket|food",
    "labels": {
      "category": "food",
      "type": "essential"
    }
  }'
```

**Response:**
```json
{
  "name": "grocery-rule",
  "pattern": "grocery|supermarket|food",
  "labels": [
    {
      "key": "category",
      "value": "food",
      "href": "/api/v1/labels/1"
    },
    {
      "key": "type",
      "value": "essential",
      "href": "/api/v1/labels/2"
    }
  ]
}
```

##### Update Rule
```bash
curl -X PUT http://localhost:8080/api/v1/rules/grocery-rule \
  -H "Content-Type: application/json" \
  -d '{
    "name": "grocery-rule",
    "pattern": "grocery|supermarket|food|market",
    "labels": {
      "category": "food",
      "type": "essential",
      "subcategory": "shopping"
    }
  }'
```

##### Delete Rule
```bash
curl -X DELETE http://localhost:8080/api/v1/rules/grocery-rule
```

**Response:**
```json
{
  "message": "Rule deleted successfully"
}
```

#### Labels

##### List All Labels
```bash
curl -X GET http://localhost:8080/api/v1/labels
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "key": "category",
      "value": "food",
      "href": "/api/v1/labels/1",
      "rules": ["grocery-rule", "restaurant-rule"],
      "transaction_count": 25
    },
    {
      "id": 2,
      "key": "category",
      "value": "transport",
      "href": "/api/v1/labels/2",
      "rules": ["uber-rule", "gas-rule"],
      "transaction_count": 12
    }
  ]
}
```

##### Create Label
```bash
curl -X POST http://localhost:8080/api/v1/labels \
  -H "Content-Type: application/json" \
  -d '{
    "key": "category",
    "value": "entertainment"
  }'
```

**Response:**
```json
{
  "id": 3,
  "key": "category",
  "value": "entertainment",
  "href": "/api/v1/labels/3",
  "rules": [],
  "transaction_count": 0
}
```

##### Delete Label
```bash
curl -X DELETE http://localhost:8080/api/v1/labels/3
```

**Response:**
```json
{
  "message": "Label deleted successfully"
}
```

#### Analytics

##### Get Summary Statistics
```bash
curl -X GET http://localhost:8080/api/v1/summary
```

**Response:**
```json
{
  "total_transactions": 156,
  "total_amount": 4567.89,
  "from": "2024-01-01T00:00:00Z",
  "to": "2024-01-31T23:59:59Z",
  "by_label": [
    {
      "label": "category:food",
      "count": 25,
      "total_amount": 1234.56,
      "average_amount": 49.38
    },
    {
      "label": "category:transport",
      "count": 12,
      "total_amount": 567.89,
      "average_amount": 47.32
    }
  ],
  "by_type": {
    "debit": {
      "count": 140,
      "total_amount": 3456.78
    },
    "credit": {
      "count": 16,
      "total_amount": 5678.90
    }
  }
}
```

#### Error Responses

All endpoints return structured error responses:

```json
{
  "error": "Validation failed",
  "message": "Invalid transaction data",
  "details": [
    {
      "field": "amount",
      "message": "Amount must be greater than 0"
    },
    {
      "field": "date",
      "message": "Date is required"
    }
  ]
}
```

**Resource-specific Error Examples:**

Creating duplicate transaction:
```json
{
  "error": "transaction 123 already exists"
}
```

Updating non-existent rule:
```json
{
  "error": "rule example-rule not found"
}
```

Resource not found:
```json
{
  "error": "transaction with hash abc123 not found"
}
```

Common HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data or resource already exists
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (deprecated in favor of 400 for duplicates)
- `500 Internal Server Error` - Server error

## 🗂️ Project Structure

```
finante/
├── cmd/                    # Command-line interface
│   └── serve.go            # Server start command
├── internal/               # Private application code
│   ├── config/             # Configuration management
│   ├── datastore/          # Database layer
│   │   └── pg/             # PostgreSQL implementation
│   ├── entity/             # Domain entities
│   ├── handlers/           # HTTP handlers
│   │   └── v1/             # API v1 handlers
│   ├── server/             # HTTP server setup
│   └── services/           # Business logic layer
├── pkg/                    # Public libraries
│   ├── context/            # Context utilities
│   ├── logger/             # Logging setup
│   ├── migrations/         # Database migrations
│   ├── parser/             # Text parsing utilities
│   └── reader/             # File reading utilities
├── ui/                     # React frontend
└── main.go                 # Application entry point
```

## 🔄 Development Workflow

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/services/...
```

### Database Migrations
```bash
# Create new migration
goose -dir pkg/migrations/sql create migration_name sql

# Run migrations
make migrate

# Check migration status
goose -dir pkg/migrations/sql status
```

### Code Generation
Some code is auto-generated using optgen for configuration builders:
```bash
# Regenerate auto-generated code
go generate ./...
```

### Building
```bash
# Build binary
make build

# Build for different platforms
GOOS=linux GOARCH=amd64 make build
```

## 🧪 Testing

The project includes comprehensive testing:

- **Unit Tests**: Service layer and business logic testing
- **Integration Tests**: Database operations testing
- **API Tests**: HTTP endpoint testing

Run tests with:
```bash
make test
```

## 🐳 Docker Support

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o finante .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/finante .
COPY --from=builder /app/pkg/migrations/sql ./migrations/
CMD ["./finante", "serve"]
```

## 🔗 Relationship Management

The system automatically manages complex relationships between transactions, rules, and labels:

### Automatic Label Creation
- When creating transactions or rules, the system automatically creates labels that don't exist
- Existing labels are reused across multiple transactions and rules
- Label relationships are properly maintained during updates

### Rule-Label Associations
- Rules can be associated with multiple labels
- When a rule matches a transaction, all associated labels are applied
- Relationship integrity is maintained when rules or labels are updated

### Transaction-Label Relationships
- Transactions can have labels applied manually or through rule matching
- Each label association tracks whether it was applied by a rule or manually
- Updating transactions properly manages existing label relationships

### Update Behavior
- **Create Operations**: Fail if resource already exists (transactions, rules)
- **Update Operations**: 
  - Transactions: Creates new if doesn't exist (upsert behavior)
  - Rules: Fails if rule doesn't exist (strict update)
- **Relationship Updates**: Old relationships are removed and new ones are created atomically

## 📊 Excel Import Format

The system can import transactions from Excel files with improved error handling and relationship management:

| Date       | Description          | Debit  | Credit |
|------------|---------------------|--------|--------|
| 15/01/2024 | Grocery Store       | 45.50  |        |
| 16/01/2024 | Salary Payment      |        | 2500.00|

- **Date**: DD/MM/YYYY format
- **Description**: Transaction description (used for rule matching)
- **Debit**: Debit amount (expense)
- **Credit**: Credit amount (income)

### Import Behavior
- **Duplicate Detection**: Import automatically detects existing transactions by hash
- **Smart Updates**: Existing transactions are updated, new transactions are created
- **Rule Application**: Rules are applied to all transactions during import
- **Error Handling**: Individual transaction errors don't stop the entire import process
- **Progress Tracking**: Detailed import results with created/updated/error counts

## 🏷️ Label System

The system uses a structured key-value label system for categorizing transactions:

### Label Structure
- **Key**: The category type (e.g., "category", "type", "merchant")
- **Value**: The specific value within that category (e.g., "food", "essential", "grocery")

### Common Label Keys
- `category`: Primary classification (food, transport, income, etc.)
- `type`: Secondary classification (essential, luxury, recurring, etc.)
- `merchant`: Vendor or service provider
- `location`: Geographic location
- `project`: For business expense tracking

### Examples
```json
{
  "category": "food",
  "type": "essential",
  "merchant": "grocery_store"
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style
- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Add documentation comments for all public functions
- Write tests for new functionality

## 📄 License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check existing documentation
- Review the API examples above

## 🔮 Future Enhancements

- [ ] Mobile app support
- [ ] Advanced analytics and reporting
- [ ] Multi-currency support
- [ ] Budget planning and tracking
- [ ] Bank API integrations
- [ ] Enhanced rule engine with ML capabilities
- [ ] Real-time notifications
- [ ] Data export in various formats
- [x] ~~Label hierarchy and inheritance~~ **Completed**: Comprehensive relationship management
- [x] ~~Custom label validation rules~~ **Completed**: Improved service layer validation
- [ ] Bulk operations API for transactions and rules
- [ ] Transaction scheduling and recurring payments
- [ ] Advanced filtering and search capabilities
- [ ] Audit trail for all data changes
- [ ] API rate limiting and authentication
- [ ] Performance monitoring and metrics
