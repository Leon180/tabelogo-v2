# API Reference: QuickSearchByPlaceID

## Endpoint Details

**Method**: `GET`  
**Path**: `/api/v1/restaurants/quick-search/{place_id}`  
**Authentication**: None (currently)  
**Content-Type**: `application/json`

---

## Request Specification

### Path Parameters

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `place_id` | string | Yes | Google Place ID | `ChIJN1t_tDeuEmsRUsoyG83frY4` |

### Query Parameters

None

### Request Headers

| Header | Required | Description |
|--------|----------|-------------|
| `Accept` | No | Response format (default: `application/json`) |

---

## Response Specification

### Success Response (200 OK)

**Headers**:
- `Content-Type`: `application/json`
- `X-Cache-Status`: `HIT` or `MISS`
- `X-Data-Source`: `CACHE` or `MAP_SERVICE`
- `X-Data-Age`: Age in seconds (e.g., `"3600s"`)

**Body**:
```json
{
  "restaurant": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Sushi Dai",
    "source": "google",
    "external_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
    "address": "5 Chome-2-1 Tsukiji, Chuo City, Tokyo 104-0045, Japan",
    "latitude": 35.6654,
    "longitude": 139.7707,
    "rating": 4.5,
    "price_range": "$$",
    "cuisine_type": "Sushi",
    "phone": "+81 3-1234-5678",
    "website": "https://example.com",
    "view_count": 1523,
    "created_at": "2025-12-01T10:00:00Z",
    "updated_at": "2025-12-06T15:30:00Z"
  }
}
```

### Error Responses

#### 400 Bad Request
Invalid Place ID format.

```json
{
  "error": "invalid_place_id",
  "message": "Place ID is required"
}
```

#### 404 Not Found
Restaurant not found for the given Place ID.

```json
{
  "error": "not_found",
  "message": "Restaurant not found"
}
```

#### 500 Internal Server Error
Server error occurred.

```json
{
  "error": "internal_error",
  "message": "Failed to search restaurant"
}
```

#### 503 Service Unavailable
Map Service is unavailable and no cached data exists.

```json
{
  "error": "service_unavailable",
  "message": "Map Service unavailable and no cached data"
}
```

---

## Data Models

### Restaurant

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Unique restaurant identifier |
| `name` | string | Restaurant name |
| `source` | string | Data source (e.g., "google") |
| `external_id` | string | External ID from source (Google Place ID) |
| `address` | string | Full address |
| `latitude` | number | Latitude coordinate |
| `longitude` | number | Longitude coordinate |
| `rating` | number | Rating (0.0 - 5.0) |
| `price_range` | string | Price range (e.g., "$", "$$", "$$$") |
| `cuisine_type` | string | Type of cuisine |
| `phone` | string | Phone number |
| `website` | string | Website URL |
| `view_count` | number | Number of times viewed |
| `created_at` | string (ISO 8601) | Creation timestamp |
| `updated_at` | string (ISO 8601) | Last update timestamp |

### ErrorResponse

| Field | Type | Description |
|-------|------|-------------|
| `error` | string | Error code |
| `message` | string | Human-readable error message |

---

## Performance Characteristics

### Cache-First Strategy

This endpoint implements a cache-first approach:

1. **Cache Hit** (data < 3 days old):
   - Returns cached data immediately
   - Response time: **<50ms**
   - Headers: `X-Cache-Status: HIT`, `X-Data-Source: CACHE`

2. **Cache Miss** (no data or data > 3 days old):
   - Fetches fresh data from Map Service
   - Updates cache
   - Response time: **<200ms**
   - Headers: `X-Cache-Status: MISS`, `X-Data-Source: MAP_SERVICE`

3. **Graceful Degradation** (Map Service unavailable):
   - Returns stale cached data if available
   - Response time: **<50ms**
   - Headers: `X-Cache-Status: HIT`, `X-Data-Source: CACHE`
   - Check `X-Data-Age` to determine staleness

### Expected Response Times

| Scenario | Response Time | Cache Status |
|----------|---------------|--------------|
| Cache hit (fresh) | <50ms | HIT |
| Cache miss | <200ms | MISS |
| Map Service timeout | <50ms | HIT (stale) |

---

## Rate Limiting

**Current Status**: No rate limiting implemented

**Recommendation**: Implement client-side caching and debouncing for search inputs.

---

## Examples

### cURL

