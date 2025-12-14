# Load Testing Guide

## Prerequisites

Install k6:
```bash
# macOS
brew install k6

# Or download from https://k6.io/docs/getting-started/installation/
```

## Running Load Tests

### Basic Test
```bash
# Run the load test
k6 run scripts/load-test.js
```

### Custom Configuration
```bash
# Run with 200 virtual users for 5 minutes
k6 run --vus 200 --duration 5m scripts/load-test.js

# Run with specific stages
k6 run --stage 1m:50,3m:100,1m:0 scripts/load-test.js
```

### Output to InfluxDB (for Grafana visualization)
```bash
k6 run --out influxdb=http://localhost:8086/k6 scripts/load-test.js
```

## Test Scenarios

The load test includes:
- **Ramp-up**: 30s to 10 users → 1m to 50 users → 2m to 100 users
- **Sustained Load**: 1m at 100 users
- **Ramp-down**: 30s to 0 users

## Success Criteria

- ✅ **P95 Response Time**: < 200ms
- ✅ **Error Rate**: < 5%
- ✅ **Cache Hit Rate**: > 70%

## Interpreting Results

### Good Performance
```
✓ status is 200........................: 100%
✓ response time < 500ms................: 100%
http_req_duration......................: avg=45ms  p(95)=120ms
cache_hit_rate.........................: 78%
```

### Poor Performance (needs optimization)
```
✗ status is 200........................: 95%
✗ response time < 500ms................: 85%
http_req_duration......................: avg=250ms p(95)=850ms
cache_hit_rate.........................: 45%
```

## Monitoring During Tests

1. **Watch Prometheus metrics**:
   ```bash
   watch -n 1 'curl -s http://localhost:18082/metrics | grep restaurant_'
   ```

2. **Check Grafana dashboard**:
   - Open: http://localhost:3000
   - Dashboard: "Restaurant Service - Cache Performance"

3. **Monitor logs**:
   ```bash
   docker-compose logs -f restaurant-service
   ```

## Common Issues

### High Error Rate
- Check Map Service is running
- Verify database connections
- Check for timeout issues

### Low Cache Hit Rate
- Increase test duration (cache needs time to warm up)
- Use more duplicate place IDs in test data
- Check TTL configuration

### Slow Response Times
- Check database query performance
- Monitor Map Service response times
- Look for connection pool exhaustion

## Advanced Testing

### Stress Test (find breaking point)
```bash
k6 run --vus 500 --duration 10m scripts/load-test.js
```

### Spike Test (sudden traffic surge)
```bash
k6 run --stage 10s:0,1m:1000,10s:0 scripts/load-test.js
```

### Soak Test (long-duration stability)
```bash
k6 run --vus 100 --duration 2h scripts/load-test.js
```
