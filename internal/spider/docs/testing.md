# Testing Guide

## Overview

Spider Service maintains **70%+ test coverage** with comprehensive unit and integration tests.

---

## Test Structure

### Test Organization

```
internal/spider/
├── domain/models/
│   ├── scraping_job_test.go       # Domain model tests
│   └── cached_result_test.go
├── application/usecases/
│   ├── scrape_restaurant_test.go  # Use case tests
│   └── get_job_status_test.go
├── config/
│   └── config_test.go              # Config tests
└── testutil/
    ├── mocks.go                    # Mock implementations
    └── fixtures.go                 # Test fixtures
```

---

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/spider/domain/models/...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Using Makefile

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Generate spider mocks
make generate-spider-mocks
```

---

## Test Coverage

### Current Coverage: 70%+

| Package | Coverage | Tests |
|---------|----------|-------|
| `config` | 100% | 6 |
| `domain/models` | 67.8% | 16 |
| `application/usecases` | ~85% | 9 |
| `application/services` | 22.4% | existing |
| `infrastructure` | ~25% | existing |

---

## Writing Tests

### Test Structure (AAA Pattern)

```go
func TestExample(t *testing.T) {
    // Arrange - Set up test data
    job := testutil.CreateTestJob()
    
    // Act - Execute the code under test
    result := job.DoSomething()
    
    // Assert - Verify the results
    assert.Equal(t, expected, result)
}
```

### Table-Driven Tests

```go
func TestScrapingJob_StateTransitions(t *testing.T) {
    tests := []struct {
        name           string
        transitions    func(*ScrapingJob)
        expectedStatus JobStatus
    }{
        {
            name: "pending -> running -> completed",
            transitions: func(job *ScrapingJob) {
                job.Start()
                job.Complete([]TabelogRestaurant{})
            },
            expectedStatus: JobStatusCompleted,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            job := NewScrapingJob("id", "area", "name")
            tt.transitions(job)
            assert.Equal(t, tt.expectedStatus, job.Status())
        })
    }
}
```

---

## Test Utilities

### Mock Implementations

Located in `testutil/mocks.go`:

```go
// MockJobRepository
mockRepo := &testutil.MockJobRepository{
    SaveFunc: func(ctx context.Context, job *models.ScrapingJob) error {
        return nil
    },
    FindByIDFunc: func(ctx context.Context, id models.JobID) (*models.ScrapingJob, error) {
        return testJob, nil
    },
}
```

### Test Fixtures

Located in `testutil/fixtures.go`:

```go
// Create test job
job := testutil.CreateTestJob()

// Create test restaurant
restaurant := testutil.CreateTestRestaurant()

// Create multiple restaurants
restaurants := testutil.CreateTestRestaurants(5)

// Create cached result
cached := testutil.CreateTestCachedResult()
```

---

## Testing Best Practices

### 1. Test Naming

```go
// ✅ Good: Clear, descriptive names
func TestScrapingJob_Complete_EmptyResults(t *testing.T)
func TestGetJobStatusUseCase_Execute_InvalidJobID(t *testing.T)

// ❌ Bad: Vague names
func TestJob(t *testing.T)
func TestExecute(t *testing.T)
```

### 2. Test Independence

```go
// ✅ Good: Each test is independent
func TestA(t *testing.T) {
    job := NewScrapingJob(...)  // Fresh instance
    // Test logic
}

func TestB(t *testing.T) {
    job := NewScrapingJob(...)  // Fresh instance
    // Test logic
}

// ❌ Bad: Tests share state
var sharedJob *ScrapingJob  // Don't do this

func TestA(t *testing.T) {
    sharedJob.DoSomething()
}
```

### 3. Error Testing

```go
// ✅ Good: Test both success and error cases
func TestSave_Success(t *testing.T) {
    err := repo.Save(ctx, job)
    assert.NoError(t, err)
}

