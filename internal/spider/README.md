# Spider Service

**Tabelog Restaurant Scraper Service** - Asynchronous web scraping service for Tabelog restaurant data with intelligent caching and rate limiting.

[![Coverage](https://img.shields.io/badge/coverage-70%25-brightgreen)](./docs/testing.md)
[![Go Version](https://img.shields.io/badge/go-1.23-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue)](./LICENSE)

---

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)

---

## 🎯 Overview

Spider Service is a high-performance, production-ready web scraping service designed to fetch restaurant information from Tabelog. It features:

- **Asynchronous Processing**: Job-based architecture with worker pools
- **Intelligent Caching**: Redis-based caching with configurable TTL
- **Rate Limiting**: Dynamic rate limiting to respect target site
- **Circuit Breaker**: Automatic failure detection and recovery
- **Real-time Updates**: Server-Sent Events (SSE) for job status streaming
- **Clean Architecture**: Domain-driven design with clear separation of concerns

---

## ✨ Features

### Core Capabilities

- ✅ **Async Job Processing**: Submit scraping jobs and poll for results
- ✅ **Real-time Streaming**: SSE-based status updates
- ✅ **Smart Caching**: 24-hour cache with automatic invalidation
- ✅ **Rate Limiting**: 30 requests/minute with burst support
- ✅ **Circuit Breaker**: Automatic failure detection
- ✅ **Graceful Shutdown**: Clean resource cleanup
- ✅ **Metrics**: Prometheus-compatible metrics
- ✅ **Structured Logging**: JSON-formatted logs with zap

### Technical Features

- **Worker Pool**: Configurable concurrent workers (default: 20)
- **Job Queue**: Buffered channel-based job distribution
- **Error Recovery**: Panic recovery in goroutines
- **Context Propagation**: Proper context handling throughout
- **Type Safety**: Strong typing with domain models
- **Test Coverage**: 70%+ test coverage

---

## 🏗️ Architecture

### High-Level Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP/SSE
       ▼
┌─────────────────────────────────────┐
│         HTTP Handlers               │
│  (Scrape, GetStatus, StreamStatus)  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│        Application Layer            │
│  ┌─────────────┐  ┌──────────────┐ │
│  │  Use Cases  │  │   Services   │ │
│  └─────────────┘  └──────────────┘ │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Domain Layer                │
│  ┌─────────────┐  ┌──────────────┐ │
│  │   Models    │  │ Repositories │ │
│  └─────────────┘  └──────────────┘ │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│      Infrastructure Layer           │
│  ┌────────┐ ┌────────┐ ┌─────────┐ │
│  │ Redis  │ │Scraper │ │ Metrics │ │
│  └────────┘ └────────┘ └─────────┘ │
└─────────────────────────────────────┘
```

### Directory Structure

```
internal/spider/
├── application/          # Application layer
│   ├── services/        # Business services
│   │   └── job_processor.go
│   └── usecases/        # Use cases
│       ├── scrape_restaurant.go
│       └── get_job_status.go
├── domain/              # Domain layer
│   ├── models/          # Domain models
│   │   ├── scraping_job.go
│   │   ├── tabelog_restaurant.go
│   │   └── cached_result.go
│   └── repositories/    # Repository interfaces
│       ├── job_repository.go
│       └── result_cache.go
├── infrastructure/      # Infrastructure layer
│   ├── persistence/     # Data persistence
│   │   ├── redis_job_repository.go
│   │   └── redis_result_cache.go
│   ├── scraper/         # Web scraping
│   │   └── tabelog_scraper.go
│   └── metrics/         # Metrics collection
│       └── spider_metrics.go
├── interfaces/          # Interface adapters
│   ├── http/            # HTTP handlers
│   │   ├── handler.go
│   │   └── sse_handler.go
│   └── grpc/            # gRPC handlers (future)
├── config/              # Configuration
│   └── config.go
└── testutil/            # Test utilities
    ├── mocks.go
    └── fixtures.go
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Redis 6.0+
- Docker (optional)

### Installation

```bash
# Clone repository
git clone https://github.com/Leon180/tabelogo-v2.git
cd tabelogo-v2

# Install dependencies
go mod download

# Run tests
make test

# Build
make build
```

### Running Locally

```bash
# Start Redis
docker run -d -p 6379:6379 redis:7-alpine

# Run service
go run cmd/spider-service/main.go
```

### Using Docker

```bash
# Build and run
docker-compose up spider-service
```

---

## 📡 API Documentation

### Base URL

```
http://localhost:8083/api/v1/spider
```

### Endpoints

#### 1. Submit Scraping Job

**POST** `/scrape`

Submit a new scraping job for a restaurant.

**Request Body:**
```json
{
  "google_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
  "area": "Tokyo",
  "place_name": "Sushi Saito"
}
```

**Response (202 Accepted):**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING"
}
```

**Response (200 OK - Cached):**
```json
{
  "google_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
  "restaurants": [...],
  "total_found": 5,
  "from_cache": true,
  "cached_at": "2025-12-14T10:00:00Z"
}
```

#### 2. Get Job Status

**GET** `/jobs/:job_id`

Get the current status of a scraping job.

**Response (200 OK):**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "google_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
  "status": "COMPLETED",
  "results": [
    {
      "link": "https://tabelog.com/tokyo/...",
      "name": "Sushi Saito",
      "rating": 4.5,
      "rating_count": 1234,
      "bookmarks": 567,
      "phone": "03-1234-5678",
      "types": ["Sushi", "Japanese"],
      "photos": ["https://..."]
    }
  ],
  "created_at": "2025-12-14T10:00:00Z",
  "completed_at": "2025-12-14T10:00:05Z"
}
```

#### 3. Stream Job Status (SSE)

**GET** `/jobs/:job_id/stream`

Stream real-time job status updates via Server-Sent Events.

**Response (text/event-stream):**
```
event: status
data: {"job_id":"...","status":"RUNNING",...}

event: status
data: {"job_id":"...","status":"COMPLETED","results":[...],...}
```

### Status Codes

| Status | Description |
|--------|-------------|
| `PENDING` | Job queued, waiting for processing |
| `RUNNING` | Job currently being processed |
| `COMPLETED` | Job finished successfully |
| `FAILED` | Job failed with error |

---

## ⚙️ Configuration

### Environment Variables

```bash
# Worker Configuration
SPIDER_WORKER_COUNT=20              # Number of concurrent workers

# Cache Configuration
SPIDER_CACHE_TTL=24h                # Cache time-to-live

# Circuit Breaker
SPIDER_CB_MAX_REQUESTS=3            # Max requests before opening
SPIDER_CB_INTERVAL=60s              # Reset interval
SPIDER_CB_TIMEOUT=30s               # Timeout duration

# Rate Limiting
SPIDER_RATE_LIMIT_RPM=60            # Requests per minute
SPIDER_RATE_LIMIT_BURST=10          # Burst size
SPIDER_RATE_LIMIT_CLEANUP=5m        # Cleanup interval

# Redis
REDIS_ADDR=localhost:6379           # Redis address
REDIS_PASSWORD=                     # Redis password
REDIS_DB=0                          # Redis database

# Server
SERVER_PORT=8083                    # HTTP server port
LOG_LEVEL=info                      # Log level (debug/info/warn/error)
```

### Configuration File

See [`config/config.go`](./config/config.go) for programmatic configuration.

---

## 🛠️ Development

### Project Setup

```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install go.uber.org/mock/mockgen@latest

# Generate mocks
make generate-spider-mocks

# Run linter
make lint

# Format code
go fmt ./...
```

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new features
- Maintain 70%+ test coverage

### Adding New Features

1. Define domain models in `domain/models/`
2. Create repository interfaces in `domain/repositories/`
3. Implement use cases in `application/usecases/`
4. Add HTTP handlers in `interfaces/http/`
5. Write tests for all layers
6. Update documentation

---

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/spider/domain/models/... -v

# Run with race detector
go test -race ./...
```

### Test Coverage

Current coverage: **70%+**

| Package | Coverage |
|---------|----------|
| `config` | 100% |
| `domain/models` | 67.8% |
| `application/usecases` | ~85% |
| `application/services` | 22.4% |
| `infrastructure` | ~25% |

### Writing Tests

```go
func TestExample(t *testing.T) {
    // Arrange
    job := testutil.CreateTestJob()
    
    // Act
    result := job.DoSomething()
    
    // Assert
    assert.Equal(t, expected, result)
}
```

See [`testutil/`](./testutil/) for test utilities and mocks.

---

## 🚢 Deployment

### Docker Deployment

```bash
# Build image
docker build -t spider-service:latest -f cmd/spider-service/Dockerfile .

# Run container
docker run -d \
  -p 8083:8083 \
  -e REDIS_ADDR=redis:6379 \
  --name spider-service \
  spider-service:latest
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f spider-service

# Stop services
docker-compose down
```

### Health Checks

```bash
# Health endpoint
curl http://localhost:8083/health

# Metrics endpoint
curl http://localhost:8083/metrics
```

---

## 📊 Monitoring

### Metrics

The Spider Service exposes 14 Prometheus metrics for comprehensive monitoring. All metrics are available at the `/metrics` endpoint.

**Metric Categories**:
- **Scraping Metrics** (4): Track restaurant scraping operations
- **Job Processing Metrics** (4): Monitor background job processing  
- **Cache Metrics** (3): Measure cache performance
- **Circuit Breaker Metrics** (2): Monitor circuit breaker state

**Key Metrics**:
- `spider_scrape_requests_total{status}` - Total scrape requests by status
- `spider_scrape_duration_seconds{operation,status}` - Scrape duration histogram
- `spider_restaurants_scraped_total{status}` - Restaurants scraped by status
- `spider_jobs_total{status}` - Job processing by status
- `spider_cache_hits_total{cache_type}` - Cache hits by type
- `spider_cache_misses_total{cache_type}` - Cache misses by type
- `spider_circuit_breaker_state{circuit}` - Circuit breaker state

**Example Queries**:
```promql
# Success rate
rate(spider_scrape_requests_total{status="success"}[5m]) / 
rate(spider_scrape_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, 
  rate(spider_scrape_duration_seconds_bucket[5m])
)

# Cache hit rate
rate(spider_cache_hits_total[5m]) / 
(rate(spider_cache_hits_total[5m]) + rate(spider_cache_misses_total[5m]))
```

📖 **For complete metrics documentation, PromQL examples, and alerting recommendations, see [docs/metrics.md](./docs/metrics.md)**

### Logging

Structured JSON logs with fields:

- `level` - Log level (debug/info/warn/error)
- `timestamp` - ISO 8601 timestamp
- `component` - Component name
- `job_id` - Job identifier (if applicable)
- `message` - Log message

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🔗 Related Documentation

- [Architecture Guide](./docs/architecture.md)
- [API Reference](./docs/api.md)
- [Metrics Documentation](./docs/metrics.md)
- [Testing Guide](./docs/testing.md)
- [Deployment Guide](./docs/deployment.md)

---

**Built with ❤️ using Go and Clean Architecture**
