# Grafana Dashboard Setup Guide - Spider Service

**Last Updated**: 2025-12-15

---

## 📋 Prerequisites

- Grafana running at `http://localhost:3001`
- Prometheus running at `http://prometheus:9090` (internal Docker network)
- Spider Service running and exposing metrics

---

## 🚀 Step-by-Step Setup

### Step 1: Access Grafana

1. Open browser: `http://localhost:3001`
2. Login credentials:
   - **Username**: `admin`
   - **Password**: `admin`
3. (Optional) Skip password change or set new password

---

### Step 2: Add Prometheus Data Source

#### 2.1 Navigate to Data Sources
1. Click **☰** (hamburger menu) in top left
2. Click **Connections** → **Data sources**
3. Click **Add data source** button

#### 2.2 Select Prometheus
1. Find and click **Prometheus** in the list

#### 2.3 Configure Prometheus
Fill in the following settings:

**Connection**:
- **Name**: `Prometheus` (default)
- **URL**: `http://prometheus:9090`
  - ⚠️ Use `prometheus` (service name) not `localhost`
  - This is the internal Docker network address

**HTTP**:
- **Access**: `Server (default)`

**Auth**: Leave all unchecked

**Scrape interval**: `15s` (default)

#### 2.4 Save & Test
1. Scroll to bottom
2. Click **Save & test**
3. Should see: ✅ "Data source is working"

---

### Step 3: Create Spider Service Dashboard

#### 3.1 Create New Dashboard
1. Click **☰** → **Dashboards**
2. Click **New** → **New Dashboard**
3. Click **Add visualization**
4. Select **Prometheus** data source

#### 3.2 Dashboard Settings
1. Click **⚙️** (gear icon) in top right
2. Set **Dashboard name**: `Spider Service - Overview`
3. Add **Description**: `Monitoring dashboard for Spider Service metrics`
4. Click **Save dashboard**

---

## 📊 Example Panels

### Panel 1: Request Rate

**Panel Title**: `Scrape Requests per Second`

**PromQL Query**:
```promql
rate(spider_scrape_requests_total[5m])
```

**Visualization**: Time series (Line chart)

**Legend**: `{{status}}`

**Panel Options**:
- **Unit**: `requests/sec`
- **Min**: `0`

---

### Panel 2: Success Rate

**Panel Title**: `Success Rate (%)`

**PromQL Query**:
```promql
100 * (
  rate(spider_scrape_requests_total{status="success"}[5m]) / 
  rate(spider_scrape_requests_total[5m])
)
```

**Visualization**: Stat (Single value)

**Panel Options**:
- **Unit**: `percent (0-100)`
- **Thresholds**:
  - Red: `< 80`
  - Yellow: `< 95`
  - Green: `>= 95`

---

### Panel 3: Latency (p50, p95, p99)

**Panel Title**: `Scrape Duration (Percentiles)`

**PromQL Queries** (add 3 queries):

Query A (p50):
```promql
histogram_quantile(0.50, 
  rate(spider_scrape_duration_seconds_bucket[5m])
)
```

Query B (p95):
```promql
histogram_quantile(0.95, 
  rate(spider_scrape_duration_seconds_bucket[5m])
)
```

Query C (p99):
```promql
histogram_quantile(0.99, 
  rate(spider_scrape_duration_seconds_bucket[5m])
)
```

**Visualization**: Time series

**Legend**: 
- Query A: `p50`
- Query B: `p95`
- Query C: `p99`

**Panel Options**:
- **Unit**: `seconds (s)`
- **Min**: `0`

---

### Panel 4: Cache Hit Rate

**Panel Title**: `Cache Hit Rate`

**PromQL Query**:
```promql
100 * (
  rate(spider_cache_hits_total[5m]) / 
  (rate(spider_cache_hits_total[5m]) + rate(spider_cache_misses_total[5m]))
)
```

**Visualization**: Gauge

**Panel Options**:
- **Unit**: `percent (0-100)`
- **Min**: `0`
- **Max**: `100`
- **Thresholds**:
  - Red: `< 30`
  - Yellow: `< 60`
  - Green: `>= 60`

---

### Panel 5: Error Rate

**Panel Title**: `Error Rate by Type`

**PromQL Query**:
```promql
rate(spider_scrape_errors_total[5m])
```

**Visualization**: Time series (Stacked area)

**Legend**: `{{error_type}}`

**Panel Options**:
- **Unit**: `errors/sec`
- **Min**: `0`

---

### Panel 6: Job Queue Length

**Panel Title**: `Job Queue Length`

**PromQL Query**:
```promql
spider_job_queue_length
```

