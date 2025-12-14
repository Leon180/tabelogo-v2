# Spider Service Architecture

## Overview

Spider Service follows **Clean Architecture** principles with clear separation of concerns across four main layers.

---

## Architectural Layers

### 1. Domain Layer (Core)

**Purpose**: Contains business logic and domain models

**Components**:
- **Models**: Core business entities
  - `ScrapingJob`: Job aggregate root
  - `TabelogRestaurant`: Restaurant value object
  - `CachedResult`: Cache entity
  
- **Repositories**: Data access interfaces
  - `JobRepository`: Job persistence
  - `ResultCacheRepository`: Cache operations

**Rules**:
- ✅ No external dependencies
- ✅ Pure business logic
- ✅ Framework-agnostic
- ❌ No infrastructure concerns

### 2. Application Layer

**Purpose**: Orchestrates business workflows

**Components**:
- **Use Cases**: Business operations
  - `ScrapeRestaurantUseCase`: Initiate scraping
  - `GetJobStatusUseCase`: Query job status
  
- **Services**: Domain services
  - `JobProcessor`: Async job processing
  - `DynamicRateLimiter`: Rate limiting logic

**Rules**:
- ✅ Depends on domain layer
- ✅ Coordinates workflows
- ❌ No framework dependencies
- ❌ No infrastructure details

### 3. Infrastructure Layer

**Purpose**: Implements technical capabilities

**Components**:
- **Persistence**: Data storage
  - `RedisJobRepository`: Redis-based job storage
  - `RedisResultCache`: Redis caching
  
- **Scraper**: Web scraping
  - `TabelogScraper`: Tabelog scraper implementation
  
- **Metrics**: Observability
  - `SpiderMetrics`: Prometheus metrics

**Rules**:
- ✅ Implements domain interfaces
- ✅ Framework-specific code
- ✅ External integrations
- ❌ No business logic

### 4. Interface Layer

**Purpose**: Exposes service capabilities

**Components**:
- **HTTP**: REST API
  - `SpiderHandler`: HTTP endpoints
  - SSE streaming support
  
- **gRPC**: Future RPC interface

**Rules**:
- ✅ Thin adapters
- ✅ Protocol-specific logic
- ❌ No business logic
- ❌ No direct domain access

---

## Data Flow

### Scraping Request Flow

```
1. Client Request
   ↓
2. HTTP Handler (Interface Layer)
   ↓
3. Use Case (Application Layer)
   ↓
4. Domain Model Creation
   ↓
5. Repository Save (Infrastructure)
   ↓
6. Job Queue Submission
   ↓
7. Worker Processing
   ↓
8. Scraper Execution
   ↓
9. Result Storage
   ↓
10. Cache Update
```

### SSE Streaming Flow

```
1. Client Connects (SSE)
   ↓
2. HTTP Handler Opens Stream
   ↓
3. Periodic Status Polling
   ↓
4. Use Case Query
   ↓
5. Repository Fetch
   ↓
6. Status Event Emission
   ↓
7. Stream Closure on Completion
```

---

## Key Design Patterns

### 1. Repository Pattern

**Purpose**: Abstract data access

```go
type JobRepository interface {
    Save(ctx context.Context, job *ScrapingJob) error
    FindByID(ctx context.Context, id JobID) (*ScrapingJob, error)
    // ...
}
```

**Benefits**:
- Testable (easy mocking)
- Swappable implementations
- Clear contracts

### 2. Use Case Pattern

**Purpose**: Encapsulate business operations

```go
type ScrapeRestaurantUseCase struct {
    jobRepo      JobRepository
    jobProcessor *JobProcessor
    logger       *zap.Logger
}

func (uc *ScrapeRestaurantUseCase) Execute(
    ctx context.Context,
    req ScrapeRestaurantRequest,
) (*ScrapeRestaurantResponse, error) {
    // Business logic here
}
```

**Benefits**:
- Single responsibility
- Clear input/output
- Testable workflows

### 3. Worker Pool Pattern

**Purpose**: Concurrent job processing

```go
type JobProcessor struct {
    workerCount int
    jobQueue    chan JobID
    workers     []*Worker
}
```

**Benefits**:
- Controlled concurrency
- Resource management
- Graceful shutdown

