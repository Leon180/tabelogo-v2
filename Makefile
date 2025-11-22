.PHONY: help init up down restart logs ps clean build test lint proto migrate

# Variables
DOCKER_COMPOSE = docker-compose -f deployments/docker-compose/docker-compose.yml
SERVICES = auth-service restaurant-service booking-service spider-service mail-service map-service api-gateway

## help: 顯示此幫助訊息
help:
	@echo "可用指令："
	@echo "  make init            - 初始化專案（建立 .env、安裝依賴）"
	@echo "  make up              - 啟動所有 Docker 容器 (完整系統)"
	@echo "  make down            - 停止並移除所有容器 (完整系統)"
	@echo "  make restart         - 重啟所有容器"
	@echo "  make logs            - 查看所有容器日誌"
	@echo "  make ps              - 查看容器狀態"
	@echo "  make clean           - 清理所有容器和 volumes"
	@echo "  make build           - 建置所有微服務"
	@echo "  make test            - 執行所有測試"
	@echo "  make lint            - 執行程式碼檢查"
	@echo "  make proto           - 生成 Protocol Buffers 程式碼"
	@echo "  make migrate-up      - 執行資料庫 migrations"
	@echo "  make migrate-down    - 回滾資料庫 migrations"
	@echo ""
	@echo "Auth Service 指令 (本地開發):"
	@echo "  make auth-up         - 啟動 Auth Service (Port 18080/19090)"
	@echo "  make auth-down       - 停止 Auth Service"
	@echo "  make auth-restart    - 重啟 Auth Service"
	@echo "  make auth-logs       - 查看 Auth Service 日誌"
	@echo "  make auth-ps         - 查看 Auth Service 狀態"
	@echo "  make auth-clean      - 清理 Auth Service 容器和資料"
	@echo "  make auth-shell      - 進入 Auth Service 容器"
	@echo "  make auth-db         - 連接到 Auth Service PostgreSQL"
	@echo "  make auth-redis      - 連接到 Auth Service Redis"
	@echo "  make auth-build      - 建置 Auth Service Docker Image"

## init: 初始化專案
init:
	@echo "=> 初始化專案..."
	@if [ ! -f .env ]; then cp .env.example .env && echo "已建立 .env 檔案"; fi
	@echo "=> 初始化完成！"

## up: 啟動所有 Docker 容器
up:
	@echo "=> 啟動所有微服務..."
	$(DOCKER_COMPOSE) up -d
	@echo "=> 所有服務已啟動"
	@echo "=> Auth Service HTTP: http://localhost:8080"
	@echo "=> Auth Service gRPC: localhost:9090"
	@echo "=> Grafana: http://localhost:3000 (admin/admin)"
	@echo "=> Prometheus: http://localhost:9090"

## down: 停止並移除所有容器
down:
	@echo "=> 停止所有容器..."
	$(DOCKER_COMPOSE) down

## restart: 重啟所有容器
restart: down up

## logs: 查看所有容器日誌
logs:
	$(DOCKER_COMPOSE) logs -f

## ps: 查看容器狀態
ps:
	$(DOCKER_COMPOSE) ps

## clean: 清理所有容器和 volumes
clean:
	@echo "=> 清理所有容器和 volumes..."
	$(DOCKER_COMPOSE) down -v --remove-orphans
	@echo "=> 清理完成"

## build: 建置所有微服務
build:
	@echo "=> 建置所有微服務..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd cmd/$$service && go build -o ../../bin/$$service . && cd ../..; \
	done
	@echo "=> 建置完成"

## test: 執行所有測試
test:
	@echo "=> 執行測試..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd cmd/$$service && go test ./... -v && cd ../..; \
	done
	@echo "=> 測試完成"

## lint: 執行程式碼檢查
lint:
	@echo "=> 執行 golangci-lint..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd cmd/$$service && golangci-lint run && cd ../..; \
	done
	@echo "=> Lint 完成"

## proto: 生成 Protocol Buffers 程式碼
proto:
	@echo "=> 生成 protobuf 程式碼..."
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh
	@echo "=> Protobuf 程式碼生成完成"

