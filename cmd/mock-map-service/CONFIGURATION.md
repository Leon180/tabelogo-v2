# Mock Service Configuration

## Environment Variables

### Latency Simulation

Control response delay to simulate real Google API behavior:

```bash
# Disable latency (default - fast mode)
MOCK_LATENCY_ENABLED=false

# Enable latency simulation
MOCK_LATENCY_ENABLED=true
MOCK_LATENCY_MIN_MS=100    # Minimum latency in milliseconds
MOCK_LATENCY_MAX_MS=300    # Maximum latency in milliseconds
```

## Usage Examples

### Fast Mode (Recommended for K6 Testing)

```bash
# No latency - responses in <10ms
docker run -p 8085:8085 \
  -e MOCK_LATENCY_ENABLED=false \
  mock-map-service
```

**Use for**:
- ✅ K6 load testing
- ✅ Development
- ✅ CI/CD pipelines
- ✅ Finding code bottlenecks

### Realistic Mode (For Integration Testing)

```bash
# Simulate Google API latency (100-300ms)
docker run -p 8085:8085 \
  -e MOCK_LATENCY_ENABLED=true \
  -e MOCK_LATENCY_MIN_MS=100 \
  -e MOCK_LATENCY_MAX_MS=300 \
  mock-map-service
```

**Use for**:
- ✅ Testing timeout handling
- ✅ Testing retry logic
- ✅ Integration tests
- ✅ Simulating production conditions

### Slow Mode (For Stress Testing)

```bash
# Simulate slow API (500-1000ms)
docker run -p 8085:8085 \
  -e MOCK_LATENCY_ENABLED=true \
  -e MOCK_LATENCY_MIN_MS=500 \
  -e MOCK_LATENCY_MAX_MS=1000 \
  mock-map-service
```

**Use for**:
- ✅ Testing degraded performance
- ✅ Testing circuit breakers
- ✅ Finding timeout issues

## Docker Compose Examples

### Fast Mode

```yaml
mock-map-service:
  image: mock-map-service
  ports:
    - "8085:8085"
  environment:
    - MOCK_LATENCY_ENABLED=false
```

### Realistic Mode

```yaml
mock-map-service:
  image: mock-map-service
  ports:
    - "8085:8085"
  environment:
    - MOCK_LATENCY_ENABLED=true
    - MOCK_LATENCY_MIN_MS=100
    - MOCK_LATENCY_MAX_MS=300
```

## Testing Strategy

### Phase 1: Fast Mode (Find Code Bottlenecks)

```bash
# Use fast mode to isolate YOUR code's performance
MOCK_LATENCY_ENABLED=false k6 run tests/k6/restaurant_test.js
```

**Goal**: Find bottlenecks in:
- Database queries
- Redis operations
- Your business logic
- Connection pools

### Phase 2: Realistic Mode (Test Integration)

```bash
# Use realistic latency to test how your code handles slow APIs
MOCK_LATENCY_ENABLED=true \
MOCK_LATENCY_MIN_MS=100 \
MOCK_LATENCY_MAX_MS=300 \
k6 run tests/k6/restaurant_test.js
```

**Goal**: Test:
- Timeout configurations
- Retry logic
- Circuit breakers
- User experience

### Phase 3: Slow Mode (Stress Test)

```bash
# Use slow mode to test extreme conditions
MOCK_LATENCY_ENABLED=true \
MOCK_LATENCY_MIN_MS=500 \
MOCK_LATENCY_MAX_MS=1000 \
k6 run tests/k6/restaurant_test.js
```

**Goal**: Test:
- System behavior under stress
- Graceful degradation
- Error handling

## Recommendations

### For K6 Load Testing
```bash
MOCK_LATENCY_ENABLED=false  # ✅ Recommended
```
**Why**: Isolate your code's performance, fast iterations

### For Development
```bash
MOCK_LATENCY_ENABLED=false  # ✅ Recommended
```
**Why**: Fast feedback loop

### For CI/CD
```bash
MOCK_LATENCY_ENABLED=false  # ✅ Recommended
```
**Why**: Fast, stable tests

### For Integration Testing
```bash
MOCK_LATENCY_ENABLED=true   # ✅ Recommended
MOCK_LATENCY_MIN_MS=100
MOCK_LATENCY_MAX_MS=300
```
**Why**: Realistic simulation

## Monitoring Latency

Check current configuration:

```bash
curl http://localhost:8085/health
```

Response includes latency status:
```json
{
  "status": "healthy",
  "service": "mock-map-service",
  "version": "1.0.0",
  "latency": {
    "enabled": true,
    "min_ms": 100,
    "max_ms": 300
  }
}
```
