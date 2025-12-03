# Restaurant Service

Restaurant management microservice for Tabelogo v2.

## Features

- Restaurant CRUD operations
- Favorite management (add/remove, notes, tags, visit tracking)
- Multi-source restaurant deduplication (Google Places, Tabelog)
- Location-based search with geospatial queries
- Cuisine type filtering
- View count tracking
- Domain-Driven Design (DDD) architecture
- PostgreSQL persistence with GORM
- HTTP REST API with Gin
- Comprehensive unit tests (98% coverage)

## API Endpoints

### HTTP REST API (Port 18082)

#### Restaurant Operations
- `POST /api/v1/restaurants` - Create a new restaurant
- `GET /api/v1/restaurants/:id` - Get restaurant by ID
- `GET /api/v1/restaurants/external/:source/:id` - Get restaurant by external ID (source: google/tabelog)
- `PUT /api/v1/restaurants/:id` - Update restaurant
- `DELETE /api/v1/restaurants/:id` - Delete restaurant
- `GET /api/v1/restaurants` - List restaurants (with pagination)
- `GET /api/v1/restaurants/search` - Search restaurants by name
- `GET /api/v1/restaurants/nearby` - Find restaurants by location
- `GET /api/v1/restaurants/cuisine/:type` - Find restaurants by cuisine type
- `POST /api/v1/restaurants/:id/view` - Increment view count

#### Favorite Operations
- `POST /api/v1/favorites` - Add restaurant to favorites
- `DELETE /api/v1/favorites/:restaurant_id` - Remove from favorites
- `GET /api/v1/favorites` - Get user's favorites
- `PUT /api/v1/favorites/:restaurant_id/notes` - Update favorite notes
- `POST /api/v1/favorites/:restaurant_id/tags` - Add tag to favorite
- `DELETE /api/v1/favorites/:restaurant_id/tags/:tag` - Remove tag from favorite
- `POST /api/v1/favorites/:restaurant_id/visit` - Record a visit
- `GET /api/v1/favorites/:restaurant_id/check` - Check if restaurant is favorited

#### Health Check
- `GET /health` - Service health check

## Quick Start

### Using Docker Compose

```bash
# Start all services (PostgreSQL, Redis, Restaurant Service)
cd deployments/docker-compose
docker-compose up -d restaurant-service

# View logs
docker-compose logs -f restaurant-service

# Stop services
docker-compose down
```

### Local Development

```bash
# Navigate to service directory
cd cmd/restaurant-service

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
# Make sure to set:
# - DATABASE_NAME=restaurant_db
# - DATABASE_PORT=5433 (to avoid conflict with auth DB on 5432)

# Run database migrations
cd ../../migrations/restaurant
psql -U postgres -d restaurant_db -f 001_initial_schema.sql

# Run the service
cd ../../cmd/restaurant-service
go run main.go
```

### Build Binary

```bash
# Build (from cmd/restaurant-service directory)
GOWORK=off go build -o ../../bin/restaurant-service .

# Run
../../bin/restaurant-service
```

## Environment Variables

See `.env` for all available configuration options.

### Required Variables

- `DATABASE_NAME` - Database name (default: restaurant_db)
- `DATABASE_PORT` - PostgreSQL port (default: 5433)

### Optional Variables

- `SERVER_ENVIRONMENT` - Environment mode (development/staging/production)
- `LOG_LEVEL` - Log level (debug/info/warn/error)
- `SERVER_PORT` - HTTP server port (default: 18082)
- `DATABASE_HOST` - PostgreSQL host (default: localhost)
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_DB` - Redis database number (default: 1)

## Testing

```bash
# Run all unit tests
cd /Users/lileon/goproject/tabelogov2
go test ./internal/restaurant/domain/model/... ./internal/restaurant/application/...

# Run with coverage
go test -coverprofile=coverage.out ./internal/restaurant/domain/model/... ./internal/restaurant/application/...
go tool cover -html=coverage.out

# Run with verbose output
go test -v ./internal/restaurant/domain/model/... ./internal/restaurant/application/...

