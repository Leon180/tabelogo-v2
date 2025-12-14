# Frontend Migration Guide: Restaurant Service Integration

## Overview

This guide helps frontend developers migrate from calling the Map Service directly to using the Restaurant Service's `QuickSearchByPlaceID` endpoint.

### Why Migrate?

- **80% Cost Reduction**: Cache-first strategy reduces Google Maps API calls
- **Faster Response Times**: Cached responses return in <50ms
- **Better Reliability**: Graceful degradation with stale data fallback
- **Unified API**: Single source of truth for restaurant data

---

## API Changes

### Before: Map Service (Direct Call)

```typescript
// OLD: Direct gRPC call to Map Service
import { MapServiceClient } from '@/services/map';

const mapClient = new MapServiceClient();
const place = await mapClient.quickSearch({ placeId: 'ChIJ...' });
```

**Response**:
```typescript
interface Place {
  placeId: string;
  name: string;
  address: string;
  location: {
    latitude: number;
    longitude: number;
  };
  rating?: number;
  // ... other Map Service fields
}
```

---

### After: Restaurant Service (HTTP REST)

```typescript
// NEW: HTTP call to Restaurant Service
import { RestaurantService } from '@/services/restaurant';

const restaurantService = new RestaurantService();
const restaurant = await restaurantService.quickSearchByPlaceId('ChIJ...');
```

**Response**:
```typescript
interface Restaurant {
  id: string;                    // UUID
  name: string;
  source: string;                // "google"
  external_id: string;           // Google Place ID
  address: string;
  latitude: number;
  longitude: number;
  rating: number;
  price_range: string;
  cuisine_type: string;
  phone: string;
  website: string;
  view_count: number;            // NEW: popularity metric
  created_at: string;            // ISO 8601
  updated_at: string;            // ISO 8601
}
```

**Response Headers** (useful for debugging):
- `X-Cache-Status`: `HIT` or `MISS`
- `X-Data-Source`: `CACHE` or `MAP_SERVICE`
- `X-Data-Age`: Age in seconds (e.g., `"3600s"`)

---

## Migration Steps

### Step 1: Update API Client

Create or update `services/restaurant.ts`:

```typescript
import axios, { AxiosInstance } from 'axios';

export interface Restaurant {
  id: string;
  name: string;
  source: string;
  external_id: string;
  address: string;
  latitude: number;
  longitude: number;
  rating: number;
  price_range: string;
  cuisine_type: string;
  phone: string;
  website: string;
  view_count: number;
  created_at: string;
  updated_at: string;
}

export interface RestaurantResponse {
  restaurant: Restaurant;
}

export class RestaurantService {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_RESTAURANT_SERVICE_URL) {
    this.client = axios.create({
      baseURL,
      timeout: 5000,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  async quickSearchByPlaceId(placeId: string): Promise<Restaurant> {
    const response = await this.client.get<RestaurantResponse>(
      `/api/v1/restaurants/quick-search/${placeId}`
    );
    return response.data.restaurant;
  }

  // Helper to check if data is from cache
  isFromCache(response: any): boolean {
    return response.headers['x-cache-status'] === 'HIT';
  }

  // Helper to get data age
  getDataAge(response: any): number {
    const ageStr = response.headers['x-data-age'];
    return parseInt(ageStr) || 0;
  }
}
```

### Step 2: Update Environment Variables

Add to `.env.local`:

```bash
NEXT_PUBLIC_RESTAURANT_SERVICE_URL=http://localhost:8083
```

For production:

```bash
NEXT_PUBLIC_RESTAURANT_SERVICE_URL=https://api.yourdomain.com/restaurant-service
```

### Step 3: Update Components

**Before**:
```typescript
// components/RestaurantDetail.tsx
import { useEffect, useState } from 'react';
import { MapServiceClient } from '@/services/map';

export function RestaurantDetail({ placeId }: { placeId: string }) {
  const [place, setPlace] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPlace = async () => {
      try {
        const mapClient = new MapServiceClient();
        const data = await mapClient.quickSearch({ placeId });
        setPlace(data);
      } catch (error) {
        console.error('Failed to fetch place:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchPlace();
  }, [placeId]);

  if (loading) return <div>Loading...</div>;
  if (!place) return <div>Not found</div>;

  return (
    <div>
      <h1>{place.name}</h1>
      <p>{place.address}</p>
      <p>Rating: {place.rating}</p>
    </div>
  );
}
```

