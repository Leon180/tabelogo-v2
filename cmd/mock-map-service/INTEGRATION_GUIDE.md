# Mock Service Integration Guide

## 如何讓現有服務使用 Mock

有兩種方式可以讓您的服務使用 Mock Map Service：

---

## 方式 1: Docker Compose（推薦）

### Step 1: 添加 Mock Service 到 docker-compose.yml

在 `deployments/docker-compose/docker-compose.yml` 中添加：

```yaml
  # Mock Map Service (for testing)
  mock-map-service:
    build:
      context: ../..
      dockerfile: deployments/docker/Dockerfile.mock-map-service
    container_name: tabelogo-mock-map-service
    environment:
      # Fast mode (default - for K6 testing)
      MOCK_LATENCY_ENABLED: "false"
      # Realistic mode (uncomment for integration testing)
      # MOCK_LATENCY_ENABLED: "true"
      # MOCK_LATENCY_MIN_MS: "100"
      # MOCK_LATENCY_MAX_MS: "300"
    ports:
      - "8085:8085"
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8085/health"]
      interval: 10s
      timeout: 5s
      retries: 3
    networks:
      - tabelogo-network
    restart: unless-stopped
```

### Step 2: 更新 Map Service 環境變數

在 `map-service` 配置中添加 Mock 模式支持：

```yaml
  map-service:
    environment:
      # ... 現有配置 ...
      
      # Mock Mode (set to true to use mock service)
      USE_MOCK_API: "${USE_MOCK_API:-false}"
      MOCK_API_BASE_URL: "http://mock-map-service:8085"
```

### Step 3: 啟動服務

```bash
# 啟動所有服務（包括 Mock）
docker-compose up -d

# 或只啟動 Mock 服務
docker-compose up -d mock-map-service

# 驗證 Mock 服務運行
curl http://localhost:8085/health
```

### Step 4: 切換到 Mock 模式

```bash
# 方法 A: 環境變數
export USE_MOCK_API=true
docker-compose up -d map-service

# 方法 B: .env 文件
echo "USE_MOCK_API=true" >> .env
docker-compose up -d
```

---

## 方式 2: 本地運行（開發用）

### Step 1: 啟動 Mock Service

```bash
cd cmd/mock-map-service
go run main.go

# 或使用 Docker
docker build -t mock-map-service -f deployments/docker/Dockerfile.mock-map-service .
docker run -p 8085:8085 mock-map-service
```

### Step 2: 配置 Map Service

更新 Map Service 的環境變數：

```bash
export USE_MOCK_API=true
export MOCK_API_BASE_URL=http://localhost:8085
export GOOGLE_MAPS_API_KEY=dummy  # Mock 模式不需要真實 key

# 重啟 Map Service
cd cmd/map-service
go run main.go
```

---

## 驗證 Mock 是否生效

### 測試 1: 健康檢查

```bash
curl http://localhost:8085/health
```

預期輸出：
```json
{
  "status": "healthy",
  "service": "mock-map-service",
  "version": "1.0.0"
}
```

### 測試 2: 搜索測試

```bash
curl -X POST http://localhost:8085/v1/places:searchText \
  -H "Content-Type: application/json" \
  -d '{"textQuery": "ramen"}'
```

預期輸出：包含 mock 餐廳數據

### 測試 3: 通過 Restaurant Service 測試

```bash
# 獲取認證 token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234!"}' \
  | jq -r '.access_token')

# 測試 quick search（使用 mock place ID）
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:18082/api/v1/restaurants/quick-search/mock_tokyo_ramen_1
```

預期：返回 mock 餐廳數據

---

## K6 測試配置

### 使用 Mock 的 K6 測試腳本

```javascript
// tests/k6/restaurant_mock_test.js

import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    load_test: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '1m', target: 0 },
      ],
    },
  },
};

// Mock place IDs
const MOCK_PLACES = [
  'mock_tokyo_ramen_1',
  'mock_osaka_sushi_1',
  'mock_kyoto_tempura_1',
  'mock_fukuoka_ramen_1',
  'mock_tokyo_sushi_1',
];

export function setup() {
  // Get auth token
  const res = http.post('http://localhost:8080/api/v1/auth/login', 
    JSON.stringify({
      email: 'test@example.com',
      password: 'Test1234!',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  
  return { token: res.json('access_token') };
}

export default function(data) {
  const placeId = MOCK_PLACES[Math.floor(Math.random() * MOCK_PLACES.length)];
  
  const res = http.get(
    `http://localhost:18082/api/v1/restaurants/quick-search/${placeId}`,
    { headers: { 'Authorization': `Bearer ${data.token}` } }
  );
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
    'has restaurant data': (r) => r.json('restaurant') !== undefined,
  });
}
```

運行測試：
```bash
# 確保 Mock 服務運行
docker-compose up -d mock-map-service

# 運行 K6 測試（無 API 費用！）
k6 run tests/k6/restaurant_mock_test.js
```

---

## 切換模式

### 開發/測試環境（使用 Mock）

```bash
# .env
USE_MOCK_API=true
MOCK_API_BASE_URL=http://mock-map-service:8085
MOCK_LATENCY_ENABLED=false
```

### 生產環境（使用真實 API）

```bash
# .env.production
USE_MOCK_API=false
GOOGLE_MAPS_API_KEY=your-real-api-key
```

### 集成測試（Mock + 延遲）

```bash
# .env.integration
USE_MOCK_API=true
MOCK_API_BASE_URL=http://mock-map-service:8085
MOCK_LATENCY_ENABLED=true
MOCK_LATENCY_MIN_MS=100
MOCK_LATENCY_MAX_MS=300
```

---

## 故障排查

### 問題 1: Mock 服務無法連接

```bash
# 檢查服務狀態
docker-compose ps mock-map-service

# 查看日誌
docker-compose logs mock-map-service

# 測試連接
curl http://localhost:8085/health
```

### 問題 2: Map Service 仍在使用真實 API

```bash
# 檢查環境變數
docker-compose exec map-service env | grep MOCK

# 確認配置
docker-compose exec map-service cat /app/config.yaml
```

### 問題 3: 找不到 Mock 數據

```bash
# 查看可用的 Mock place IDs
curl http://localhost:8085/v1/places:searchText \
  -H "Content-Type: application/json" \
  -d '{"textQuery": ""}'
```

---

## 下一步

1. ✅ 添加 Mock Service 到 docker-compose.yml
2. ✅ 更新 Map Service 支持 Mock 模式
3. ✅ 創建 K6 測試腳本
4. ✅ 運行測試驗證

需要我幫您完成這些步驟嗎？
