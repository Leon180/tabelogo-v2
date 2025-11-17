.PHONY: help init up down restart logs ps clean build test lint proto migrate

# Variables
DOCKER_COMPOSE = docker-compose -f deployments/docker-compose/docker-compose.yml
SERVICES = auth-service restaurant-service booking-service spider-service mail-service map-service api-gateway

## help: 顯示此幫助訊息
help:
	@echo "可用指令："
	@echo "  make init        - 初始化專案（建立 .env、安裝依賴）"
	@echo "  make up          - 啟動所有 Docker 容器"
	@echo "  make down        - 停止並移除所有容器"
	@echo "  make restart     - 重啟所有容器"
	@echo "  make logs        - 查看所有容器日誌"
	@echo "  make ps          - 查看容器狀態"
	@echo "  make clean       - 清理所有容器和 volumes"
	@echo "  make build       - 建置所有微服務"
	@echo "  make test        - 執行所有測試"
	@echo "  make lint        - 執行程式碼檢查"
	@echo "  make proto       - 生成 Protocol Buffers 程式碼"
	@echo "  make migrate     - 執行資料庫 migrations"

## init: 初始化專案
init:
	@echo "=> 初始化專案..."
	@if [ ! -f .env ]; then cp .env.example .env && echo "已建立 .env 檔案"; fi
	@echo "=> 初始化完成！"

## up: 啟動所有 Docker 容器
up:
	@echo "=> 啟動 Docker 容器..."
	$(DOCKER_COMPOSE) up -d
	@echo "=> 所有容器已啟動"
	@echo "=> Kafka UI: http://localhost:8080"
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
