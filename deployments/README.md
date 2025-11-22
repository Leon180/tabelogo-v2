# Tabelogo v2 - Deployment Architecture

## Directory Structure

All deployment configurations are centralized in the `deployments/` directory:

```
deployments/
└── docker-compose/
    ├── docker-compose.yml     # Main system orchestration (All services)
    └── auth-service.yml       # Auth Service local development
```

## Usage

### 1. Main System (All Services)

Use this for full system testing or integration.

```bash
# Start all services
make up

# Stop all services
make down

# View logs
make logs
```

### 2. Individual Services (Local Development)

Use this when working on a specific service to isolate dependencies and avoid port conflicts.

#### Auth Service
- **File**: `deployments/docker-compose/auth-service.yml`
- **Ports**: HTTP 18080, gRPC 19090
- **DB Ports**: Postgres 15432, Redis 16379

```bash
# Start Auth Service
make auth-up

# Stop Auth Service
make auth-down

# View logs
make auth-logs
```

## Adding New Services

1. Create a new docker-compose file in `deployments/docker-compose/<service-name>.yml`.
2. Ensure it uses unique ports for local development (e.g., 18081, 19091).
3. Add the service to the main `deployments/docker-compose/docker-compose.yml`.
4. Update `Makefile` with new targets (e.g., `make <service>-up`).