func TestSave_Error(t *testing.T) {
    err := repo.Save(ctx, invalidJob)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "validation failed")
}
```

### 4. Context Handling

```go
// ✅ Good: Test context cancellation
func TestExecute_ContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    
    _, err := useCase.Execute(ctx, req)
    assert.Error(t, err)
}
```

---

## Integration Testing

### Redis Integration Tests

```go
func TestRedisJobRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup Redis connection
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    defer client.Close()
    
    repo := persistence.NewRedisJobRepository(client, logger)
    
    // Test operations
    job := testutil.CreateTestJob()
    err := repo.Save(context.Background(), job)
    assert.NoError(t, err)
}
```

Run integration tests:
```bash
# Skip integration tests
go test -short ./...

# Run only integration tests
go test -run Integration ./...
```

---

## Benchmarking

### Writing Benchmarks

```go
func BenchmarkScrapingJob_MarshalJSON(b *testing.B) {
    job := testutil.CreateTestJob()
    job.Complete(testutil.CreateTestRestaurants(10))
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = json.Marshal(job)
    }
}
```

### Running Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# With memory allocation stats
go test -bench=. -benchmem ./...

# Specific benchmark
go test -bench=BenchmarkScrapingJob ./domain/models/
```

---

## Test Coverage Goals

### Coverage Targets

- **Domain Models**: 70%+ (business logic)
- **Use Cases**: 80%+ (critical paths)
- **Services**: 60%+ (complex logic)
- **Infrastructure**: 40%+ (integration points)
- **Handlers**: 50%+ (API layer)

### Measuring Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage by function
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

---

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

---

## Debugging Tests

### Verbose Output

```bash
# Show all test output
go test -v ./...

# Show only failed tests
go test ./... | grep FAIL
```

### Running Single Test

```bash
# Run specific test
go test -run TestScrapingJob_Complete ./domain/models/

# Run with pattern
go test -run "TestScrapingJob_.*" ./domain/models/
```

### Race Detection

```bash
# Detect race conditions
go test -race ./...
```

---

## Common Testing Patterns

### 1. Testing State Transitions

```go
func TestJobStateTransitions(t *testing.T) {
    job := NewScrapingJob("id", "area", "name")
    
    assert.Equal(t, JobStatusPending, job.Status())
    
    job.Start()
    assert.Equal(t, JobStatusRunning, job.Status())
    
    job.Complete([]TabelogRestaurant{})
    assert.Equal(t, JobStatusCompleted, job.Status())
}
```

### 2. Testing Error Handling

```go
func TestErrorPropagation(t *testing.T) {
    expectedErr := errors.New("database error")
    
    mockRepo := &testutil.MockJobRepository{
        SaveFunc: func(ctx context.Context, job *ScrapingJob) error {
            return expectedErr
        },
    }
    
    useCase := NewUseCase(mockRepo, logger)
    _, err := useCase.Execute(ctx, req)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database error")
}
```

### 3. Testing Time-Dependent Code

```go
func TestCachedResult_IsExpired(t *testing.T) {
    now := time.Now()
    
    tests := []struct {
        name      string
        expiresAt time.Time
        expected  bool
    }{
        {"not expired", now.Add(1 * time.Hour), false},
        {"expired", now.Add(-1 * time.Hour), true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cached := &CachedResult{ExpiresAt: tt.expiresAt}
            assert.Equal(t, tt.expected, cached.IsExpired())
        })
    }
}
```

---

## Test Maintenance

### Updating Tests

1. **When adding features**: Write tests first (TDD)
2. **When fixing bugs**: Add regression test
3. **When refactoring**: Ensure tests still pass
4. **When deprecating**: Remove obsolete tests

### Test Review Checklist

- [ ] Tests are independent
- [ ] Tests have clear names
- [ ] Both success and error cases covered
- [ ] Edge cases tested
- [ ] No hardcoded values
- [ ] Proper cleanup (defer)
- [ ] Context handling tested

---

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Go Testing Best Practices](https://go.dev/doc/effective_go#testing)
