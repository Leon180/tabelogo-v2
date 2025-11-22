# Auth Service

Authentication and authorization microservice for Tabelogo v2.

## Features

- User registration and login
- JWT-based authentication (Access Token + Refresh Token)
- Password hashing with bcrypt
- Dual protocol support: gRPC + HTTP REST API
- Redis-based token storage
- PostgreSQL user persistence

## API Endpoints

### HTTP REST API (Port 8080)

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get tokens
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/validate` - Validate access token
- `GET /health` - Health check

### gRPC API (Port 9090)

- `Register` - Register a new user
- `Login` - Login and get tokens
- `RefreshToken` - Refresh access token
- `ValidateToken` - Validate access token

## Quick Start

### Using Docker Compose

```bash
# Start all services (PostgreSQL, Redis, Auth Service)
docker-compose up -d

# View logs
docker-compose logs -f auth-service

# Stop services
docker-compose down
```

### Local Development

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your configuration
# Make sure to set:
# - DB_NAME=auth_db
# - JWT_SECRET (min 32 characters)

# Run the service
go run main.go
```

### Build Binary

```bash
# Build
GOWORK=off go build -o ../../bin/auth-service .

# Run
../../bin/auth-service
```

## Environment Variables

See `.env.example` for all available configuration options.

### Required Variables

- `DB_NAME` - Database name
- `JWT_SECRET` - JWT signing secret (minimum 32 characters)

### Optional Variables

- `ENVIRONMENT` - Environment mode (development/staging/production)
- `LOG_LEVEL` - Log level (debug/info/warn/error)
- `SERVER_PORT` - HTTP server port (default: 8080)
- `GRPC_PORT` - gRPC server port (default: 9090)
- `DB_HOST` - PostgreSQL host (default: localhost)
- `REDIS_HOST` - Redis host (default: localhost)

## Testing

```bash
# Run unit tests
make test-unit

# Run integration tests (requires Docker)
make test-integration

# Run all tests
make test-all

# Generate coverage report
make test-coverage
```

## Architecture

This service follows Domain-Driven Design (DDD) principles:

```
cmd/auth-service/
├── main.go                 # Entry point
├── Dockerfile             # Container definition
├── docker-compose.yml     # Local development stack
└── .env.example          # Environment template

internal/auth/
├── domain/               # Domain layer (entities, repositories)
│   ├── model/           # Domain entities
│   ├── repository/      # Repository interfaces
│   └── errors/          # Domain errors
├── infrastructure/      # Infrastructure layer
│   ├── postgres/       # PostgreSQL implementation
│   └── redis/          # Redis implementation
├── application/        # Application layer
│   └── service.go     # Business logic
└── interfaces/         # Interface layer
    ├── grpc/          # gRPC handlers
    └── http/          # HTTP handlers
```

## Dependencies

- PostgreSQL 15+
- Redis 7+
- Go 1.23+

## License

See root LICENSE file.
