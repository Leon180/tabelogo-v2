# Restaurant Service - Deployment Guide

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose installed
- Go 1.24+ (for local development)
- PostgreSQL 15+ (port 5433)
- Redis 7+ (DB 1)

---

## 📦 Deployment Options

### Option 1: Docker Compose (Recommended)

#### 1.1 Start All Services
```bash
cd deployments/docker-compose
docker-compose up -d restaurant-service
```

This will:
- ✅ Start PostgreSQL (restaurant_db on port 5433)
- ✅ Start Redis (DB 1)
- ✅ Run database migrations automatically
- ✅ Start Restaurant Service on port **18082**

#### 1.2 Check Service Health
```bash
curl http://localhost:18082/health
# Expected: {"status":"healthy"}
```

#### 1.3 View Logs
```bash
docker-compose logs -f restaurant-service
```

#### 1.4 Stop Service
```bash
docker-compose stop restaurant-service
```

---

### Option 2: Local Binary

#### 2.1 Build Binary
```bash
go build -o bin/restaurant-service ./cmd/restaurant-service/
```

#### 2.2 Set Environment Variables
```bash
export SERVER_PORT=18082
export DATABASE_HOST=localhost
export DATABASE_PORT=5433
export DATABASE_NAME=restaurant_db
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_DB=1
```

Or use `.env.restaurant` file:
```bash
source .env.restaurant
```

#### 2.3 Run Migrations
```bash
migrate -path migrations/restaurant \
  -database "postgresql://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" \
  up
```

#### 2.4 Start Service
```bash
./bin/restaurant-service
```

---

## 🧪 Testing

### Integration Tests
```bash
# Make sure service is running on port 18082
./scripts/test-restaurant-service.sh
```

**Test Coverage**:
- ✅ Health check
- ✅ Create restaurant from Google Maps
- ✅ Prevent duplicate external IDs
- ✅ Create same restaurant from Tabelog (different source)
- ✅ Get restaurant by ID
- ✅ Search restaurants
- ✅ Add to favorites
- ✅ Get user favorites
- ✅ Prevent duplicate favorites

### Manual API Testing

#### Create Restaurant (Google Maps)
```bash
curl -X POST http://localhost:18082/api/v1/restaurants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sushi Dai",
    "source": "google",
    "external_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
    "address": "Tokyo, Chuo-ku, Tsukiji",
    "latitude": 35.6654,
    "longitude": 139.7707,
    "rating": 4.5,
    "price_range": "$$",
    "cuisine_type": "Sushi",
    "phone": "03-3547-6797"
  }'
```

#### Search Restaurants
```bash
curl "http://localhost:18082/api/v1/restaurants/search?q=sushi&limit=10"
```

#### Add to Favorites
```bash
curl -X POST http://localhost:18082/api/v1/favorites \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "restaurant_id": "<restaurant-id-from-previous-response>"
  }'
```

---

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `18082` | HTTP server port |
| `ENVIRONMENT` | `development` | Environment (development/production) |
| `LOG_LEVEL` | `debug` | Log level (debug/info/warn/error) |
| `DATABASE_HOST` | `localhost` | PostgreSQL host |
| `DATABASE_PORT` | `5433` | PostgreSQL port |
| `DATABASE_NAME` | `restaurant_db` | Database name |
| `DATABASE_USER` | `postgres` | Database user |
| `DATABASE_PASSWORD` | `postgres` | Database password |
| `DATABASE_MAX_OPEN_CONNS` | `100` | Max open connections |
| `DATABASE_MAX_IDLE_CONNS` | `10` | Max idle connections |
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_DB` | `1` | Redis database number |

---

## 🗄️ Database

### Schema
- **Table**: `restaurants`
  - Primary Key: `id` (UUID)
  - Unique Index: `(source, external_id)` WHERE deleted_at IS NULL
  - Supports: Google Maps, Tabelog, OpenTable sources

- **Table**: `user_favorites`
  - Primary Key: `id` (UUID)
  - Unique Index: `(user_id, restaurant_id)` WHERE deleted_at IS NULL

### Manual Migration
```bash
# Up
migrate -path migrations/restaurant \
  -database "postgresql://user:pass@host:5433/restaurant_db?sslmode=disable" \
  up

# Down (rollback)
migrate -path migrations/restaurant \
  -database "postgresql://user:pass@host:5433/restaurant_db?sslmode=disable" \
  down 1
```

---

## 🌐 API Endpoints

### Base URL
```
http://localhost:18082/api/v1
```

### Endpoints

#### 1. Health Check
```
GET /health
```

#### 2. Create Restaurant
```
POST /restaurants
Content-Type: application/json

