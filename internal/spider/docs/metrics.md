# Spider Service Metrics Documentation

**Last Updated**: 2025-12-15  
**Version**: 1.0

---

## Overview

The Spider Service exposes 14 Prometheus metrics across 4 categories:
- **Scraping Metrics** (4): Track restaurant scraping operations
- **Job Processing Metrics** (4): Monitor background job processing
- **Cache Metrics** (3): Measure cache performance
- **Circuit Breaker Metrics** (2): Monitor circuit breaker state

All metrics use the `spider_` prefix and follow Prometheus naming conventions.

---

## Scraping Metrics

### 1. `spider_scrape_requests_total`

**Type**: Counter  
**Labels**: `status`  
**Description**: Total number of scrape requests by status

**Label Values**:
- `success`: Successful scrape
- `failed`: Failed scrape
- `cached`: Returned from cache

**Example Queries**:
```promql
# Total scrape requests per second
rate(spider_scrape_requests_total[5m])

# Success rate
rate(spider_scrape_requests_total{status="success"}[5m]) / 
rate(spider_scrape_requests_total[5m])

# Cache hit rate
rate(spider_scrape_requests_total{status="cached"}[5m]) / 
rate(spider_scrape_requests_total[5m])
```

**Alerting**:
```yaml
# Low success rate
- alert: SpiderHighFailureRate
  expr: |
    rate(spider_scrape_requests_total{status="failed"}[5m]) / 
    rate(spider_scrape_requests_total[5m]) > 0.1
  for: 5m
  annotations:
    summary: "Spider service has high failure rate (> 10%)"
```

---

### 2. `spider_scrape_duration_seconds`

**Type**: Histogram  
**Labels**: `operation`, `status`  
**Description**: Duration of scrape operations in seconds

**Label Values**:
- `operation`: `search`, `details`
- `status`: `success`, `failure`

**Buckets**: Default Prometheus buckets (0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10)

**Example Queries**:
```promql
# Average scrape duration by operation
rate(spider_scrape_duration_seconds_sum{status="success"}[5m]) / 
rate(spider_scrape_duration_seconds_count{status="success"}[5m])

# 95th percentile latency for search operations
histogram_quantile(0.95, 
  rate(spider_scrape_duration_seconds_bucket{operation="search"}[5m])
)

# Compare success vs failure duration
rate(spider_scrape_duration_seconds_sum{status="success"}[5m]) / 
rate(spider_scrape_duration_seconds_count{status="success"}[5m])
vs
rate(spider_scrape_duration_seconds_sum{status="failure"}[5m]) / 
rate(spider_scrape_duration_seconds_count{status="failure"}[5m])
```

**Alerting**:
```yaml
# Slow scraping
- alert: SpiderSlowScraping
  expr: |
    histogram_quantile(0.95, 
      rate(spider_scrape_duration_seconds_bucket[5m])
    ) > 10
  for: 5m
  annotations:
    summary: "Spider scraping is slow (p95 > 10s)"
```

---

### 3. `spider_scrape_errors_total`

**Type**: Counter  
**Labels**: `error_type`  
**Description**: Total number of scrape errors by type

**Label Values**:
- `network`: Network errors
- `parse`: HTML parsing errors
- `not_found`: Restaurant not found
- `circuit_breaker`: Circuit breaker open
- `timeout`: Request timeout

**Example Queries**:
```promql
# Error rate by type
rate(spider_scrape_errors_total[5m])

# Most common error type
topk(3, rate(spider_scrape_errors_total[5m]))

# Network error percentage
rate(spider_scrape_errors_total{error_type="network"}[5m]) / 
rate(spider_scrape_errors_total[5m])
```

**Alerting**:
```yaml
# High error rate
- alert: SpiderHighErrorRate
  expr: rate(spider_scrape_errors_total[5m]) > 1
  for: 5m
  annotations:
    summary: "Spider service has high error rate (> 1/s)"

# Circuit breaker frequently open
- alert: SpiderCircuitBreakerActive
  expr: |
    rate(spider_scrape_errors_total{error_type="circuit_breaker"}[5m]) > 0.1
  for: 2m
  annotations:
    summary: "Circuit breaker is frequently opening"
```

---

### 4. `spider_restaurants_scraped_total`

**Type**: Counter  
**Labels**: `status`  
**Description**: Total number of restaurants scraped by status

**Label Values**:
- `success`: Successfully scraped
- `failure`: Failed to scrape

**Example Queries**:
```promql
# Restaurants scraped per second
rate(spider_restaurants_scraped_total{status="success"}[5m])

# Success rate
rate(spider_restaurants_scraped_total{status="success"}[5m]) / 
rate(spider_restaurants_scraped_total[5m])

# Total restaurants scraped today
increase(spider_restaurants_scraped_total{status="success"}[24h])
```

---

## Job Processing Metrics

### 5. `spider_jobs_total`

**Type**: Counter  
**Labels**: `status`  
**Description**: Total number of jobs by status

**Label Values**:
- `pending`: Job created
- `running`: Job in progress
- `completed`: Job finished successfully
- `failed`: Job failed

