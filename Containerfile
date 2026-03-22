# Stage 1: Build React frontend
FROM docker.io/node:22 AS frontend-builder

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
FROM registry.access.redhat.com/ubi9/go-toolset AS backend-builder

ARG GIT_SHA

# Copy go mod files for better caching
COPY go.mod go.sum ./
RUN go mod download

USER 0

# Copy source code
COPY . .

# Build the application
RUN GOOS=linux go build -ldflags="-X main.sha=${GIT_SHA}" -o /tmp/ducket .

# Stage 3: Final runtime image
FROM registry.access.redhat.com/ubi9/ubi-minimal

RUN microdnf install -y ca-certificates tzdata && \
    microdnf clean all

WORKDIR /app

# Copy built application from backend builder
COPY --from=backend-builder /tmp/ducket .

# Copy built frontend from frontend builder
COPY --from=frontend-builder /app/ui/dist ./ui/dist

# Create data directory
RUN mkdir -p /app/data && \
    chown -R 1001:0 /app

USER 1001

# Expose port
EXPOSE 8080

# Default command
CMD ["./ducket", "run", "--db-uri=/app/data/ducket.db", "--server-port=8080", "--server-gin-mode=release", "--server-mode=prod", "--statics-folder=./ui/dist"]