# Run specific test
go test -v -run TestRestaurantService_CreateRestaurant_Success ./internal/restaurant/application/...

# Generate coverage report
go tool cover -func=coverage.out | tail -1
# Expected output: total: (statements) 98.0%
```

## Test Coverage

Current test coverage: **98.0%** (exceeds 90% requirement)

- **Domain Layer**: 98.2% (60 tests)
  - Location Value Object: 100%
  - Restaurant Aggregate Root: 98%
  - Favorite Aggregate Root: 100%

- **Application Layer**: 97.9% (62 tests)
  - 19 Restaurant service methods tested
  - 9 Favorite service methods tested
  - All error scenarios covered

See [TESTING_SUMMARY.md](/Users/lileon/goproject/tabelogov2/internal/restaurant/TESTING_SUMMARY.md) for detailed coverage report.

## Architecture

This service follows Domain-Driven Design (DDD) principles:

```
cmd/restaurant-service/
├── main.go                # Entry point with Uber FX
├── Dockerfile            # Container definition
├── entrypoint.sh         # Container startup script
├── .env                  # Environment configuration
├── go.mod                # Go module definition
└── README.md             # This file

internal/restaurant/
├── domain/              # Domain layer (pure business logic)
│   ├── model/          # Domain entities (Restaurant, Favorite, Location)
│   ├── repository/     # Repository interfaces
│   └── errors/         # Domain errors
├── infrastructure/     # Infrastructure layer
│   ├── postgres/      # PostgreSQL repository implementations
│   └── redis/         # Redis cache (future)
├── application/       # Application layer
│   ├── service.go    # Business logic orchestration
│   └── dto.go        # Request/Response DTOs
└── interfaces/        # Interface layer
    └── http/         # HTTP handlers (Gin)

migrations/restaurant/
└── 001_initial_schema.sql  # Database schema
```

## Database Schema

### Restaurants Table
- Primary key: `id` (UUID)
- Unique constraint: `(source, external_id)` for deduplication
- Geospatial index on `(latitude, longitude)` for location-based queries
- Full-text search index on `name` and `address`

### User Favorites Table
- Primary key: `id` (UUID)
- Unique constraint: `(user_id, restaurant_id)`
- Support for notes, tags (TEXT[]), visit tracking
- Soft delete support

## External ID Deduplication

The service prevents duplicate restaurants from different sources:

```go
// (source, external_id) is unique
Source:     "google"     | "tabelog"
ExternalID: "ChIJ..."    | "https://tabelog.com/..."
```

This ensures that Map Service (Google) and Spider Service (Tabelog) don't create duplicates.

## Integration Points

### Map Service Integration
- Map Service creates restaurants with `source=google` and `external_id=place_id`
- Checks existence via `GET /api/v1/restaurants/external/google/:place_id`

### Spider Service Integration
- Spider Service creates restaurants with `source=tabelog` and `external_id=url`
- Checks existence via `GET /api/v1/restaurants/external/tabelog/:url`

### Frontend Integration
- User favorites management
- Restaurant search and filtering
- Location-based discovery

## Dependencies

- PostgreSQL 15+
- Redis 7+ (for caching, optional)
- Go 1.24+
- PostGIS extension (for geospatial queries)

## Development Notes

### Running Migrations

```bash
# Connect to PostgreSQL
psql -U postgres -h localhost -p 5433 -d restaurant_db

# Run migrations
\i migrations/restaurant/001_initial_schema.sql
```

### Testing with curl

```bash
# Create restaurant
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

# Get restaurant by external ID
curl http://localhost:18082/api/v1/restaurants/external/google/ChIJTest123

# Add to favorites
curl -X POST http://localhost:18082/api/v1/favorites \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "uuid-here",
    "restaurant_id": "uuid-here"
  }'
```

## Production Deployment

See [Dockerfile](Dockerfile) for containerized deployment.

The service includes:
- Multi-stage Docker build for minimal image size
- Automatic database migrations on startup
- Health check endpoint
- Graceful shutdown handling

## License

See root LICENSE file.
