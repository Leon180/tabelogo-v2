# Spider Service Dockerfile
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Copy API proto generated files
COPY api/ ./api/

# Copy internal packages needed by spider service
COPY internal/spider/ ./internal/spider/

# Copy spider service
COPY cmd/spider-service/ ./cmd/spider-service/

# Build the service (with GOWORK=off to ignore go.work file)
WORKDIR /app/cmd/spider-service
RUN GOWORK=off go mod download
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -ldflags="-w -s" -o /app/bin/spider-service .

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

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
