# Multi-stage build for Finante application
# Stage 1: Build React frontend
FROM docker.io/node:22-alpine AS frontend-builder

ARG GIT_SHA
ENV GIT_SHA=${GIT_SHA}

WORKDIR /app/ui

# Copy package files for better caching
COPY ui/package*.json ./
RUN npm ci

# Copy source and build
COPY ui/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM docker.io/golang:1.25-alpine AS backend-builder

ARG GIT_SHA

RUN apk add --no-cache build-base

WORKDIR /app

# Copy go mod files for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN GOOS=linux go build -ldflags="-X main.sha=${GIT_SHA}" -o finante .

# Stage 3: Final runtime image
FROM alpine:3.21

# Create non-root user for security
RUN addgroup -g 1001 -S finante && \
    adduser -S finante -u 1001 -G finante

WORKDIR /app

# Copy built application from backend builder
COPY --from=backend-builder /app/finante .

# Copy database migrations
COPY --from=backend-builder /app/internal/store/migrations/sql ./migrations/

# Copy built frontend from frontend builder
COPY --from=frontend-builder /app/ui/dist ./ui/dist

# Create directory for uploads and logs
RUN mkdir -p /app/uploads /app/logs && \
    chown -R finante:finante /app

# Switch to non-root user
USER finante

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD wget -qO- http://localhost:8080/healthz || exit 1

# Default command
CMD ["./finante", "serve", "--server-port=8080", "--server-gin-mode=release", "--server-mode=prod", "--statics-folder=./ui/dist"]