## migrate-up: 執行所有資料庫 migrations (up)
migrate-up:
	@echo "=> 執行資料庫 migrations..."
	@echo "執行 auth DB migration..."
	@migrate -path migrations/auth -database "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable" up || true
	@echo "執行 restaurant DB migration..."
	@migrate -path migrations/restaurant -database "postgres://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" up || true
	@echo "執行 booking DB migration..."
	@migrate -path migrations/booking -database "postgres://postgres:postgres@localhost:5434/booking_db?sslmode=disable" up || true
	@echo "執行 spider DB migration..."
	@migrate -path migrations/spider -database "postgres://postgres:postgres@localhost:5435/spider_db?sslmode=disable" up || true
	@echo "執行 mail DB migration..."
	@migrate -path migrations/mail -database "postgres://postgres:postgres@localhost:5436/mail_db?sslmode=disable" up || true
	@echo "=> Migrations 完成"

## migrate-down: 回滾所有資料庫 migrations
migrate-down:
	@echo "=> 回滾資料庫 migrations..."
	@migrate -path migrations/auth -database "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable" down || true
	@migrate -path migrations/restaurant -database "postgres://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" down || true
	@migrate -path migrations/booking -database "postgres://postgres:postgres@localhost:5434/booking_db?sslmode=disable" down || true
	@migrate -path migrations/spider -database "postgres://postgres:postgres@localhost:5435/spider_db?sslmode=disable" down || true
	@migrate -path migrations/mail -database "postgres://postgres:postgres@localhost:5436/mail_db?sslmode=disable" down || true
	@echo "=> Rollback 完成"

## dev: 啟動開發環境
dev: init up
	@echo "=> 開發環境已啟動！"

## test-unit: 執行單元測試
test-unit:
	@echo "=> 執行單元測試..."
	@go test -v -short ./internal/auth/application/...
	@echo "=> 單元測試完成"

## test-integration: 執行整合測試
test-integration:
	@echo "=> 啟動測試環境..."
	@docker-compose -f docker-compose.test.yml up -d
	@echo "=> 等待服務就緒..."
	@sleep 5
	@echo "=> 執行整合測試..."
	@TEST_DB_HOST=localhost TEST_DB_PORT=5433 TEST_REDIS_ADDR=localhost:6380 \
		go test -v ./tests/integration/...
	@echo "=> 關閉測試環境..."
	@docker-compose -f docker-compose.test.yml down
	@echo "=> 整合測試完成"

## test-all: 執行所有測試
test-all: test-unit test-integration
	@echo "=> 所有測試完成"

## test-coverage: 執行測試並生成覆蓋率報告
test-coverage:
	@echo "=> 生成測試覆蓋率報告..."
	@go test -v -coverprofile=coverage.out ./internal/auth/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "=> 覆蓋率報告已生成: coverage.html"

## auth-build: 建置 Auth Service Docker Image
auth-build:
	@echo "=> 建置 Auth Service Docker Image..."
	@docker build -f cmd/auth-service/Dockerfile -t tabelogo-auth-service:latest .
	@echo "=> Auth Service Image 建置完成"

## auth-up: 啟動 Auth Service 及其依賴 (PostgreSQL, Redis)
auth-up:
	@echo "=> 啟動 Auth Service..."
	@docker-compose -f deployments/docker-compose/auth-service.yml up -d
	@echo "=> Auth Service 已啟動"
	@echo "=> HTTP API: http://localhost:18080"
	@echo "=> gRPC API: localhost:19090"
	@echo "=> 查看日誌: make auth-logs"

## auth-down: 停止 Auth Service
auth-down:
	@echo "=> 停止 Auth Service..."
	@docker-compose -f deployments/docker-compose/auth-service.yml down
	@echo "=> Auth Service 已停止"

## auth-restart: 重啟 Auth Service
auth-restart: auth-down auth-up

## auth-logs: 查看 Auth Service 日誌
auth-logs:
	@docker-compose -f deployments/docker-compose/auth-service.yml logs -f auth-service

## auth-ps: 查看 Auth Service 狀態
auth-ps:
	@docker-compose -f deployments/docker-compose/auth-service.yml ps

## auth-clean: 清理 Auth Service 容器和資料
auth-clean:
	@echo "=> 清理 Auth Service..."
	@docker-compose -f deployments/docker-compose/auth-service.yml down -v
	@echo "=> Auth Service 已清理"

## auth-shell: 進入 Auth Service 容器
auth-shell:
	@docker exec -it tabelogo-auth-service sh

## auth-db: 連接到 Auth Service PostgreSQL
auth-db:
	@docker exec -it tabelogo-postgres-auth-dev psql -U postgres -d auth_db

## auth-redis: 連接到 Auth Service Redis
auth-redis:
	@docker exec -it tabelogo-redis-auth-dev redis-cli