**Example Queries**:
```promql
# Job completion rate
rate(spider_jobs_total{status="completed"}[5m])

# Job failure rate
rate(spider_jobs_total{status="failed"}[5m]) / 
rate(spider_jobs_total[5m])

# Jobs created vs completed
rate(spider_jobs_total{status="pending"}[5m]) vs 
rate(spider_jobs_total{status="completed"}[5m])
```

**Alerting**:
```yaml
# High job failure rate
- alert: SpiderHighJobFailureRate
  expr: |
    rate(spider_jobs_total{status="failed"}[5m]) / 
    rate(spider_jobs_total[5m]) > 0.2
  for: 5m
  annotations:
    summary: "Job failure rate is high (> 20%)"
```

---

### 6. `spider_job_duration_seconds`

**Type**: Histogram  
**Labels**: `status`  
**Description**: Duration of job processing in seconds

**Label Values**:
- `completed`: Successfully completed jobs
- `failed`: Failed jobs

**Buckets**: [0.5, 1, 2, 5, 10, 30, 60]

**Example Queries**:
```promql
# Average job duration
rate(spider_job_duration_seconds_sum[5m]) / 
rate(spider_job_duration_seconds_count[5m])

# 95th percentile job duration
histogram_quantile(0.95, 
  rate(spider_job_duration_seconds_bucket[5m])
)

# Jobs taking > 30 seconds
rate(spider_job_duration_seconds_bucket{le="30"}[5m])
```

**Alerting**:
```yaml
# Slow job processing
- alert: SpiderSlowJobProcessing
  expr: |
    histogram_quantile(0.95, 
      rate(spider_job_duration_seconds_bucket{status="completed"}[5m])
    ) > 60
  for: 10m
  annotations:
    summary: "Job processing is slow (p95 > 60s)"
```

---

### 7. `spider_worker_pool_size`

**Type**: Gauge  
**Labels**: None  
**Description**: Number of active workers in the pool

**Example Queries**:
```promql
# Current worker pool size
spider_worker_pool_size

# Worker pool utilization over time
spider_worker_pool_size
```

**Alerting**:
```yaml
# Worker pool at capacity
- alert: SpiderWorkerPoolFull
  expr: spider_worker_pool_size >= 20
  for: 5m
  annotations:
    summary: "Worker pool is at capacity"
```

---

### 8. `spider_job_queue_length`

**Type**: Gauge  
**Labels**: None  
**Description**: Number of pending jobs in the queue

**Example Queries**:
```promql
# Current queue length
spider_job_queue_length

# Queue growth rate
deriv(spider_job_queue_length[5m])

# Average queue length
avg_over_time(spider_job_queue_length[1h])
```

**Alerting**:
```yaml
# Queue backing up
- alert: SpiderQueueBacklog
  expr: spider_job_queue_length > 100
  for: 5m
  annotations:
    summary: "Job queue has significant backlog (> 100 jobs)"

# Queue growing rapidly
- alert: SpiderQueueGrowing
  expr: deriv(spider_job_queue_length[5m]) > 2
  for: 5m
  annotations:
    summary: "Job queue is growing rapidly"
```

---

## Cache Metrics

### 9. `spider_cache_hits_total`

**Type**: Counter  
**Labels**: `cache_type`  
**Description**: Total number of cache hits by cache type

**Label Values**:
- `result`: Result cache
- `job`: Job cache (future)

**Example Queries**:
```promql
# Cache hit rate
rate(spider_cache_hits_total[5m])

# Cache hit rate by type
rate(spider_cache_hits_total{cache_type="result"}[5m])

# Overall cache hit ratio
rate(spider_cache_hits_total[5m]) / 
(rate(spider_cache_hits_total[5m]) + rate(spider_cache_misses_total[5m]))
```

---

### 10. `spider_cache_misses_total`

**Type**: Counter  
**Labels**: `cache_type`  
**Description**: Total number of cache misses by cache type

**Label Values**:
- `result`: Result cache
- `job`: Job cache (future)

**Example Queries**:
```promql
# Cache miss rate
rate(spider_cache_misses_total[5m])

# Cache effectiveness
rate(spider_cache_hits_total{cache_type="result"}[5m]) / 
(rate(spider_cache_hits_total{cache_type="result"}[5m]) + 
 rate(spider_cache_misses_total{cache_type="result"}[5m]))
```

**Alerting**:
```yaml
# Low cache hit rate
- alert: SpiderLowCacheHitRate
  expr: |
    rate(spider_cache_hits_total{cache_type="result"}[5m]) / 
    (rate(spider_cache_hits_total{cache_type="result"}[5m]) + 
     rate(spider_cache_misses_total{cache_type="result"}[5m])) < 0.3
  for: 10m
  annotations:
    summary: "Cache hit rate is low (< 30%)"
```

---

### 11. `spider_cache_size_bytes`

**Type**: Gauge  
**Labels**: `cache_type`  
**Description**: Approximate size of cache in bytes by cache type

**Label Values**:
- `result`: Result cache
- `job`: Job cache (future)

**Example Queries**:
```promql
# Current cache size in MB
spider_cache_size_bytes / 1024 / 1024

# Cache size growth
deriv(spider_cache_size_bytes[1h])
```

