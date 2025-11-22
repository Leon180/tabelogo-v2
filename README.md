# Multi-Source Restaurant Aggregator (Tabelogo V2)

A microservices-based restaurant information aggregator platform that integrates multiple restaurant data sources, providing restaurant search, booking, and review functionalities.

## 🏗 Architecture Features

- **Microservices Architecture**: Independent development, deployment, and scaling for each service
- **Database per Service**: Each microservice has its own independent database instance
- **DDD Design**: Domain-Driven Design with clear layered architecture
- **Event-Driven**: Event-driven architecture using Kafka
- **gRPC Communication**: Efficient gRPC communication between services
- **Full Monitoring**: Observability with Prometheus + Grafana + Jaeger

## 🎯 Core Services

| Service | Port | Database | Description |
|---------|------|----------|-------------|
| API Gateway | 8080 | - | Unified entry point, routing, authentication |
| Auth Service | 8081/9081 | auth_db | User authentication and authorization |
| Restaurant Service | 8082/9082 | restaurant_db | Restaurant data management |
| Booking Service | 8083/9083 | booking_db | Booking functionality |
| Spider Service | 8084/9084 | spider_db | Crawler service |
| Mail Service | 8085/9085 | mail_db | Email notification |
| Map Service | 8086/9086 | - | Maps and navigation |

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.24+
- Make

### Local Development Setup

```bash
# 1. Clone repository
git clone https://github.com/Leon180/tabelogo-v2.git
cd tabelogov2

# 2. Initialize project (create .env file)
make init

# 3. Start all infrastructure (PostgreSQL, Redis, Kafka, etc.)
make up

# 4. Check container status
make ps
```

### Available Make Commands

```bash
make help          # Show all available commands
make init          # Initialize project
make up            # Start all Docker containers
make down          # Stop all containers
make restart       # Restart all containers
make logs          # View container logs
make ps            # View container status
make clean         # Clean up all containers and volumes
make build         # Build all microservices
make test          # Run all tests
make lint          # Run code linter
make migrate-up    # Run database migrations
make migrate-down  # Rollback database migrations
```

## 🗄️ Database Architecture

### Database per Service Principle

Each microservice has its own independent PostgreSQL database instance:

| Database | Port | Usage |
|----------|------|-------|
| auth_db | 5432 | User authentication data |
| restaurant_db | 5433 | Restaurant master data |
| booking_db | 5434 | Booking data |
| spider_db | 5435 | Crawler jobs and results |
| mail_db | 5436 | Email queue and logs |

### Redis Configuration

Different Redis Database Numbers are used to distinguish services:

- DB 0: Auth Service (Session, Token Blacklist)
- DB 1: Restaurant Service (Restaurant Cache)
- DB 2: Booking Service (Booking Cache)
- DB 3: Spider Service (Rate Limiting, Distributed Lock)
- DB 4: API Gateway (Rate Limiting, API Cache)

## 📊 Monitoring & Observability

- **Kafka UI**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090

## 🔧 Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin
- **gRPC**: Protocol Buffers
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Message Queue**: Apache Kafka
- **Monitoring**: Prometheus + Grafana + Jaeger
- **Logging**: Zap + OpenTelemetry
- **Containerization**: Docker + Docker Compose

## 📁 Project Structure

```
tabelogov2/
├── cmd/                      # Entry points for each microservice (independent go.mod)
├── internal/                 # Internal code separated by service
├── pkg/                      # Shared packages (independent go.mod)
├── api/proto/                # gRPC Protocol Buffers definitions
├── migrations/               # Database migrations for each service
├── deployments/              # Docker & Kubernetes configurations
├── scripts/                  # Build and deployment scripts
├── tests/                    # Tests
└── docs/                     # Documentation
```

Detailed architecture documentation: [architecture.md](docs/architecture.md)

## 🔐 Environment Variables

Copy `.env.example` to `.env` and modify the settings:

```bash
cp .env.example .env
```

Important variables:
- `JWT_SECRET`: JWT signing secret (Must change for production)
- `GOOGLE_MAPS_API_KEY`: Google Maps API Key
- `SMTP_*`: Email service settings

## 🧪 Testing

```bash
# Run tests for all services
make test

# Run tests for a specific service
cd cmd/auth-service && go test ./... -v
```

## 📝 License

MIT License

## 👥 Author

Leon Li
