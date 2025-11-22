# Tabelogo v2 - Docker Compose Architecture

## Architecture Overview

This project uses a microservices architecture with Docker Compose for container orchestration. There are two types of docker-compose configurations:

### 1. Root `docker-compose.yml` - Full System

**Purpose:** Starts the entire Tabelogo v2 system, including all microservices and infrastructure.

**Included Services:**
- **Infrastructure**
  - PostgreSQL (Independent for Auth, Restaurant, Booking)
  - Redis (Shared Cache)
  - Kafka + Zookeeper (Message Queue)
  - Prometheus (Monitoring)
  - Grafana (Visualization)

- **Microservices**
  - Auth Service (8080/HTTP, 9090/gRPC)
  - Restaurant Service (To be implemented)
  - Booking Service (To be implemented)
  - API Gateway (To be implemented)

**How to Start:**
```bash
# In project root
make up              # Start all services
make down            # Stop all services
make ps              # View service status
make logs            # View all logs
```

### 2. Service Directory `cmd/*/docker-compose.yml` - Single Service Development

**Purpose:** Used only for local development and testing of a single service.

**Features:**
- Uses different ports to avoid conflicts
- Starts only the service and its direct dependencies
- Suitable for rapid iteration

**Example - Auth Service Local Development:**
```bash
cd cmd/auth-service
docker-compose up -d    # Start Auth Service (Ports 18080/19090)
docker-compose down     # Stop
```

Or using Makefile:
```bash
make auth-up            # Start Auth Service (Local Dev Mode)
make auth-down          # Stop
make auth-logs          # View logs
```

## Port Allocation

### Root docker-compose (Production Mode)
| Service | HTTP | gRPC | Other |
|---------|------|------|-------|
| Auth Service | 8080 | 9090 | - |
| Restaurant Service | 8081 | 9091 | - |
| Booking Service | 8082 | 9092 | - |
| API Gateway | 8000 | - | - |
| PostgreSQL (Auth) | 5432 | - | - |
| PostgreSQL (Restaurant) | 5433 | - | - |
| PostgreSQL (Booking) | 5434 | - | - |
| Redis | 6379 | - | - |
| Kafka | - | - | 9092 |
| Prometheus | - | - | 9090 |
| Grafana | - | - | 3000 |

### Service Directory docker-compose (Development Mode)
| Service | HTTP | gRPC | DB | Redis |
|---------|------|------|----|-------|
| Auth Service | 18080 | 19090 | 15432 | 16379 |

## Usage Scenarios

### Scenario 1: Full System Testing
```bash
# Start entire system
make up

# Test inter-service communication
curl http://localhost:8080/health
```

### Scenario 2: Single Service Development
```bash
# Develop only Auth Service
cd cmd/auth-service
docker-compose up -d

# Or use Makefile
make auth-up
```

### Scenario 3: Adding New Microservice
1. Create service in `cmd/new-service/`
2. Create `cmd/new-service/Dockerfile`
3. Create `cmd/new-service/docker-compose.yml` (For development)
4. Add service definition to root `docker-compose.yml`
5. Update `Makefile` with relevant commands

## Network Architecture

All services are in the same Docker network `tabelogo-network` and can communicate with each other via service names:

```yaml
# Example: Restaurant Service calling Auth Service
AUTH_SERVICE_URL: http://auth-service:8080
AUTH_SERVICE_GRPC: auth-service:9090
```

## Data Persistence

All databases and caches use Docker Volumes for persistence:

```bash
# List volumes
docker volume ls | grep tabelogo

# Clean all data (Dangerous!)
make clean
```

## Best Practices

1. **Development**: Use service directory docker-compose
2. **Integration Testing**: Use root docker-compose
3. **Production Deployment**: Use Kubernetes (Future)
4. **Port Conflicts**: Ensure development mode uses different ports

## Troubleshooting

### Port Already in Use
```bash
# Check port usage
lsof -i :8080

# Use development mode (different ports)
make auth-up
```

### Service Fails to Start
```bash
# View logs
make logs

# Or specific service
make auth-logs
```

### Clean and Restart
```bash
# Stop and remove all containers and data
make clean

# Restart
make up
```