{
  "name": "string",
  "source": "google|tabelog|opentable",
  "external_id": "string",
  "address": "string",
  "latitude": 0.0,
  "longitude": 0.0,
  "rating": 0.0,
  "price_range": "string",
  "cuisine_type": "string",
  "phone": "string",
  "website": "string"
}
```

#### 3. Get Restaurant
```
GET /restaurants/:id
```

#### 4. Search Restaurants
```
GET /restaurants/search?q=<query>&limit=<limit>&offset=<offset>
```

#### 5. Add to Favorites
```
POST /favorites
Content-Type: application/json

{
  "user_id": "uuid",
  "restaurant_id": "uuid"
}
```

#### 6. Get User Favorites
```
GET /users/:userId/favorites
```

---

## 🐛 Troubleshooting

### Service won't start
1. **Check PostgreSQL connection**:
   ```bash
   psql -h localhost -p 5433 -U postgres -d restaurant_db
   ```

2. **Check Redis connection**:
   ```bash
   redis-cli -h localhost -p 6379
   SELECT 1
   ```

3. **Check logs**:
   ```bash
   docker-compose logs restaurant-service
   ```

### Database connection error
- ✅ Verify `DATABASE_PORT=5433` (not 5432)
- ✅ Check if postgres-restaurant container is running
- ✅ Ensure migrations are successful

### Port already in use
```bash
# Check what's using port 18082
lsof -i :18082

# Kill process if needed
kill -9 <PID>
```

### Migration failed
```bash
# Check current migration version
migrate -path migrations/restaurant \
  -database "postgresql://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" \
  version

# Force version (use with caution)
migrate -path migrations/restaurant \
  -database "postgresql://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" \
  force <version>
```

---

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:18082/health
```

### Database Status
```bash
docker exec -it tabelogo-postgres-restaurant \
  psql -U postgres -d restaurant_db -c "SELECT COUNT(*) FROM restaurants;"
```

### Redis Cache Status
```bash
docker exec -it tabelogo-redis redis-cli
> SELECT 1
> KEYS map:*
```

---

## 🔐 Security Considerations

### Production Checklist
- [ ] Change default database password
- [ ] Enable SSL/TLS for database connections
- [ ] Set `ENVIRONMENT=production`
- [ ] Configure Redis password
- [ ] Enable rate limiting
- [ ] Add authentication middleware
- [ ] Configure CORS properly
- [ ] Enable request logging
- [ ] Set up firewall rules

### Environment Variables (Production)
```bash
ENVIRONMENT=production
LOG_LEVEL=info
DATABASE_PASSWORD=<strong-password>
REDIS_PASSWORD=<redis-password>
JWT_SECRET=<32-character-secret>
```

---

## 🔄 Integration with Other Services

### Map Service → Restaurant Service
**Flow**: Google Places API → Map Service → Restaurant Service

```go
// Map Service calls Restaurant Service to persist data
POST http://restaurant-service:18082/api/v1/restaurants
{
  "source": "google",
  "external_id": "<Google Place ID>",
  "name": "Restaurant Name",
  ...
}
```

### Spider Service → Restaurant Service
**Flow**: Tabelog Crawler → Spider Service → Restaurant Service

```go
// Spider Service calls Restaurant Service after scraping
POST http://restaurant-service:18082/api/v1/restaurants
{
  "source": "tabelog",
  "external_id": "https://tabelog.com/...",
  "name": "餐廳名稱",
  ...
}
```

### Frontend → Restaurant Service
**Flow**: User Search → Frontend → Restaurant Service

```typescript
// Frontend searches restaurants
GET http://restaurant-service:18082/api/v1/restaurants/search?q=sushi
```

---

## 📝 Changelog

### v1.0.0 (2025-12-02)
- ✅ Initial release
- ✅ Complete DDD architecture
- ✅ Restaurant & Favorite management
- ✅ External ID deduplication
- ✅ HTTP REST API (6 endpoints)
- ✅ Docker support
- ✅ Integration tests

---

## 📚 Additional Resources

- [Implementation Summary](IMPLEMENTATION_SUMMARY.md) - Detailed architecture documentation
- [Architecture.md](../../architecture.md) - Overall project architecture
- [Migration Summary](../../migrations/MIGRATIONS_SUMMARY.md) - Database schema details

---

## 🎉 Success!

Your Restaurant Service is now ready to:
- ✅ Receive data from Map Service (Google Places API)
- ✅ Receive data from Spider Service (Tabelog crawler)
- ✅ Serve restaurant data to Frontend
- ✅ Manage user favorites
- ✅ Prevent duplicate entries with `(source, external_id)` composite key

**Service Status**: 🟢 Production Ready (HTTP API)
**Future Enhancements**: gRPC support, Unit tests, Swagger docs
