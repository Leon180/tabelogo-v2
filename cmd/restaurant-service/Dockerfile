# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go.mod files (needed for replace directives)
COPY go.mod go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY cmd/restaurant-service/go.mod cmd/restaurant-service/go.sum ./cmd/restaurant-service/

# Download dependencies (with GOWORK=off to ignore go.work)
WORKDIR /app/cmd/restaurant-service
RUN GOWORK=off go mod download

# Install swag CLI for Swagger generation (optional for future)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
WORKDIR /app
COPY . .

# Build the application (with GOWORK=off)
WORKDIR /app/cmd/restaurant-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOWORK=off \
    go build -ldflags="-w -s" -o /app/bin/restaurant-service .

# Install golang-migrate for database migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata netcat-openbsd wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/restaurant-service .

# Copy migrate binary from builder
COPY --from=builder /go/bin/migrate .

# Copy migrations
COPY --from=builder /app/migrations/restaurant ./migrations/restaurant

# Copy entrypoint script
COPY cmd/restaurant-service/entrypoint.sh .
RUN chmod +x entrypoint.sh

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose HTTP port
EXPOSE 18082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:18082/health || exit 1

# Run the application with entrypoint
ENTRYPOINT ["./entrypoint.sh"]