```bash
# Basic request
curl -X GET "http://localhost:8083/api/v1/restaurants/quick-search/ChIJN1t_tDeuEmsRUsoyG83frY4" \
  -H "Accept: application/json"

# With verbose output to see headers
curl -v -X GET "http://localhost:8083/api/v1/restaurants/quick-search/ChIJN1t_tDeuEmsRUsoyG83frY4"
```

### JavaScript/TypeScript (fetch)

```typescript
async function quickSearch(placeId: string) {
  const response = await fetch(
    `http://localhost:8083/api/v1/restaurants/quick-search/${placeId}`,
    {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    }
  );

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  // Check cache status
  const cacheStatus = response.headers.get('X-Cache-Status');
  const dataAge = response.headers.get('X-Data-Age');
  console.log(`Cache: ${cacheStatus}, Age: ${dataAge}`);

  const data = await response.json();
  return data.restaurant;
}

// Usage
try {
  const restaurant = await quickSearch('ChIJN1t_tDeuEmsRUsoyG83frY4');
  console.log(restaurant.name);
} catch (error) {
  console.error('Failed to fetch restaurant:', error);
}
```

### JavaScript/TypeScript (axios)

```typescript
import axios from 'axios';

const client = axios.create({
  baseURL: 'http://localhost:8083',
  timeout: 5000,
});

async function quickSearch(placeId: string) {
  try {
    const response = await client.get(
      `/api/v1/restaurants/quick-search/${placeId}`
    );

    // Access headers
    console.log('Cache Status:', response.headers['x-cache-status']);
    console.log('Data Age:', response.headers['x-data-age']);

    return response.data.restaurant;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      if (error.response?.status === 404) {
        console.error('Restaurant not found');
      } else if (error.response?.status === 503) {
        console.error('Service unavailable');
      }
    }
    throw error;
  }
}
```

### React Hook

```typescript
import { useState, useEffect } from 'react';
import axios from 'axios';

interface Restaurant {
  id: string;
  name: string;
  address: string;
  latitude: number;
  longitude: number;
  rating: number;
  // ... other fields
}

export function useRestaurant(placeId: string) {
  const [restaurant, setRestaurant] = useState<Restaurant | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [cacheStatus, setCacheStatus] = useState<string>('');

  useEffect(() => {
    const fetchRestaurant = async () => {
      try {
        setLoading(true);
        const response = await axios.get(
          `http://localhost:8083/api/v1/restaurants/quick-search/${placeId}`
        );
        setRestaurant(response.data.restaurant);
        setCacheStatus(response.headers['x-cache-status']);
        setError(null);
      } catch (err) {
        setError(err as Error);
      } finally {
        setLoading(false);
      }
    };

    if (placeId) {
      fetchRestaurant();
    }
  }, [placeId]);

  return { restaurant, loading, error, cacheStatus };
}

// Usage in component
function RestaurantDetail({ placeId }: { placeId: string }) {
  const { restaurant, loading, error, cacheStatus } = useRestaurant(placeId);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  if (!restaurant) return <div>Not found</div>;

  return (
    <div>
      <h1>{restaurant.name}</h1>
      <p>{restaurant.address}</p>
      <small>Cache: {cacheStatus}</small>
    </div>
  );
}
```

---

## Monitoring & Debugging

### Response Headers

Always check these headers for debugging:

```
X-Cache-Status: HIT | MISS
X-Data-Source: CACHE | MAP_SERVICE
X-Data-Age: <seconds>s
```

### Logging

The backend logs all requests with:
- Place ID
- Client IP
- Cache hit/miss status
- Response time
- Errors (if any)

### Metrics

Available Prometheus metrics:
- `restaurant_quick_search_total` - Total requests
- `restaurant_cache_hits_total` - Cache hits
- `restaurant_map_service_calls_total` - Map Service calls

---

## Best Practices

1. **Cache Response Headers**: Store cache status for analytics
2. **Handle Stale Data**: Check `X-Data-Age` and show indicator if > 24 hours
3. **Error Handling**: Implement proper error handling for all status codes
4. **Debouncing**: Debounce search inputs to reduce API calls
5. **Loading States**: Show loading indicators during requests
6. **Retry Logic**: Implement exponential backoff for failed requests

---

## Support & Resources

- **Swagger UI**: http://localhost:8083/swagger/index.html
- **Postman Collection**: `docs/postman/restaurant-service.json`
- **Migration Guide**: `docs/frontend-migration-guide.md`
- **Backend Repository**: Contact backend team for access
