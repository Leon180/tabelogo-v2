# Swagger UI 訪問指南

## 🎯 快速開始

### 本地開發環境 (VSCode)

```bash
# 1. 啟動依賴服務
docker-compose -f deployments/docker-compose/auth-service.yml up -d postgres-auth redis-auth

# 2. 在 VSCode 按 F5 啟動 Auth Service

# 3. 訪問 Swagger UI
open http://localhost:8081/auth-service/swagger/index.html
```

### Docker 環境

```bash
# 方法 1: 一鍵重建（推薦）
make auth-rebuild

# 方法 2: 手動步驟
make auth-build
make auth-up

# 等待 30 秒後訪問
open http://localhost:18080/auth-service/swagger/index.html
```

## 📊 端口對比

| 環境 | HTTP | gRPC | Swagger URL |
|------|------|------|-------------|
| **本地開發** | 8081 | 9091 | http://localhost:8081/auth-service/swagger/index.html |
| **Docker** | 18080 | 19090 | http://localhost:18080/auth-service/swagger/index.html |

## 🔧 Makefile 命令

### Swagger 文檔生成
```bash
make swagger-auth    # 生成 Swagger 文檔
```

### Docker 相關
```bash
make auth-build      # 構建 Docker 鏡像
make auth-rebuild    # 完整重建（停止→構建→啟動→測試）
make auth-up         # 啟動服務
make auth-down       # 停止服務
make auth-restart    # 重啟服務
make auth-logs       # 查看日誌
make auth-ps         # 查看狀態
make auth-clean      # 清理容器和數據
make auth-shell      # 進入容器
```

### 本地開發
```bash
make auth-dev        # 本地開發模式（自動生成 Swagger + 啟動）
```

## 🌐 可用的 Swagger 端點

### 本地開發 (Port 8081)
- Swagger UI: http://localhost:8081/auth-service/swagger/index.html
- Swagger JSON: http://localhost:8081/auth-service/swagger/doc.json
- 快捷訪問: http://localhost:8081/swagger

### Docker (Port 18080)
- Swagger UI: http://localhost:18080/auth-service/swagger/index.html
- Swagger JSON: http://localhost:18080/auth-service/swagger/doc.json
- 快捷訪問: http://localhost:18080/swagger

## 🐛 故障排除

### 本地開發問題

#### 問題: 404 Not Found
```bash
# 1. 確認 Swagger 文檔已生成
ls -la internal/auth/docs/

# 2. 重新生成
make swagger-auth

# 3. 重啟服務 (Shift+F5, 然後 F5)
```

#### 問題: 重定向循環
```bash
# 清除瀏覽器緩存
# Chrome: F12 → 右鍵刷新按鈕 → 清空緩存並硬性重新載入

# 或使用無痕模式
```

### Docker 環境問題

#### 問題: 404 Not Found
```bash
# 1. 檢查容器內的文件
docker exec -it tabelogo-auth-service sh
ls -la /app/internal/auth/docs/

# 2. 如果文件不存在，重新構建
make auth-rebuild
```

#### 問題: 服務啟動失敗
```bash
# 查看日誌
make auth-logs

# 檢查容器狀態
make auth-ps

# 清理並重新啟動
make auth-clean
make auth-rebuild
```

#### 問題: 端口已被佔用
```bash
# 查看端口使用情況
lsof -i :18080
lsof -i :19090

# 停止衝突的服務或修改 docker-compose.yml 中的端口
```

## 📁 文件結構

```
/app/ (Docker 容器)
├── auth-service          # 二進制文件
└── internal/
    └── auth/
        └── docs/
            ├── docs.go       # Swagger 元數據
            ├── index.html    # Swagger UI
            ├── swagger.json  # OpenAPI JSON
            └── swagger.yaml  # OpenAPI YAML
```

## 🔍 驗證修復

```bash
# 測試本地開發環境
curl http://localhost:8081/auth-service/swagger/index.html

# 測試 Docker 環境
curl http://localhost:18080/auth-service/swagger/index.html

# 測試重定向
curl -L http://localhost:18080/swagger
```

## 📚 相關文檔

- [VSCode 調試指南](docs/VSCODE_DEBUG_GUIDE.md)
- [快速啟動指南](docs/QUICK_START.md)
- [Docker Swagger 修復](docs/DOCKER_SWAGGER_FIX.md)
- [Swagger 故障排除](docs/SWAGGER_TROUBLESHOOTING.md)
- [最終解決方案](docs/FINAL_SOLUTION_SUMMARY.md)

## 🆘 獲取幫助

```bash
# 查看所有可用命令
make help

# 查看特定命令的說明
make help | grep auth
```

## 📝 注意事項

1. **首次運行**: 需要下載 Docker 鏡像，可能需要幾分鐘
2. **端口衝突**: Docker 使用 18080/19090，本地開發使用 8081/9091
3. **數據持久化**: 使用 `make auth-clean` 會刪除所有數據庫數據
4. **自動測試**: `make auth-rebuild` 會自動測試所有 Swagger 端點

---

**最後更新**: 2025-11-25
**版本**: 1.0.0
