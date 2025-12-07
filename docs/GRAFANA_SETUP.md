# Grafana Setup Guide

## Access Grafana

**URL**: http://localhost:3001  
**Username**: admin  
**Password**: admin

## Step 1: Add Prometheus Data Source

1. Open Grafana at http://localhost:3001
2. Login with admin/admin (skip password change if prompted)
3. Navigate to **Connections** > **Data sources** (or use the gear icon ⚙️)
4. Click **Add data source**
5. Select **Prometheus**
6. Configure:
   - **Name**: Prometheus
   - **URL**: `http://prometheus:9090`
   - **Access**: Server (default)
7. Click **Save & test**
8. You should see "Data source is working" ✅

## Step 2: Import Restaurant Service Dashboard

1. Navigate to **Dashboards** (four squares icon)
2. Click **New** > **Import**
3. Upload the dashboard JSON file:
   ```
   deployments/grafana/dashboards/restaurant-service.json
   ```
4. Or paste the JSON content directly
5. Select **Prometheus** as the datasource
6. Click **Import**

## Dashboard Panels

The Restaurant Service dashboard includes:

### Performance Metrics
- **Cache Hit Rate**: Target >70% for cost savings
- **Response Time (p95)**: Should be <500ms
- **Request Rate**: Traffic by endpoint

### Service Health
- **Map Service API Calls**: Success vs Error rates
- **Error Rate**: Alert if >1%
- **Stale Data Returns**: Fallback usage count

### Cache Metrics
- **Total Restaurants Cached**: Cache size
- **Sync Duration**: Map Service call time

## Alerts (Optional)

You can set up alerts for:
- Cache hit rate < 70%
- Error rate > 1%
- Response time > 500ms (p95)
- Map Service API errors

## Troubleshooting

### Prometheus Not Showing Data
```bash
# Check Prometheus targets
curl http://localhost:9091/api/v1/targets

# Verify Restaurant Service metrics
curl http://localhost:18082/metrics
```

### Grafana Can't Connect to Prometheus
- Ensure both containers are on the same network: `tabelogo-network`
- Use `http://prometheus:9090` (not localhost) for datasource URL
- Check docker logs: `docker logs tabelogo-grafana`

### No Metrics Showing
- Wait 15-30 seconds for first scrape
- Make some API calls to generate metrics:
  ```bash
  curl http://localhost:18082/api/v1/restaurants/quick-search/ChIJN5Nz71W3j4ARhx5bwpTQEGg
  ```

## Useful Prometheus Queries

```promql
# Cache hit rate
rate(restaurant_cache_hits_total[5m]) / (rate(restaurant_cache_hits_total[5m]) + rate(restaurant_cache_misses_total[5m])) * 100

# Error rate
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) * 100

# p95 response time
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Map Service call rate
rate(restaurant_map_service_calls_total[5m])
```

## Next Steps

1. ✅ Access Grafana
2. ✅ Add Prometheus datasource
3. ✅ Import dashboard
4. Monitor your metrics!
5. Set up alerts (optional)
6. Create custom dashboards for specific needs
