import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const cacheHitRate = new Rate('cache_hit_rate');
const responseTime = new Trend('response_time');

// Test configuration
export const options = {
    stages: [
        { duration: '30s', target: 10 },   // Ramp up to 10 users
        { duration: '1m', target: 50 },    // Ramp up to 50 users
        { duration: '2m', target: 100 },   // Ramp up to 100 users
        { duration: '1m', target: 100 },   // Stay at 100 users
        { duration: '30s', target: 0 },    // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95% of requests should be below 200ms
        http_req_failed: ['rate<0.05'],   // Error rate should be below 5%
        cache_hit_rate: ['rate>0.7'],     // Cache hit rate should be above 70%
    },
};

// Sample place IDs for testing (mix of cache hits and misses)
const placeIds = [
    'ChIJN1t_tDeuEmsRUsoyG83frY4', // Google Sydney
    'ChIJrTLr-GyuEmsRBfy61i59si0', // Sydney Opera House
    'ChIJ3S-JXmauEmsRUcIaWtf4MzE', // Sydney Harbour Bridge
    'ChIJP3Sa8ziYEmsRUKgyFmh9AQM', // Bondi Beach
    'ChIJISz8NjyuEmsRFTQ9Iw7Ear8', // Darling Harbour
    'ChIJ21P2pu65EmsRfBLOG7HRiY8', // Taronga Zoo
    'ChIJrTLr-GyuEmsRBfy61i59si0', // Duplicate for cache hit
    'ChIJN1t_tDeuEmsRUsoyG83frY4', // Duplicate for cache hit
];

export default function () {
    // Select a random place ID
    const placeId = placeIds[Math.floor(Math.random() * placeIds.length)];

    // Make request to QuickSearchByPlaceID endpoint
    const url = `http://localhost:18082/api/v1/restaurants/quick-search?place_id=${placeId}`;
    const startTime = new Date();

    const response = http.get(url, {
        headers: {
            'Content-Type': 'application/json',
        },
    });

    const duration = new Date() - startTime;
    responseTime.add(duration);

    // Check response
    const success = check(response, {
        'status is 200': (r) => r.status === 200,
        'response has data': (r) => r.json('data') !== undefined,
        'response time < 500ms': (r) => r.timings.duration < 500,
    });

    // Track cache hits (if response time is very fast, likely a cache hit)
    if (duration < 50) {
        cacheHitRate.add(1);
    } else {
        cacheHitRate.add(0);
    }

    // Random sleep between 1-3 seconds
    sleep(Math.random() * 2 + 1);
}

// Summary function to display results
export function handleSummary(data) {
    return {
        'stdout': textSummary(data, { indent: ' ', enableColors: true }),
        'load-test-results.json': JSON.stringify(data),
    };
}

function textSummary(data, options) {
    const indent = options.indent || '';
    const enableColors = options.enableColors || false;

    let summary = '\n';
    summary += `${indent}✓ Test completed\n`;
    summary += `${indent}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n`;

    summary += `${indent}Requests:\n`;
    summary += `${indent}  Total: ${data.metrics.http_reqs.values.count}\n`;
    summary += `${indent}  Failed: ${data.metrics.http_req_failed.values.rate * 100}%\n\n`;

    summary += `${indent}Response Time:\n`;
    summary += `${indent}  Avg: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms\n`;
    summary += `${indent}  P95: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms\n`;
    summary += `${indent}  P99: ${data.metrics.http_req_duration.values['p(99)'].toFixed(2)}ms\n\n`;

    if (data.metrics.cache_hit_rate) {
        summary += `${indent}Cache Performance:\n`;
        summary += `${indent}  Hit Rate: ${(data.metrics.cache_hit_rate.values.rate * 100).toFixed(2)}%\n\n`;
    }

    return summary;
}
