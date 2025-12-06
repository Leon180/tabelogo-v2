#!/bin/bash
# Metrics Verification Test Script

echo "=== Restaurant Service Metrics Verification ==="
echo ""

# 1. Check if metrics endpoint is accessible
echo "1. Testing metrics endpoint..."
METRICS_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:18082/metrics)
if [ "$METRICS_RESPONSE" = "200" ]; then
    echo "✅ Metrics endpoint accessible (HTTP 200)"
else
    echo "❌ Metrics endpoint failed (HTTP $METRICS_RESPONSE)"
    exit 1
fi
echo ""

# 2. Check if Restaurant Service metrics are exposed
echo "2. Checking Restaurant Service metrics..."
METRICS=$(curl -s http://localhost:18082/metrics | grep "^restaurant_")

if echo "$METRICS" | grep -q "restaurant_cache_hits_total"; then
    echo "✅ restaurant_cache_hits_total found"
else
    echo "❌ restaurant_cache_hits_total NOT found"
fi

if echo "$METRICS" | grep -q "restaurant_cache_misses_total"; then
    echo "✅ restaurant_cache_misses_total found"
else
    echo "❌ restaurant_cache_misses_total NOT found"
fi

if echo "$METRICS" | grep -q "restaurant_map_service_calls_total"; then
    echo "✅ restaurant_map_service_calls_total found"
else
    echo "❌ restaurant_map_service_calls_total NOT found"
fi

if echo "$METRICS" | grep -q "restaurant_sync_duration_seconds"; then
    echo "✅ restaurant_sync_duration_seconds found"
else
    echo "❌ restaurant_sync_duration_seconds NOT found"
fi

if echo "$METRICS" | grep -q "restaurant_stale_data_returns_total"; then
    echo "✅ restaurant_stale_data_returns_total found"
else
    echo "❌ restaurant_stale_data_returns_total NOT found"
fi
echo ""

# 3. Show current metric values
echo "3. Current metric values:"
echo "------------------------"
curl -s http://localhost:18082/metrics | grep -E "restaurant_(cache|map_service|stale|sync)" | grep -v "^#"
echo ""

# 4. Test metric labels
echo "4. Checking metric labels..."
if curl -s http://localhost:18082/metrics | grep -q 'restaurant_map_service_calls_total{status="success"}'; then
    echo "✅ Map Service calls metric has 'success' label"
fi

if curl -s http://localhost:18082/metrics | grep -q 'restaurant_map_service_calls_total{status="error"}'; then
    echo "✅ Map Service calls metric has 'error' label"
fi
echo ""

echo "=== Verification Complete ==="