### 4. Circuit Breaker Pattern

**Purpose**: Failure protection

```go
type CircuitBreaker struct {
    state      State
    failures   int
    lastFailed time.Time
}
```

**Benefits**:
- Automatic recovery
- Prevents cascade failures
- Configurable thresholds

---

## Dependency Rules

### Dependency Direction

```
Interfaces → Application → Domain
     ↓
Infrastructure
```

**Rules**:
1. Inner layers never depend on outer layers
2. Domain has no dependencies
3. Infrastructure implements domain interfaces
4. Interfaces depend on application

### Dependency Injection

Using **Uber FX** for DI:

```go
fx.New(
    fx.Provide(
        // Infrastructure
        persistence.NewRedisJobRepository,
        persistence.NewRedisResultCache,
        
        // Application
        services.NewJobProcessor,
        usecases.NewScrapeRestaurantUseCase,
        
        // Interfaces
        http.NewSpiderHandler,
    ),
    fx.Invoke(registerHTTPHandlers),
)
```

---

## Concurrency Model

### Worker Pool Architecture

```
                    ┌──────────────┐
                    │  Job Queue   │
                    │  (buffered)  │
                    └──────┬───────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
          ▼                ▼                ▼
    ┌─────────┐      ┌─────────┐      ┌─────────┐
    │Worker 1 │      │Worker 2 │  ... │Worker N │
    └────┬────┘      └────┬────┘      └────┬────┘
         │                │                │
         └────────────────┼────────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │   Scraper    │
                   └──────────────┘
```

### Goroutine Management

**Lifecycle**:
1. Start workers on service startup
2. Workers listen on job queue
3. Process jobs concurrently
4. Graceful shutdown on context cancellation

**Safety**:
- Panic recovery in all goroutines
- Context propagation
- Proper cleanup on shutdown

---

## Error Handling Strategy

### Error Types

1. **Domain Errors**: Business rule violations
2. **Infrastructure Errors**: External system failures
3. **Validation Errors**: Input validation failures

### Error Propagation

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to save job: %w", err)
}
```

### Error Recovery

- Circuit breaker for external calls
- Retry logic with exponential backoff
- Panic recovery in goroutines

---

## Caching Strategy

### Cache Layers

1. **Application Cache**: In-memory (future)
2. **Distributed Cache**: Redis (current)

### Cache Policy

- **TTL**: 24 hours default
- **Invalidation**: On new scrape
- **Key Format**: `spider:result:{google_id}`

### Cache Flow

```
Request → Check Cache → Hit? → Return
                      ↓ Miss
                   Scrape → Store → Return
```

---

## Testing Strategy

### Test Pyramid

```
        ┌─────────┐
        │   E2E   │  (Few)
        ├─────────┤
        │Integration│ (Some)
        ├─────────┤
        │   Unit   │  (Many)
        └─────────┘
```

### Test Coverage by Layer

- **Domain**: 67.8% (high value)
- **Application**: ~85% (critical paths)
- **Infrastructure**: ~25% (integration tests)
- **Interfaces**: 0% (deferred)

### Mocking Strategy

- Mock all external dependencies
- Use interfaces for testability
- Provide test fixtures

---

## Performance Considerations

### Optimization Points

1. **Worker Pool**: Tune worker count
2. **Rate Limiting**: Balance speed vs. politeness
3. **Caching**: Reduce redundant scrapes
4. **Connection Pooling**: Redis connections

### Scalability

- **Horizontal**: Multiple service instances
- **Vertical**: Increase worker count
- **Cache**: Redis cluster for high load

---

## Security Considerations

1. **Rate Limiting**: Prevent abuse
2. **Input Validation**: Sanitize all inputs
3. **Error Messages**: No sensitive data leakage
4. **Logging**: Sanitize logged data

---

## Future Enhancements

### Planned Features

1. **gRPC Interface**: For service-to-service communication
2. **Distributed Tracing**: OpenTelemetry integration
3. **Advanced Caching**: Multi-level cache
4. **Job Priorities**: Priority queue
5. **Batch Processing**: Bulk scraping support

### Architectural Evolution

1. **Event Sourcing**: For audit trail
2. **CQRS**: Separate read/write models
3. **Saga Pattern**: For complex workflows

---

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
