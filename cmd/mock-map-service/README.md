# Mock Map Service

A lightweight mock service that simulates Google Maps API responses for testing and development.

## Purpose

- **Avoid API costs** during K6 load testing
- **Fast responses** for development
- **Controllable test data**
- **No external dependencies**

## Endpoints

### Health Check
```bash
GET /health
```

### Text Search (Mock)
```bash
POST /v1/places:searchText
Content-Type: application/json

{
  "textQuery": "ramen tokyo"
}
```

### Place Details (Mock)
```bash
GET /v1/places/:placeId
```

### Photo (Mock)
```bash
GET /v1/:photoName/media
```

## Running Locally

```bash
# Install dependencies
cd cmd/mock-map-service
go mod download

# Run the service
go run main.go

# Service will start on http://localhost:8085
```

## Running with Docker

```bash
# Build
docker build -t mock-map-service -f deployments/docker/Dockerfile.mock-map-service .

# Run
docker run -p 8085:8085 mock-map-service
```

## Test Data

Mock data is loaded from `testdata/places.json`. You can customize this file to add more test restaurants.

Default mock places:
- `mock_tokyo_ramen_1` - 一蘭拉麵 (Tokyo)
- `mock_osaka_sushi_1` - すしざんまい (Osaka)
- `mock_kyoto_tempura_1` - 天ぷら京都 (Kyoto)
- `mock_fukuoka_ramen_1` - 博多一風堂 (Fukuoka)
- `mock_tokyo_sushi_1` - 築地寿司 (Tokyo)

## Usage with Map Service

Set environment variable to use mock instead of real Google API:

```bash
MOCK_MODE=true
MOCK_BASE_URL=http://localhost:8085
```

## K6 Testing

```javascript
// Use mock service in K6 tests
const MOCK_PLACE_IDS = [
  'mock_tokyo_ramen_1',
  'mock_osaka_sushi_1',
  'mock_kyoto_tempura_1',
];

export default function() {
  const placeId = MOCK_PLACE_IDS[Math.floor(Math.random() * MOCK_PLACE_IDS.length)];
  http.get(`http://localhost:18082/api/v1/restaurants/quick-search/${placeId}`);
}
```

## Features

- ✅ No API costs
- ✅ Fast response times (<10ms)
- ✅ Customizable test data
- ✅ Compatible with Google Maps API format
- ✅ Easy to extend

## License

MIT
