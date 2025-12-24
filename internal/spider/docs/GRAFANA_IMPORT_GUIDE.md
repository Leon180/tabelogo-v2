# Grafana Dashboard Import Guide

**Quick Start**: Import pre-configured Spider Service dashboard in 2 minutes!

---

## 📦 What's Included

The `grafana-dashboard-spider-service.json` file contains a complete dashboard with:

✅ **8 Pre-configured Panels**:
1. Scrape Requests per Second (Time series)
2. Success Rate (Gauge)
3. Cache Hit Rate (Gauge)
4. Scrape Duration Percentiles (Time series - p50, p95, p99)
5. Job Queue Length (Time series)
6. Error Rate by Type (Stacked area)
7. Active Workers (Stat)
8. Restaurants Scraped in 24h (Stat)

✅ **Auto-refresh**: 30 seconds  
✅ **Time range**: Last 6 hours  
✅ **Dark theme**: Enabled

---

## 🚀 Import Steps

### Step 1: Ensure Prometheus Data Source Exists

Before importing, make sure you have Prometheus configured:

1. Go to Grafana: `http://localhost:3001`
2. Login (admin/admin)
3. Navigate to **☰** → **Connections** → **Data sources**
4. Check if **Prometheus** exists
   - If not, click **Add data source** → **Prometheus**
   - Set URL: `http://prometheus:9090`
   - Click **Save & test**

---

### Step 2: Import Dashboard

#### Method 1: Via UI (Recommended)

1. **Open Grafana**: `http://localhost:3001`

2. **Navigate to Dashboards**:
   - Click **☰** (hamburger menu)
   - Click **Dashboards**

3. **Import Dashboard**:
   - Click **New** → **Import**
   - Click **Upload JSON file**
   - Select: `internal/spider/docs/grafana-dashboard-spider-service.json`
   - Or paste the JSON content directly

4. **Configure Import**:
   - **Name**: `Spider Service - Overview` (pre-filled)
   - **Folder**: Select or create folder (e.g., "Spider Service")
   - **UID**: `spider-service-overview` (pre-filled)
   - **Prometheus**: Select your Prometheus data source

5. **Import**:
   - Click **Import**
   - Dashboard will open automatically!

#### Method 2: Via File Path

If you prefer to copy the file path:

```bash
# Full path to dashboard JSON
/Users/lileon/goproject/tabelogov2/internal/spider/docs/grafana-dashboard-spider-service.json
```

---

### Step 3: Verify Dashboard

After import, you should see:

✅ **8 panels** displaying data  
✅ **Auto-refresh** every 30 seconds  
✅ **No errors** in panels

**If panels show "No data"**:
1. Check Spider Service is running: `http://localhost:18084/metrics`
2. Check Prometheus is scraping: `http://localhost:9091/targets`
3. Wait 30 seconds for first scrape

---

## 🎨 Dashboard Layout

```
┌─────────────────────────────────────────────────────────────────┐
│                    Spider Service - Overview                    │
├─────────────────┬─────────────┬─────────────┬──────────────────┤
│ Requests/sec    │ Success %   │ Cache Hit % │                  │
│ (Time series)   │ (Gauge)     │ (Gauge)     │                  │
├─────────────────┴─────────────┴─────────────┴──────────────────┤
│ Latency Percentiles (p50, p95, p99)  │ Job Queue Length        │
│ (Time series)                         │ (Time series)           │
├───────────────────────────────────────┼─────────────┬───────────┤
│ Error Rate by Type                    │ Workers     │ Scraped   │
│ (Stacked area)                        │ (Stat)      │ (Stat)    │
└───────────────────────────────────────┴─────────────┴───────────┘
```

---

## 🔧 Customization

### Change Time Range

- Click **time picker** in top right
- Select: Last 1h, 6h, 12h, 24h, or custom

### Adjust Refresh Rate