**Visualization**: Time series

**Panel Options**:
- **Unit**: `short`
- **Min**: `0`

---

### Panel 7: Worker Pool Utilization

**Panel Title**: `Active Workers`

**PromQL Query**:
```promql
spider_worker_pool_size
```

**Visualization**: Stat

**Panel Options**:
- **Unit**: `short`
- **Color**: Value-based

---

### Panel 8: Restaurants Scraped

**Panel Title**: `Restaurants Scraped (Last 24h)`

**PromQL Query**:
```promql
increase(spider_restaurants_scraped_total{status="success"}[24h])
```

**Visualization**: Stat

**Panel Options**:
- **Unit**: `short`
- **Color**: Green

---

## 🎨 Dashboard Layout Suggestions

### Row 1: Overview (4 panels)
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Request Rate    │ Success Rate    │ Cache Hit Rate  │ Active Workers  │
│ (Time series)   │ (Stat)          │ (Gauge)         │ (Stat)          │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

### Row 2: Performance (2 panels)
```
┌──────────────────────────────────┬──────────────────────────────────┐
│ Latency Percentiles              │ Job Queue Length                 │
│ (Time series - p50, p95, p99)    │ (Time series)                    │
└──────────────────────────────────┴──────────────────────────────────┘
```

### Row 3: Errors & Metrics (2 panels)
```
┌──────────────────────────────────┬──────────────────────────────────┐
│ Error Rate by Type               │ Restaurants Scraped (24h)        │
│ (Stacked area)                   │ (Stat)                           │
└──────────────────────────────────┴──────────────────────────────────┘
```

---

## 🔧 Panel Creation Steps

### For Each Panel:

1. **Add Panel**:
   - Click **Add** → **Visualization**
   - Select **Prometheus** data source

2. **Configure Query**:
   - Paste PromQL query in **Metrics browser**
   - Click **Run queries** to test

3. **Set Visualization**:
   - Select visualization type (Time series, Stat, Gauge, etc.)
   - Configure legend format

4. **Panel Options** (right sidebar):
   - Set **Title**
   - Set **Unit**
   - Configure **Thresholds** (if applicable)
   - Set **Min/Max** values

5. **Save**:
   - Click **Apply** in top right
   - Click **Save dashboard** (💾 icon)

---

## 📱 Dashboard Variables (Optional)

Add variables for filtering:

### Variable 1: Status
- **Name**: `status`
- **Type**: Query
- **Query**: `label_values(spider_scrape_requests_total, status)`
- **Multi-value**: Yes
- **Include All**: Yes

### Variable 2: Operation
- **Name**: `operation`
- **Type**: Query
- **Query**: `label_values(spider_scrape_duration_seconds, operation)`
- **Multi-value**: Yes
- **Include All**: Yes

**Usage in queries**:
```promql
rate(spider_scrape_requests_total{status=~"$status"}[5m])
```

---

## 🔔 Setting Up Alerts (Optional)

### Alert 1: High Failure Rate

1. Edit **Success Rate** panel
2. Click **Alert** tab
3. Click **Create alert rule from this panel**
4. Configure:
   - **Condition**: `WHEN last() OF query(A) IS BELOW 80`
   - **For**: `5m`
   - **Summary**: `Spider service has high failure rate`

### Alert 2: Slow Performance

1. Edit **Latency p95** panel
2. Create alert:
   - **Condition**: `WHEN last() OF query(B) IS ABOVE 10`
   - **For**: `5m`
   - **Summary**: `Spider scraping is slow (p95 > 10s)`

---

## 💡 Tips

1. **Auto-refresh**: Set dashboard refresh to `30s` or `1m`
2. **Time range**: Use `Last 6 hours` or `Last 24 hours`
3. **Save often**: Click 💾 to save dashboard changes
4. **Export**: Share dashboard JSON via **⚙️** → **JSON Model**
5. **Templating**: Use variables for dynamic filtering

---

## 🐛 Troubleshooting

### "No data" in panels
- Check Prometheus is scraping Spider Service
- Visit `http://localhost:9091/targets` to verify
- Ensure Spider Service is running: `http://localhost:18084/metrics`

### "Data source not found"
- Verify Prometheus data source is configured
- Check URL is `http://prometheus:9090` (not localhost)

### Queries return errors
- Test query in Prometheus UI: `http://localhost:9091`
- Check metric names match exactly
- Verify time range has data

---

## 📚 Additional Resources

- [Grafana Documentation](https://grafana.com/docs/)
- [PromQL Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Spider Service Metrics](../metrics.md)

---

**Ready to create your dashboard!** 🎨
