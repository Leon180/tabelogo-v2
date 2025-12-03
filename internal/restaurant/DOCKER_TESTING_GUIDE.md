# Restaurant Service - Docker Testing Guide

**Date:** 2025-12-03
**Purpose:** Guide for testing Restaurant Service with Docker and Docker Compose

---

## Prerequisites

1. **Docker Desktop** must be running
2. **Docker Compose** installed (comes with Docker Desktop)
3. Ports available:
   - `18082` - Restaurant Service HTTP
   - `5433` - PostgreSQL (mapped to internal 5432)
   - `6379` - Redis

---

## Quick Start - Docker Compose Testing

### Option 1: Full Integration Test (Recommended)

This script tests the complete Restaurant Service with all dependencies:

```bash
# Navigate to project root
cd /Users/lileon/goproject/tabelogov2

# Make sure Docker is running, then execute:
./scripts/test-restaurant-docker.sh
```

**What this script does:**
1. Stops any existing containers
2. Builds the Restaurant Service Docker image
3. Starts PostgreSQL and Redis
4. Waits for databases to be ready
5. Starts Restaurant Service
6. Runs automated migrations
7. Tests health endpoint
8. Tests restaurant creation API
9. Tests restaurant retrieval by ID
10. Tests restaurant retrieval by external ID
11. Tests restaurant search

**Expected Output:**
```
============================================
All tests passed! ✓
============================================

Service is running at:
  - Health Check: http://localhost:18082/health
  - Restaurants API: http://localhost:18082/api/v1/restaurants
```

---

### Option 2: Manual Docker Compose

For manual testing and development:

```bash
# Navigate to docker-compose directory
cd deployments/docker-compose

# Start all services
docker-compose up -d

# Or start only restaurant service and dependencies
docker-compose up -d postgres-restaurant redis restaurant-service

# View logs
docker-compose logs -f restaurant-service

# Check status
docker-compose ps

# Stop services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

---

## Docker Image Build Only

To just build and test the Docker image without docker-compose:

```bash
# Navigate to project root
cd /Users/lileon/goproject/tabelogov2

# Run build script
./scripts/build-restaurant-service.sh
```

This verifies:
- Dockerfile syntax
- Go build process
- Multi-stage build
- Binary size optimization

**Manual build command:**
```bash
docker build -f cmd/restaurant-service/Dockerfile -t tabelogo-restaurant-service:test .
```

---

## Testing Checklist

Before committing Restaurant Service changes, verify:

### 1. Unit Tests (98% Coverage)
```bash
go test ./internal/restaurant/domain/model/... ./internal/restaurant/application/... -v
```
**Expected:** All 92 tests pass

### 2. Docker Build
```bash
./scripts/build-restaurant-service.sh
```
**Expected:** Image builds successfully

### 3. Docker Compose Integration
```bash
./scripts/test-restaurant-docker.sh
```
**Expected:** All 13 test steps pass

### 4. Manual API Testing

Once service is running via docker-compose:

**Health Check:**
```bash
curl http://localhost:18082/health
```

**Create Restaurant:**
```bash
curl -X POST http://localhost:18082/api/v1/restaurants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sushi Dai",
    "source": "google",
    "external_id": "ChIJTest123",
    "address": "Tsukiji Fish Market, Tokyo",
    "latitude": 35.6654,
    "longitude": 139.7707,
    "rating": 4.5,
    "cuisine_type": "Japanese"
  }'
```

**Get Restaurant by External ID:**
```bash
curl http://localhost:18082/api/v1/restaurants/external/google/ChIJTest123
```

**List Restaurants:**
```bash
curl "http://localhost:18082/api/v1/restaurants?limit=10&offset=0"
```

---

## Docker Architecture

### Multi-Stage Build

The Dockerfile uses a two-stage build for optimization:

**Stage 1: Builder (golang:1.24-alpine)**
- Installs build dependencies
- Downloads Go modules
- Builds restaurant-service binary
- Installs golang-migrate for migrations

**Stage 2: Runtime (alpine:3.19)**
- Minimal runtime image (~50MB vs ~500MB)
- Only includes:
  - restaurant-service binary
  - migrate binary
  - Migration SQL files
  - Runtime dependencies (ca-certificates, netcat, wget)
- Non-root user (appuser:1000)

### Container Communication

```
┌─────────────────────────────────────────────────────┐
│ Docker Network: tabelogo-network                     │
│                                                       │
│  ┌──────────────────┐      ┌──────────────────┐    │
│  │  postgres-       │:5432 │  redis           │    │
│  │  restaurant      │      │                  │    │
│  └──────────────────┘      └──────────────────┘    │
│           │                         │               │
│           └─────────┬───────────────┘               │
│                     │                               │
│           ┌─────────▼────────────┐                 │
│           │  restaurant-service  │                 │
│           │  :18082              │                 │
│           └──────────────────────┘                 │
└─────────────────────────────────────────────────────┘
                      │
              Port Mapping
                      │
         ┌────────────▼────────────┐
         │  localhost:18082        │ (HTTP API)
         │  localhost:5433         │ (PostgreSQL)
         │  localhost:6379         │ (Redis)
         └─────────────────────────┘