**Alerting**:
```yaml
# Cache size too large
- alert: SpiderCacheTooLarge
  expr: spider_cache_size_bytes > 1073741824  # 1GB
  for: 5m
  annotations:
    summary: "Cache size exceeds 1GB"
```

---

## Circuit Breaker Metrics

### 12. `spider_circuit_breaker_state`

**Type**: Gauge  
**Labels**: `circuit`  
**Description**: Circuit breaker state (0=closed, 1=open, 2=half-open)

**Label Values**:
- `circuit`: `tabelog_scraper`

**State Values**:
- `0`: Closed (normal operation)
- `1`: Open (failing, requests blocked)
- `2`: Half-open (testing recovery)

**Example Queries**:
```promql
# Current circuit breaker state
spider_circuit_breaker_state

# Time in open state
changes(spider_circuit_breaker_state{circuit="tabelog_scraper"}[1h])
```

**Alerting**:
```yaml
# Circuit breaker open
- alert: SpiderCircuitBreakerOpen
  expr: spider_circuit_breaker_state{circuit="tabelog_scraper"} == 1
  for: 2m
  annotations:
    summary: "Circuit breaker is open - service degraded"
```

---

### 13. `spider_circuit_breaker_failures_total`

**Type**: Counter  
**Labels**: `circuit`  
**Description**: Total number of circuit breaker failures

**Label Values**:
- `circuit`: `tabelog_scraper`

**Example Queries**:
```promql
# Failure rate
rate(spider_circuit_breaker_failures_total[5m])

# Total failures in last hour
increase(spider_circuit_breaker_failures_total[1h])
```

---

## Grafana Dashboard Suggestions

### Overview Dashboard

**Panels**:
1. **Request Rate**: `rate(spider_scrape_requests_total[5m])`
2. **Success Rate**: Success rate calculation
3. **Average Latency**: p50, p95, p99 latencies
4. **Error Rate**: `rate(spider_scrape_errors_total[5m])`
5. **Cache Hit Rate**: Cache effectiveness
6. **Queue Length**: `spider_job_queue_length`
7. **Worker Pool**: `spider_worker_pool_size`
8. **Circuit Breaker State**: `spider_circuit_breaker_state`

### Performance Dashboard

**Panels**:
1. **Scrape Duration Heatmap**: Histogram visualization
2. **Job Duration Distribution**: Histogram visualization
3. **Restaurants per Second**: `rate(spider_restaurants_scraped_total[5m])`
4. **Cache Performance**: Hit/miss rates by type

### Errors Dashboard

**Panels**:
1. **Error Types**: Breakdown by `error_type`
2. **Failure Rate Trend**: Over time
3. **Circuit Breaker Activity**: State changes
4. **Failed Jobs**: `rate(spider_jobs_total{status="failed"}[5m])`

---

## Best Practices

### 1. Query Time Ranges
- Use `[5m]` for real-time monitoring
- Use `[1h]` or `[24h]` for trends
- Use `increase()` for totals over time

### 2. Alerting Thresholds
- Start conservative, tune based on actual behavior
- Use `for:` clause to avoid flapping
- Include context in annotations

### 3. Label Cardinality
- Current labels are low-cardinality (safe)
- Avoid adding high-cardinality labels (e.g., job_id, google_id)
- Use `cache_type` for future expansion

### 4. Metric Retention
- Default Prometheus retention: 15 days
- Consider longer retention for capacity planning
- Use recording rules for expensive queries

---

## Recording Rules (Optional)

For frequently used queries, create recording rules:

```yaml
groups:
  - name: spider_service
    interval: 30s
    rules:
      # Success rate
      - record: spider:scrape_success_rate:5m
        expr: |
          rate(spider_scrape_requests_total{status="success"}[5m]) / 
          rate(spider_scrape_requests_total[5m])
      
      # Cache hit rate
      - record: spider:cache_hit_rate:5m
        expr: |
          rate(spider_cache_hits_total[5m]) / 
          (rate(spider_cache_hits_total[5m]) + 
           rate(spider_cache_misses_total[5m]))
      
      # Average job duration
      - record: spider:job_duration_avg:5m
        expr: |
          rate(spider_job_duration_seconds_sum[5m]) / 
          rate(spider_job_duration_seconds_count[5m])
```

---

## Troubleshooting

### High Error Rate
1. Check `spider_scrape_errors_total` by `error_type`
2. Check circuit breaker state
3. Review application logs
4. Check network connectivity to Tabelog

### Slow Performance
1. Check `spider_scrape_duration_seconds` percentiles
2. Check `spider_job_duration_seconds` percentiles
3. Review worker pool size and queue length
4. Check cache hit rate

### Low Cache Hit Rate
1. Check `spider_cache_hits_total` vs `spider_cache_misses_total`
2. Verify cache TTL settings
3. Check cache size limits
4. Review cache invalidation logic

---

**For Questions**: Contact the Spider Service team  
**Metrics Endpoint**: `http://localhost:18084/metrics`  
**Prometheus Config**: See `deployments/prometheus/prometheus.yml`
