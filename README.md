# Finante - Personal Financial Management System

A comprehensive personal financial management system built with Go and React, designed to help you track, categorize, and analyze your financial transactions through automated rule-based tagging.

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
- **Service Layer**: Business logic for transactions, rules, and tags management
- **Repository Layer**: Database operations with PostgreSQL using pgx driver
- **Rule Engine**: Pattern-based transaction categorization system
- **File Parser**: Excel file import functionality for bulk transaction loading

## 🚀 Features

### Core Functionality
- **Transaction Management**: Create, read, update, and delete financial transactions
- **Rule-Based Categorization**: Automatic transaction tagging based on configurable rules
- **Tag Management**: Organize transactions with custom tags
- **Excel Import**: Bulk import transactions from Excel spreadsheets
- **Financial Analytics**: Transaction statistics and summaries

### Technical Features
- **RESTful API**: Well-structured API endpoints with JSON responses
- **Database Migrations**: Automated schema management with goose
- **Graceful Shutdown**: Proper server lifecycle management
- **Structured Logging**: Comprehensive logging with zap
- **Configuration Management**: Environment-based configuration
- **Connection Pooling**: Optimized database connection management

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
- `POST /transactions` - Create new transaction
- `PUT /transactions/:id` - Update transaction
- `DELETE /transactions/:id` - Delete transaction

#### Rules
- `GET /rules` - List all rules
- `GET /rules/:name` - Get specific rule
- `POST /rules` - Create new rule
- `PUT /rules/:name` - Update rule
- `DELETE /rules/:name` - Delete rule

#### Tags
- `GET /tags` - List all tags with transaction counts
- `POST /tags` - Create new tag
- `DELETE /tags/:value` - Delete tag

#### Analytics
- `GET /summary` - Get transaction statistics and summaries

### Request/Response Examples

#### Transactions

##### List All Transactions
```bash
# Get all transactions
curl -X GET http://localhost:8080/api/v1/transactions

# Get transactions with filters
curl -X GET "http://localhost:8080/api/v1/transactions?start=2024-01-01&end=2024-01-31&tags=food&limit=10&offset=0"
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
      "tags": {
        "food": "grocery-rule",
      }
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
  "tags": {
    "food": "grocery-rule",
    "essentials": "grocery-rule"
  }
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
    "content": "grocery store purchase"
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
  "tags": {}
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
    "content": "grocery store purchase updated"
  }'
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
      "tags": ["food", "essentials"],
      "created_at": "2024-01-01T10:00:00Z"
    },
    {
      "name": "salary-rule",
      "pattern": "salary|wage|payroll",
      "tags": ["income", "work"],
      "created_at": "2024-01-01T10:30:00Z"
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
  "tags": ["food", "essentials"],
  "created_at": "2024-01-01T10:00:00Z"
}
```

##### Create Rule
```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "grocery-rule",
    "pattern": "grocery|supermarket|food",
    "tags": ["food", "essentials"]
  }'
```

**Response:**
```json
{
  "name": "grocery-rule",
  "pattern": "grocery|supermarket|food",
  "tags": ["food", "essentials"],
  "created_at": "2024-01-15T14:30:00Z"
}
```

##### Update Rule
```bash
curl -X PUT http://localhost:8080/api/v1/rules/grocery-rule \
  -H "Content-Type: application/json" \
  -d '{
    "name": "grocery-rule",
    "pattern": "grocery|supermarket|food|market",
    "tags": ["food", "essentials", "shopping"]
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

#### Tags

##### List All Tags
```bash
curl -X GET http://localhost:8080/api/v1/tags
```

**Response:**
```json
{
  "data": [
    {
      "value": "food",
      "rules": ["grocery-rule", "restaurant-rule"],
      "transaction_count": 25,
      "created_at": "2024-01-01T10:00:00Z"
    },
    {
      "value": "transport",
      "rules": ["uber-rule", "gas-rule"],
      "transaction_count": 12,
      "created_at": "2024-01-01T11:00:00Z"
    }
  ]
}
```

##### Create Tag
```bash
curl -X POST http://localhost:8080/api/v1/tags \
  -H "Content-Type: application/json" \
  -d '{
    "value": "entertainment"
  }'
```

**Response:**
```json
{
  "value": "entertainment",
  "rules": [],
  "transaction_count": 0,
  "created_at": "2024-01-15T14:30:00Z"
}
```

##### Delete Tag
```bash
curl -X DELETE http://localhost:8080/api/v1/tags/entertainment
```

**Response:**
```json
{
  "message": "Tag deleted successfully"
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
  "by_tag": [
    {
      "tag": "food",
      "count": 25,
      "total_amount": 1234.56,
      "average_amount": 49.38
    },
    {
      "tag": "transport",
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

Common HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
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

## 📊 Excel Import Format

The system can import transactions from Excel files with the following format:

| Date       | Description          | Debit  | Credit |
|------------|---------------------|--------|--------|
| 15/01/2024 | Grocery Store       | 45.50  |        |
| 16/01/2024 | Salary Payment      |        | 2500.00|

- **Date**: DD/MM/YYYY format
- **Description**: Transaction description (used for rule matching)
- **Debit**: Debit amount (expense)
- **Credit**: Credit amount (income)

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
