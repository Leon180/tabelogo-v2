# 後端服務環境變數配置策略

## 環境變數檔案位置

### 方案 1: 每個服務獨立的 .env (推薦用於開發) ⭐

**結構**:
```
cmd/
├── auth-service/
│   ├── .env          # Auth Service 的環境變數
│   └── main.go
├── map-service/
│   ├── .env          # Map Service 的環境變數 (包含 GOOGLE_MAPS_API_KEY)
│   └── main.go
├── restaurant-service/
│   ├── .env
│   └── main.go
└── booking-service/
    ├── .env
    └── main.go
```

**優點**:
- ✅ 每個服務的配置獨立,易於管理
- ✅ 本地開發時直接運行 `go run cmd/map-service/main.go` 即可
- ✅ 符合微服務的獨立性原則

**實作方式**:

在每個服務的 `main.go` 中加載 `.env`:

```go
package main

import (
    "log"
    "github.com/joho/godotenv"
)

func main() {
    // 加載 .env 文件
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    
    // 讀取環境變數
    apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    port := os.Getenv("PORT")
    
    // ... 啟動服務
}
```

**範例 - Map Service 的 .env**:
```bash
# cmd/map-service/.env
GOOGLE_MAPS_API_KEY=AIza...你的服務器Key
PORT=8080
LOG_LEVEL=debug
```

### 方案 2: 集中式環境變數 (推薦用於 Docker)

**結構**:
```
deployments/
└── docker-compose/
    ├── .env.development    # 開發環境
    ├── .env.production     # 生產環境
    └── docker-compose.yml
```

**docker-compose.yml 範例**:
```yaml
services:
  map-service:
    build:
      context: ../../
      dockerfile: cmd/map-service/Dockerfile
    env_file:
      - .env.${ENVIRONMENT:-development}
    environment:
      - SERVICE_NAME=map-service
      - GOOGLE_MAPS_API_KEY=${GOOGLE_MAPS_API_KEY}
    ports:
      - "8080:8080"
  
  auth-service:
    build:
      context: ../../
      dockerfile: cmd/auth-service/Dockerfile
    env_file:
      - .env.${ENVIRONMENT:-development}
    environment:
      - SERVICE_NAME=auth-service
      - DB_HOST=${AUTH_DB_HOST}
    ports:
      - "8081:8081"
```

**.env.development 範例**:
```bash
# deployments/docker-compose/.env.development

# Google Maps
GOOGLE_MAPS_API_KEY=AIza...服務器Key

# Database
AUTH_DB_HOST=postgres-auth
AUTH_DB_PORT=5432
RESTAURANT_DB_HOST=postgres-restaurant
RESTAURANT_DB_PORT=5432

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
```

## 推薦配置策略

### 開發環境 (本地運行)
使用**方案 1** - 每個服務獨立的 `.env`:
```bash
# 1. 在每個 cmd/<service>/ 目錄下創建 .env
cd cmd/map-service
cat > .env << EOF
GOOGLE_MAPS_API_KEY=AIza...
PORT=8080
EOF

# 2. 直接運行服務
go run main.go
```

### Docker 環境
使用**方案 2** - 集中式配置:
```bash
# 1. 在 deployments/docker-compose/ 創建環境文件
cd deployments/docker-compose
cp .env.example .env.development

# 2. 編輯 .env.development 填入實際值

# 3. 啟動所有服務
docker-compose --env-file .env.development up
```

## 安全最佳實踐

### 1. .gitignore 設定

確保 `.gitignore` 包含:
```
# Environment files
.env
.env.local
.env.*.local
*.env
cmd/*/.env
deployments/**/.env.*
!.env.example
```

### 2. 範例文件

為每個服務提供 `.env.example`:
```bash
# cmd/map-service/.env.example
GOOGLE_MAPS_API_KEY=your_server_api_key_here
PORT=8080
LOG_LEVEL=info
```

### 3. 生產環境

生產環境建議使用:
- **Kubernetes Secrets**
- **AWS Secrets Manager**
- **HashiCorp Vault**
- **環境變數注入** (不使用 .env 文件)

## 您當前的配置

根據您打開的文件,您已經在正確的位置創建了 `.env`:

✅ `/Users/lileon/goproject/tabelogov2/cmd/map-service/.env`

這是**正確的**!只需確保:

1. **在 `cmd/map-service/main.go` 中加載 .env**:
   ```go
   import "github.com/joho/godotenv"
   
   func main() {
       godotenv.Load() // 加載當前目錄的 .env
       // ...
   }
   ```

2. **安裝 godotenv 依賴**:
   ```bash
   cd cmd/map-service
   go get github.com/joho/godotenv
   ```

3. **測試環境變數是否載入**:
   ```go
   apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
   if apiKey == "" {
       log.Fatal("GOOGLE_MAPS_API_KEY not set")
   }
   log.Printf("API Key loaded: %s...", apiKey[:10])
   ```

## 檢查清單

- [x] 前端 `.env.local` 已配置 (客戶端 API Key)
- [x] 後端 `cmd/map-service/.env` 已配置 (服務器 API Key)
- [ ] 在 `main.go` 中加載 `.env` 文件
- [ ] 測試服務器 API Key 是否正確載入
- [ ] 將 `.env` 加入 `.gitignore`
- [ ] 創建 `.env.example` 作為範本

## 下一步

1. 確認 Map Service 的 `main.go` 有載入 `.env`
2. 啟動 Map Service 測試
3. 前端調用 Map Service 的 API 進行搜尋
