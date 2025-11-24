# Environment Configuration

Create a `.env.local` file in the `web/` directory with the following variables:

```bash
# API Service URLs
NEXT_PUBLIC_MAP_SERVICE_URL=http://localhost:8080
NEXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:8081
NEXT_PUBLIC_RESTAURANT_SERVICE_URL=http://localhost:8082
NEXT_PUBLIC_BOOKING_SERVICE_URL=http://localhost:8083

# Google Maps API Key
NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here

# Default map center (Tokyo)
NEXT_PUBLIC_DEFAULT_LAT=35.6762
NEXT_PUBLIC_DEFAULT_LNG=139.6503
```

## Getting a Google Maps API Key

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the following APIs:
   - Maps JavaScript API
   - Places API
   - Geocoding API
4. Create credentials (API Key)
5. Copy the API key to your `.env.local` file
