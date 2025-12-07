# Spider Service Dockerfile
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go workspace files
COPY go.work go.work.sum ./

# Copy pkg module
COPY pkg/ ./pkg/

# Copy root module files
COPY go.mod go.sum ./

# Copy internal packages needed by spider service
COPY internal/spider/ ./internal/spider/

# Copy spider service
COPY cmd/spider-service/ ./cmd/spider-service/

# Build the service
WORKDIR /app/cmd/spider-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o spider-service .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/cmd/spider-service/spider-service .

# Expose port
EXPOSE 8084

CMD ["./spider-service"]