```

### Environment Variables

**Docker Compose Sets:**
```yaml
ENVIRONMENT: development
LOG_LEVEL: info
SERVER_PORT: 18082
DB_HOST: postgres-restaurant  # Container hostname
DB_PORT: 5432                 # Internal port
DB_NAME: restaurant_db
REDIS_HOST: redis
REDIS_PORT: 6379
REDIS_DB: 1
```

**Local Development (.env):**
```bash
DATABASE_HOST=localhost
DATABASE_PORT=5433  # External mapped port
```

---

## Startup Sequence

The `entrypoint.sh` script ensures proper initialization:

1. **Wait for PostgreSQL** - Uses `nc` (netcat) to check connectivity
2. **Wait for Redis** - Optional, checks if REDIS_HOST is set
3. **Run Migrations** - Executes `golang-migrate` with PostgreSQL connection
4. **Start Service** - Launches restaurant-service binary

**Migration Connection String:**
```
postgresql://postgres:postgres@postgres-restaurant:5432/restaurant_db?sslmode=disable
```

---

## Troubleshooting

### Docker Daemon Not Running
```bash
# Error: Cannot connect to Docker daemon
# Solution: Start Docker Desktop
```

### Port Already in Use
```bash
# Error: port is already allocated
docker ps -a  # Check running containers
docker-compose down  # Stop services
# Or change port in docker-compose.yml
```

### Migration Failed
```bash
# View migration logs
docker-compose logs restaurant-service

# Check PostgreSQL is accessible
docker-compose exec postgres-restaurant psql -U postgres -d restaurant_db -c '\dt'

# Manually run migrations
docker-compose exec restaurant-service /app/migrate \
  -path /app/migrations/restaurant \
  -database "postgresql://postgres:postgres@postgres-restaurant:5432/restaurant_db?sslmode=disable" \
  up
```

### Service Won't Start
```bash
# Check service logs
docker-compose logs -f restaurant-service

# Check health status
docker-compose ps restaurant-service

# Enter container for debugging
docker-compose exec restaurant-service sh
```

### Database Connection Issues
```bash
# Test PostgreSQL connectivity from container
docker-compose exec restaurant-service nc -zv postgres-restaurant 5432

# Test from host
docker-compose exec postgres-restaurant pg_isready -U postgres
```

---

## Files Modified for Docker Support

### Created:
1. `cmd/restaurant-service/Dockerfile` - Multi-stage build definition
2. `cmd/restaurant-service/entrypoint.sh` - Container startup script
3. `cmd/restaurant-service/.env` - Local environment config
4. `cmd/restaurant-service/.env.example` - Environment template
5. `cmd/restaurant-service/.gitignore` - Ignore patterns
6. `scripts/test-restaurant-docker.sh` - Automated integration test
7. `scripts/build-restaurant-service.sh` - Build verification script

### Updated:
1. `deployments/docker-compose/docker-compose.yml` - Added restaurant-service
2. `cmd/restaurant-service/go.mod` - Service module definition
3. `cmd/restaurant-service/go.sum` - Dependency checksums

---

## Performance Notes

**Build Time:** ~2-3 minutes (first build)
- Subsequent builds: ~30 seconds (Docker layer caching)

**Startup Time:** ~15-20 seconds
- PostgreSQL ready: ~5 seconds
- Migrations: ~3 seconds
- Service start: ~2 seconds

**Image Size:**
- Builder stage: ~800MB (discarded)
- Final runtime: ~50MB

**Memory Usage:**
- Restaurant Service: ~50MB
- PostgreSQL: ~50MB
- Redis: ~10MB
- **Total:** ~110MB

---

## Integration with Other Services

### Map Service Integration
Map Service should connect to Restaurant Service for deduplication:

```bash
# Check if restaurant exists before creating
curl http://localhost:18082/api/v1/restaurants/external/google/{place_id}

# If 404, create new restaurant
curl -X POST http://localhost:18082/api/v1/restaurants ...
```

### Spider Service Integration
Spider Service should connect for Tabelog restaurants:

```bash
# Check if restaurant exists
curl http://localhost:18082/api/v1/restaurants/external/tabelog/{url_encoded}
```

### Frontend Integration
Frontend connects to Restaurant Service for:
- Search and discovery
- Favorites management
- Restaurant details

---

## Next Steps

After successful Docker testing:

1. ✅ Commit Restaurant Service changes
2. ⏳ Implement API Gateway (Phase 2 completion)
3. ⏳ Add Swagger/OpenAPI documentation
4. ⏳ Add monitoring and metrics (Prometheus)
5. ⏳ Add distributed tracing (Jaeger)

---

**Documentation Status:** ✅ Complete
**Last Updated:** 2025-12-03
