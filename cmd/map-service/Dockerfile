# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY cmd/map-service/go.mod cmd/map-service/go.sum ./cmd/map-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Download dependencies
WORKDIR /app/cmd/map-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY cmd/map-service/ ./cmd/map-service/
COPY internal/map/ ./internal/map/
COPY pkg/ ./pkg/

# Build the application
WORKDIR /app/cmd/map-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o map-service .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/cmd/map-service/map-service .

# Expose ports
EXPOSE 8081

# Run the application
CMD ["./map-service"]