**After**:
```typescript
// components/RestaurantDetail.tsx
import { useEffect, useState } from 'react';
import { RestaurantService } from '@/services/restaurant';

export function RestaurantDetail({ placeId }: { placeId: string }) {
  const [restaurant, setRestaurant] = useState(null);
  const [loading, setLoading] = useState(true);
  const [cacheStatus, setCacheStatus] = useState('');

  useEffect(() => {
    const fetchRestaurant = async () => {
      try {
        const service = new RestaurantService();
        const response = await service.client.get(
          `/api/v1/restaurants/quick-search/${placeId}`
        );
        setRestaurant(response.data.restaurant);
        setCacheStatus(response.headers['x-cache-status']);
      } catch (error) {
        console.error('Failed to fetch restaurant:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchRestaurant();
  }, [placeId]);

  if (loading) return <div>Loading...</div>;
  if (!restaurant) return <div>Not found</div>;

  return (
    <div>
      <h1>{restaurant.name}</h1>
      <p>{restaurant.address}</p>
      <p>Rating: {restaurant.rating}</p>
      <p>Views: {restaurant.view_count}</p>
      {/* Optional: Show cache status for debugging */}
      {process.env.NODE_ENV === 'development' && (
        <small>Cache: {cacheStatus}</small>
      )}
    </div>
  );
}
```

### Step 4: Update Error Handling

```typescript
try {
  const restaurant = await restaurantService.quickSearchByPlaceId(placeId);
  // Success
} catch (error) {
  if (error.response?.status === 404) {
    // Restaurant not found
    showNotification('Restaurant not found');
  } else if (error.response?.status === 503) {
    // Map Service unavailable and no cached data
    showNotification('Service temporarily unavailable');
  } else {
    // Other errors
    showNotification('An error occurred');
  }
}
```

---

## Testing Checklist

### Functional Testing

- [ ] Search works with valid Place IDs
- [ ] Error handling for invalid Place IDs
- [ ] Error handling for non-existent restaurants
- [ ] All restaurant fields display correctly
- [ ] Loading states work properly

### Performance Testing

- [ ] First request (cache miss) completes in <200ms
- [ ] Subsequent requests (cache hit) complete in <50ms
- [ ] Check `X-Cache-Status` header in browser DevTools

### Error Scenarios

- [ ] Test with invalid Place ID format
- [ ] Test with non-existent Place ID
- [ ] Test when backend is down (should show error)
- [ ] Test network timeout

---

## Debugging Tips

### Check Cache Status

Open browser DevTools → Network tab → Click on the request → Headers:

```
X-Cache-Status: HIT
X-Data-Source: CACHE
X-Data-Age: 3600s
```

### Enable Verbose Logging

```typescript
const service = new RestaurantService();
service.client.interceptors.response.use(response => {
  console.log('Cache Status:', response.headers['x-cache-status']);
  console.log('Data Age:', response.headers['x-data-age']);
  return response;
});
```

---

## Rollback Plan

If issues arise, you can quickly rollback:

1. **Feature Flag**: Use a feature flag to toggle between old/new implementation
2. **Gradual Rollout**: Deploy to 10% → 50% → 100% of users
3. **Monitor Metrics**: Watch error rates and response times

```typescript
const USE_RESTAURANT_SERVICE = process.env.NEXT_PUBLIC_USE_RESTAURANT_SERVICE === 'true';

if (USE_RESTAURANT_SERVICE) {
  // New implementation
  const restaurant = await restaurantService.quickSearchByPlaceId(placeId);
} else {
  // Old implementation
  const place = await mapClient.quickSearch({ placeId });
}
```

---

## Support

- **API Documentation**: http://localhost:8083/swagger/index.html
- **Postman Collection**: See `docs/postman/restaurant-service.json`
- **Backend Team**: Contact for API issues or questions

---

## FAQ

**Q: What happens if the Map Service is down?**  
A: The Restaurant Service returns cached data (even if stale) with a 200 status. Check `X-Data-Age` header to see data freshness.

**Q: How often is cached data refreshed?**  
A: Data is refreshed when it's older than 3 days or when explicitly requested.

**Q: Can I force a fresh fetch from Map Service?**  
A: Not currently supported. The cache-first strategy is automatic.

**Q: What's the rate limit?**  
A: No rate limiting currently implemented. Monitor your usage.
