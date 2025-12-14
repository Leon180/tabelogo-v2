# API Reference

## Base URL

```
http://localhost:8083/api/v1/spider
```

---

## Endpoints

### 1. Submit Scraping Job

Submit a new scraping job for restaurant data.

**Endpoint**: `POST /scrape`

**Request Headers**:
```
Content-Type: application/json
```

**Request Body**:
```json
{
  "google_id": "string",    // Required: Google Place ID
  "area": "string",          // Required: Search area (e.g., "Tokyo")
  "place_name": "string"     // Required: Restaurant name
}
```

**Success Response (202 Accepted)**:

Job submitted successfully.

```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING"
}
```

**Success Response (200 OK - Cached)**:

Results found in cache.

```json
{
  "google_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
  "restaurants": [
    {
      "link": "https://tabelog.com/tokyo/A1234/...",
      "name": "Sushi Saito",
      "rating": 4.5,
      "rating_count": 1234,
      "bookmarks": 567,
      "phone": "03-1234-5678",
      "types": ["Sushi", "Japanese"],
      "photos": [
        "https://tblg.k-img.com/restaurant/images/..."
      ]
    }
  ],
  "total_found": 5,
  "from_cache": true,
  "cached_at": "2025-12-14T10:00:00Z"
}
```

**Error Responses**:

| Status | Description | Response |
|--------|-------------|----------|
| 400 | Invalid request | `{"error": "validation error message"}` |
| 500 | Server error | `{"error": "internal server error"}` |

**Example**:

```bash
curl -X POST http://localhost:8083/api/v1/spider/scrape \
  -H "Content-Type: application/json" \
  -d '{
    "google_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
    "area": "Tokyo",
    "place_name": "Sushi Saito"
  }'
```

---

### 2. Get Job Status

Retrieve the current status of a scraping job.

**Endpoint**: `GET /jobs/:job_id`

**Path Parameters**:
- `job_id` (string, required): UUID of the job

**Success Response (200 OK)**:

```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "google_id": "ChIJN1t_tDeuEmsRUsoyG83frY4",
  "status": "COMPLETED",
  "results": [
    {
      "link": "https://tabelog.com/tokyo/...",
      "name": "Sushi Saito",
      "rating": 4.5,
      "rating_count": 1234,
      "bookmarks": 567,
      "phone": "03-1234-5678",
      "types": ["Sushi", "Japanese"],
      "photos": ["https://..."]
    }
  ],
  "created_at": "2025-12-14T10:00:00+09:00",
  "completed_at": "2025-12-14T10:00:05+09:00"
}
```

**Error Response (404 Not Found)**:

```json
{
  "error": "job not found: invalid job ID"
}
```

**Job Statuses**:

| Status | Description |
|--------|-------------|
| `PENDING` | Job is queued, waiting for processing |
| `RUNNING` | Job is currently being processed |
| `COMPLETED` | Job finished successfully with results |
| `FAILED` | Job failed with error |

**Example**:

```bash
curl http://localhost:8083/api/v1/spider/jobs/550e8400-e29b-41d4-a716-446655440000
```

---

### 3. Stream Job Status (SSE)

Stream real-time job status updates via Server-Sent Events.

**Endpoint**: `GET /jobs/:job_id/stream`

**Path Parameters**:
- `job_id` (string, required): UUID of the job

**Response Headers**:
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

**Event Format**:

```
event: status
data: {"job_id":"...","status":"RUNNING",...}

event: status
data: {"job_id":"...","status":"COMPLETED","results":[...],...}
```

**Event Types**:

| Event | Description |
|-------|-------------|
| `status` | Job status update |
| `error` | Error occurred |

**Status Update Frequency**: Every 500ms

**Stream Termination**:
- Job reaches `COMPLETED` status
- Job reaches `FAILED` status
- Client disconnects
- Server error

**Example (JavaScript)**:

