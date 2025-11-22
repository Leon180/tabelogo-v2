# Multi-Source Restaurant Aggregator Platform - Complete Architecture Design

## 1. Project Overview

A platform that aggregates restaurant information from multiple sources, providing restaurant search, booking, and other functionalities. It adopts a microservices architecture to demonstrate distributed system design and implementation capabilities.

---

## 2. Core Functional Modules

### 2.1 Functional Services
- **Auth Service**: User authentication and authorization (Implemented)
- **Booking Service**: Restaurant booking functionality (Integrating OpenTable API)
- **Map Service**: Map and navigation functionality (Integrating Google Maps API)
- **Spider Service**: Crawler microservice (Crawling Tabelog, etc.)
- **Mail Service**: Email notification service
- **Restaurant Service**: Restaurant data aggregation and query

---

## 3. Technical Architecture

### 3.1 Architecture Patterns
- **Microservices Architecture**
  - Inter-service communication: gRPC (Internal)
  - External API: RESTful API
  - Service Discovery: Consul or etcd (Planned)
  - API Gateway: Unified entry, routing, authentication

- **Domain-Driven Design (DDD)**
  - Layered Architecture: Presentation → Application → Domain → Infrastructure
  - Aggregate Root Design
  - Repository Pattern
  - Value Objects

### 3.2 Core Tech Stack
- **Language**: Go 1.24+
- **Dependency Injection**: Uber FX
- **Web Framework**: Gin
- **RPC Framework**: gRPC with Protocol Buffers
- **Message Queue**: Apache Kafka
- **Database**: PostgreSQL 15+ (Using GORM)
- **Cache**: Redis 7+
- **Concurrency**: Goroutines, Channels, Context

---

## 4. Data Layer Design

### 4.1 Database per Service Principle
Following microservices best practices, **each microservice has its own independent database** to achieve true service decoupling.

#### Database Allocation Strategy ✅

| Service | Database Name | Port | Main Tables | Description | Status |
|---------|---------------|------|-------------|-------------|--------|
| **Auth Service** | `auth_db` | **5432** | users, refresh_tokens | User authentication data | ✅ Implemented |
| **Restaurant Service** | `restaurant_db` | 5433 | restaurants, user_favorites | Restaurant master data, user favorites | Planned |
| **Booking Service** | `booking_db` | 5434 | bookings, booking_history | Booking data (Event Sourcing) | Planned |
| **Spider Service** | `spider_db` | 5435 | crawl_jobs, crawl_results | Crawler jobs and results | Planned |
| **Mail Service** | `mail_db` | 5436 | email_queue, email_logs | Email queue and logs | Planned |
| **Map Service** | No DB | - | - | Proxy for Google Maps API | - |

**Note**: For local development, Auth Service uses port **15432** to avoid conflicts if running separately.

#### Independent Redis Configuration

Each service uses a different Redis Database or independent Redis instance:

```yaml
Auth Service:     redis://redis:6379/0  (Session, Token Blacklist)
Restaurant Service: redis://redis:6379/1  (Restaurant Cache)
Booking Service:   redis://redis:6379/2  (Booking Cache)
Spider Service:    redis://redis:6379/3  (Rate Limiting, Distributed Lock)
API Gateway:       redis://redis:6379/4  (Rate Limiting, API Cache)
```

### 4.2 Cross-Service Data Query Strategy

#### 4.2.1 API Composition Pattern
When combining data from multiple services, the API Gateway or BFF (Backend for Frontend) is responsible:

**Example: Query user booking history (including restaurant info)**
```
1. API Gateway receives GET /api/v1/users/{userId}/bookings
2. Call Booking Service → Get booking list (with restaurant_id)
3. Call Restaurant Service → Batch query restaurant info by restaurant_ids
4. API Gateway combines data and returns
```

#### 4.2.2 CQRS Pattern (Command Query Responsibility Segregation)
For complex queries, create a **Read Model**:

- **Write Side (Command)**: Each microservice writes to its own DB
- **Read Side (Query)**: Sync to a dedicated query DB via events
- **Implementation**:
  - Use Kafka to send data change events
  - Query Service subscribes to events and updates Read Model (e.g., Elasticsearch)

### 4.3 Data Consistency

#### 4.3.1 Saga Pattern (Distributed Transactions)
Use **Choreography-based Saga** for cross-service transactions.

#### 4.3.2 Eventual Consistency
- Accept temporary data inconsistency
- Achieve consistency eventually via event-driven mechanisms

### 4.4 Database Technical Details

#### 4.4.1 PostgreSQL Design Standards
- **Schema Design**: 3NF
- **Primary Key**: UUID v4
- **Soft Delete**: deleted_at TIMESTAMP NULL
- **Audit Fields**: created_at, updated_at, created_by, updated_by

---

## 5. Message Queue Architecture

### 5.1 Kafka Usage Scenarios
- **Crawler Results Processing**
  - Topic: `spider-results`
- **Event-Driven Architecture**
  - Topic: `restaurant-events`
  - Topic: `booking-events`
  - Topic: `user-events`

---

## 6. API Design

### 6.1 RESTful API Standards
- **Versioning**: `/api/v1/...`
- **HTTP Methods**: GET, POST, PUT/PATCH, DELETE
- **Unified Response Format**:
```json
{
  "success": true,
  "data": {},
  "error": null,
  "meta": {
    "timestamp": "2025-11-17T10:00:00Z",
    "request_id": "uuid"
  }
}
```

### 6.2 API Documentation
- Swagger/OpenAPI 3.0

---

## 7. Authentication & Authorization

### 7.1 Authentication
- **JWT (JSON Web Token)**
  - Access Token (15 min)
  - Refresh Token (7 days, stored in Redis)
  - Token Blacklist

### 7.2 Authorization
- RBAC (Role-Based Access Control)
- Roles: Admin, User, Guest

### 7.3 Security
- Password hashing with bcrypt (cost=12)
- HTTPS/TLS 1.3
- SQL Injection protection
- XSS protection
- CSRF Token

---

## 8. Testing Strategy

### 8.1 Testing Levels
- **Unit Tests**: Go testing, testify/mock
- **Integration Tests**: Inter-service integration, testcontainers-go
- **E2E Tests**: API end-to-end testing

---

## 9. Monitoring & Observability

### 9.1 Metrics (Prometheus)
- Application metrics (HTTP/gRPC latency, error rate)
- System metrics (CPU, Memory)

### 9.2 Visualization (Grafana)
- Pre-built Dashboards

### 9.3 Logging (Zap + OpenTelemetry)
- Structured JSON logs
- Distributed Tracing (Jaeger)

---

## 10. DevOps & Deployment

### 10.1 Containerization (Docker)
- Multi-stage build
- Alpine base image

### 10.2 Orchestration (Docker Compose / K8s)
- Docker Compose for local development
- Kubernetes for production

### 10.3 CI/CD (GitHub Actions)
- Lint, Test, Build, Deploy

---

## 11. Development Workflow

### 11.1 Branching Strategy
- **main**: Production
- **develop**: Development
- **feature/***: New features
- **fix/***: Bug fixes

---

## 12. Database Schema Examples

### 12.1 Auth Service Database (`auth_db`) ✅

**Status**: Implemented
**Tables**: users, refresh_tokens

#### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
```

#### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP NULL
);
```
