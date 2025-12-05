# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod files (needed for replace directives)
COPY go.mod go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY cmd/map-service/go.mod cmd/map-service/go.sum ./cmd/map-service/

# Download dependencies (with GOWORK=off to ignore go.work)
WORKDIR /app/cmd/map-service
RUN GOWORK=off go mod download

# Install swag CLI for Swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
WORKDIR /app
COPY . .

# Generate Swagger documentation (run from service directory to limit scope)
WORKDIR /app/cmd/map-service
RUN swag init \
    --generalInfo main.go \
    --dir .,../../internal/map/interfaces/http \
    --output ../../internal/map/docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/bin/map-service .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/map-service .

# Copy Swagger docs
COPY --from=builder /app/internal/map/docs ./docs

# Expose ports (HTTP and gRPC)
EXPOSE 8081 19083

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the application
CMD ["./map-service"]