```javascript
const eventSource = new EventSource(
  'http://localhost:8083/api/v1/spider/jobs/550e8400.../stream'
);

eventSource.addEventListener('status', (event) => {
  const data = JSON.parse(event.data);
  console.log('Status:', data.status);
  
  if (data.status === 'COMPLETED') {
    console.log('Results:', data.results);
    eventSource.close();
  }
});

eventSource.addEventListener('error', (event) => {
  console.error('Error:', event);
  eventSource.close();
});
```

**Example (curl)**:

```bash
curl -N http://localhost:8083/api/v1/spider/jobs/550e8400.../stream
```

---

## Data Models

### ScrapingJob

```typescript
interface ScrapingJob {
  job_id: string;           // UUID
  google_id: string;        // Google Place ID
  status: JobStatus;        // Current status
  results?: Restaurant[];   // Results (if completed)
  error?: string;           // Error message (if failed)
  created_at: string;       // ISO 8601 timestamp
  completed_at?: string;    // ISO 8601 timestamp (if completed)
}
```

### Restaurant

```typescript
interface Restaurant {
  link: string;             // Tabelog URL
  name: string;             // Restaurant name
  rating: number;           // Rating (0.0 - 5.0)
  rating_count: number;     // Number of ratings
  bookmarks: number;        // Number of bookmarks
  phone: string;            // Phone number
  types: string[];          // Cuisine types
  photos: string[];         // Photo URLs
}
```

### JobStatus

```typescript
type JobStatus = 
  | "PENDING"    // Queued
  | "RUNNING"    // Processing
  | "COMPLETED"  // Success
  | "FAILED";    // Error
```

---

## Rate Limiting

**Limits**:
- 60 requests per minute per IP
- Burst of 10 requests allowed

**Headers**:
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59
X-RateLimit-Reset: 1702540800
```

**429 Response**:
```json
{
  "error": "rate limit exceeded",
  "retry_after": 30
}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Job doesn't exist |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |
| 503 | Service Unavailable - Circuit breaker open |

---

## Best Practices

### 1. Polling vs. Streaming

**Use Polling** when:
- Simple status checks
- Infrequent updates needed
- Client doesn't support SSE

**Use Streaming** when:
- Real-time updates required
- Long-running jobs
- Better user experience needed

### 2. Error Handling

```javascript
async function scrapeRestaurant(data) {
  try {
    const response = await fetch('/api/v1/spider/scrape', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Scrape failed:', error);
    throw error;
  }
}
```

### 3. Caching

- Check for cached results first
- `from_cache: true` indicates cached data
- Cache TTL is 24 hours
- Cache is invalidated on new scrape

---

## Examples

### Complete Workflow

```javascript
// 1. Submit job
const submitResponse = await fetch('/api/v1/spider/scrape', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    google_id: 'ChIJN1t_tDeuEmsRUsoyG83frY4',
    area: 'Tokyo',
    place_name: 'Sushi Saito'
  })
});

const { job_id } = await submitResponse.json();

// 2. Stream status updates
const eventSource = new EventSource(
  `/api/v1/spider/jobs/${job_id}/stream`
);

eventSource.addEventListener('status', (event) => {
  const job = JSON.parse(event.data);
  
  switch (job.status) {
    case 'PENDING':
      console.log('Job queued...');
      break;
    case 'RUNNING':
      console.log('Scraping in progress...');
      break;
    case 'COMPLETED':
      console.log('Results:', job.results);
      eventSource.close();
      break;
    case 'FAILED':
      console.error('Error:', job.error);
      eventSource.close();
      break;
  }
});
```

---

## Changelog

### v1.0.0 (2025-12-14)
- Initial API release
- SSE streaming support
- Redis caching
- Rate limiting
- Circuit breaker

---

For more information, see:
- [Architecture Guide](./architecture.md)
- [Testing Guide](./testing.md)
- [Deployment Guide](./deployment.md)