- Click **refresh icon** in top right
- Select: 5s, 10s, 30s, 1m, 5m, or off

### Edit Panels

1. Hover over panel title
2. Click **⋮** (three dots)
3. Click **Edit**
4. Modify query, visualization, or options
5. Click **Apply**
6. Click **💾 Save dashboard**

### Add More Panels

1. Click **Add** → **Visualization**
2. Select **Prometheus** data source
3. Enter PromQL query (see `metrics.md` for examples)
4. Configure visualization
5. Click **Apply**

---

## 📊 Example PromQL Queries

For additional panels, use these queries from `metrics.md`:

### Job Processing
```promql
# Job completion rate
rate(spider_jobs_total{status="completed"}[5m])

# Job failure rate
rate(spider_jobs_total{status="failed"}[5m])
```

### Circuit Breaker
```promql
# Circuit breaker state
spider_circuit_breaker_state{circuit="tabelog_scraper"}

# Circuit breaker failures
rate(spider_circuit_breaker_failures_total[5m])
```

### Cache Performance
```promql
# Cache size in MB
spider_cache_size_bytes{cache_type="result"} / 1024 / 1024
```

---

## 🔔 Setting Up Alerts (Optional)

### Add Alert to Success Rate Panel

1. Edit **Success Rate** panel
2. Click **Alert** tab
3. Click **Create alert rule from this panel**
4. Configure:
   ```
   Condition: WHEN last() OF query(A) IS BELOW 80
   For: 5m
   Summary: Spider service has high failure rate (< 80%)
   ```
5. Click **Save rule and exit**

### Add Alert to Latency Panel

1. Edit **Scrape Duration** panel
2. Create alert:
   ```
   Condition: WHEN last() OF query(B) IS ABOVE 10
   For: 5m
   Summary: Spider scraping is slow (p95 > 10s)
   ```

---

## 🐛 Troubleshooting

### "No data" in all panels

**Check**:
```bash
# 1. Spider Service is running
curl http://localhost:18084/metrics

# 2. Prometheus is scraping
open http://localhost:9091/targets

# 3. Metrics exist in Prometheus
open http://localhost:9091/graph
# Query: spider_scrape_requests_total
```

**Fix**:
- Ensure Spider Service is running
- Check Prometheus scrape config
- Wait 30-60 seconds for first scrape

### "Data source not found" error

**Fix**:
1. Go to **Connections** → **Data sources**
2. Add Prometheus data source
3. URL: `http://prometheus:9090`
4. Click **Save & test**
5. Re-import dashboard

### Panels show different data source

**Fix**:
1. Edit dashboard
2. Click **⚙️** (settings)
3. Click **Variables**
4. Update data source variable
5. Save dashboard

### Import fails with "UID already exists"

**Fix**:
1. During import, change UID to something unique
2. Or delete existing dashboard first
3. Then re-import

---

## 💡 Tips

1. **Star the dashboard**: Click ⭐ to add to favorites
2. **Share dashboard**: Click **Share** → Copy link
3. **Export dashboard**: **⚙️** → **JSON Model** → Copy
4. **Create folder**: Organize dashboards by service
5. **Use variables**: Add filters for dynamic dashboards

---

## 📚 Related Documentation

- [Metrics Documentation](./metrics.md) - All metrics explained
- [Grafana Setup Guide](./grafana_setup.md) - Manual setup steps
- [Prometheus Queries](./metrics.md#example-queries) - More PromQL examples

---

## ✅ Quick Checklist

- [ ] Prometheus data source configured
- [ ] Dashboard JSON file located
- [ ] Dashboard imported successfully
- [ ] All 8 panels showing data
- [ ] Auto-refresh working (30s)
- [ ] Dashboard saved
- [ ] (Optional) Alerts configured

---

**Ready to monitor!** 🎉

Your Spider Service dashboard is now live at:
`http://localhost:3001/d/spider-service-overview`
