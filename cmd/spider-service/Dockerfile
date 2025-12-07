# Spider Service Dockerfile
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go.mod files (needed for replace directives)
COPY go.mod go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY cmd/spider-service/go.mod cmd/spider-service/go.sum ./cmd/spider-service/

# Download dependencies (with GOWORK=off to ignore go.work)
WORKDIR /app/cmd/spider-service
RUN GOWORK=off go mod download

# Copy source code
WORKDIR /app
COPY . .

# Build the application (with GOWORK=off)
WORKDIR /app/cmd/spider-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOWORK=off \
    go build -ldflags="-w -s" -o /app/bin/spider-service .

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/spider-service .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 18084 19084

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:18084/health || exit 1

CMD ["./spider-service"]
