# Multi-stage build for Finante application
# Stage 1: Build React frontend
FROM docker.io/node:18-alpine AS frontend-builder

WORKDIR /app/ui

# Copy package files for better caching
COPY ui/package*.json ./
RUN npm ci --only=production

# Copy source and build
COPY ui/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM docker.io/golang:1.23-alpine AS backend-builder

WORKDIR /app

# Copy go mod files for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o finante .

# Stage 3: Final runtime image
FROM alpine:latest

# Create non-root user for security
RUN addgroup -g 1001 -S finante && \
    adduser -S finante -u 1001 -G finante

WORKDIR /app

# Copy built application from backend builder
COPY --from=backend-builder /app/finante .

# Copy database migrations
COPY --from=backend-builder /app/internal/datastore/pg/migrations/sql ./migrations/

# Copy built frontend from frontend builder
COPY --from=frontend-builder /app/ui/dist ./ui/dist

# Create directory for uploads and logs
RUN mkdir -p /app/uploads /app/logs && \
    chown -R finante:finante /app

# Switch to non-root user
USER finante

# Expose port
EXPOSE 8080

# Default command
CMD ["./finante", "serve", "--server-port=8080", "--server-gin-mode=release"] 
