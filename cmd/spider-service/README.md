# Spider Service

Web scraping microservice for Tabelog restaurant data.

## Overview

Spider Service scrapes Tabelog.com to collect restaurant information including ratings, reviews, photos, and other metadata. It provides a RESTful API for initiating scraping jobs and retrieving results.

## Features

- 🕷️ **Tabelog Scraping**: Link, content, and photo spiders
- ⚡ **Async Processing**: Non-blocking job execution
- 🔄 **Rate Limiting**: Respectful scraping (2 req/s, 500ms delay)
- 🎯 **Job Management**: Track scraping job status
- 🐳 **Docker Ready**: Containerized deployment

## Quick Start

### Local Development

```bash
# Install dependencies
go mod download

# Run service
go run main.go

# Or build and run
go build -o spider-service
./spider-service
```

### Docker

```bash
# Build image
docker build -t spider-service .

# Run container
docker run -p 8084:8084 spider-service
```

### Docker Compose

```bash
cd ../../deployments/docker-compose
docker-compose up spider-service
```

## API Endpoints

### Start Scraping Job

```bash
POST /api/v1/spider/scrape
Content-Type: application/json

{
  "google_id": "ChIJN5Nz71W3j4ARhx5bwpTQEGg",
  "area": "tokyo",
  "place_name": "Afuri Ramen"
}

# Response
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING"
}
```

### Get Job Status

```bash
GET /api/v1/spider/jobs/{job_id}

# Response
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "google_id": "ChIJN5Nz71W3j4ARhx5bwpTQEGg",
  "status": "COMPLETED",
  "results": [
    {
      "link": "https://tabelog.com/tokyo/...",
      "name": "Afuri Ramen Harajuku",
      "rating": 3.58,
      "rating_count": 1234,
      "bookmarks": 567,
      "phone": "03-1234-5678",
      "types": ["ラーメン", "つけ麺"],
      "photos": ["https://..."]
    }
  ],
  "created_at": "2025-12-07T16:00:00Z",
  "completed_at": "2025-12-07T16:00:05Z"
}
```

### Health Check

```bash
GET /health

# Response
{
  "status": "healthy",
  "service": "spider-service"
}
```

## Configuration

Environment variables:

```bash
# Server
SERVER_PORT=8084

# Logging
ENVIRONMENT=development
LOG_LEVEL=info

# Database (optional for MVP)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=spider_db
DB_USER=postgres
DB_PASSWORD=postgres

# Redis (for future job queue)
REDIS_HOST=localhost
REDIS_PORT=6379
```

## Architecture

```
internal/spider/
├── domain/          # Business logic
│   ├── models/      # Entities and value objects
│   └── repositories/
├── application/     # Use cases
│   └── usecases/
├── infrastructure/  # External dependencies
│   ├── scraper/     # Colly-based scraper
│   └── persistence/
└── interfaces/      # API layer
    └── http/
```

## Development

### Run Tests

```bash
go test ./...
```

### Build

```bash
go build -o spider-service
```

### Lint

```bash
golangci-lint run
```

## Scraping Details

### Rate Limiting

- **Requests per second**: 2
- **Delay between requests**: 500ms
- **Max concurrent**: 4
- **Timeout**: 10s
- **Retries**: 3

### User-Agent Rotation

Random User-Agent headers to avoid detection.

### Politeness

- Respects robots.txt (future)
- Random delays
- Exponential backoff on errors

## Next Steps

- [ ] PostgreSQL persistence
- [ ] Redis job queue
- [ ] Prometheus metrics
- [ ] gRPC API
- [ ] Unit tests

## License

MIT
