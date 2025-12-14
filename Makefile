.PHONY: help init up down restart logs ps clean build test lint proto migrate

# Variables
# Docker configuration
DOCKER_COMPOSE_FILE := deployments/docker-compose/docker-compose.yml
DOCKER_COMPOSE := docker-compose -f $(DOCKER_COMPOSE_FILE)

# Enable BuildKit for faster builds and caching
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1
SERVICES = auth-service restaurant-service map-service spider-service

## help: Show this help message
help:
	@echo "Available commands:"
	@echo "  make init            - Initialize project (create .env, install dependencies)"
	@echo "  make up              - Start all Docker containers (Full System)"
	@echo "  make down            - Stop and remove all containers (Full System)"
	@echo "  make restart         - Restart all containers"
	@echo "  make logs            - View logs for all containers"
	@echo "  make ps              - View container status"
	@echo "  make clean           - Clean up all containers and volumes"
	@echo "  make build           - Build all microservices"
	@echo "  make test            - Run all tests"
	@echo "  make lint            - Run code linter"
	@echo "  make proto           - Generate Protocol Buffers code"
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Rollback database migrations"
	@echo "  make swagger         - Generate Swagger documentation for all services"
	@echo "  make swagger-auth    - Generate Swagger docs for Auth Service"
	@echo ""
	@echo "Auth Service Commands (Docker):"
	@echo "  make auth-build      - Build Auth Service Docker Image"
	@echo "  make auth-rebuild    - Rebuild & restart Auth Service (with tests)"
	@echo "  make auth-up         - Start Auth Service (Port 8080/50051)"
	@echo "  make auth-down       - Stop Auth Service"
	@echo "  make auth-restart    - Restart Auth Service"
	@echo "  make auth-logs       - View Auth Service logs"
	@echo "  make auth-ps         - View Auth Service status"
	@echo "  make auth-clean      - Clean Auth Service container and data"
	@echo "  make auth-shell      - Enter Auth Service container"
	@echo "  make auth-db         - Connect to Auth Service PostgreSQL"
	@echo "  make auth-redis      - Connect to Auth Service Redis"
	@echo "  make auth-build      - Build Auth Service Docker Image"
	@echo "  make auth-dev        - Run Auth Service locally with auto Swagger gen"

## init: Initialize project
init:
	@echo "=> Initializing project..."
	@if [ ! -f .env ]; then cp .env.example .env && echo ".env file created"; fi
	@echo "=> Initialization complete!"

## up: Start all Docker containers
up:
	@echo "=> Starting all microservices..."
	$(DOCKER_COMPOSE) up -d
	@echo "=> All services started"
	@echo "=> Auth Service HTTP: http://localhost:8080"
	@echo "=> Auth Service gRPC: localhost:9090"
	@echo "=> Grafana: http://localhost:3000 (admin/admin)"
	@echo "=> Prometheus: http://localhost:9090"

## down: Stop and remove all containers
down:
	@echo "=> Stopping all containers..."
	$(DOCKER_COMPOSE) down

## restart: Restart all containers
restart: down up

## logs: View logs for all containers
logs:
	$(DOCKER_COMPOSE) logs -f

## ps: View container status
ps:
	$(DOCKER_COMPOSE) ps

## clean: Clean up all containers and volumes
clean:
	@echo "=> Cleaning up all containers and volumes..."
	$(DOCKER_COMPOSE) down -v --remove-orphans
	@echo "=> Cleanup complete"

## docker-clean: Clean unused Docker resources (safe)
docker-clean:
	@echo "=> Cleaning Docker resources..."
	@docker container prune -f
	@docker image prune -f
	@docker volume prune -f
	@docker buildx prune --keep-storage 5GB -f 2>/dev/null || true
	@echo "=> Docker cleanup complete"

## docker-clean-all: Remove ALL unused Docker resources (WARNING: destructive)
docker-clean-all:
	@echo "WARNING: This will remove ALL unused Docker resources!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [ "$$REPLY" = "y" ] || [ "$$REPLY" = "Y" ]; then \
		docker system prune -a --volumes -f; \
		echo "=> Complete cleanup done"; \
	else \
		echo "=> Cleanup cancelled"; \
	fi

## build: Build all microservices
.PHONY: build
build:
	@echo "=> Building all microservices..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

## lint: Run code linter
lint:
	@echo "=> Running golangci-lint..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd cmd/$$service && golangci-lint run && cd ../..; \
	done
	@echo "=> Lint complete"

## proto: Generate Protocol Buffers code
proto:
	@echo "=> Generating protobuf code..."
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh
	@echo "=> Protobuf code generation complete"

## migrate-up: Run all database migrations (up)
migrate-up:
	@echo "=> Running database migrations..."
	@echo "Running auth DB migration..."
	@migrate -path migrations/auth -database "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable" up || true
	@echo "Running restaurant DB migration..."
	@migrate -path migrations/restaurant -database "postgres://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" up || true
	@echo "Running booking DB migration..."
	@migrate -path migrations/booking -database "postgres://postgres:postgres@localhost:5434/booking_db?sslmode=disable" up || true
	@echo "Running spider DB migration..."
	@migrate -path migrations/spider -database "postgres://postgres:postgres@localhost:5435/spider_db?sslmode=disable" up || true
	@echo "Running mail DB migration..."
	@migrate -path migrations/mail -database "postgres://postgres:postgres@localhost:5436/mail_db?sslmode=disable" up || true
	@echo "=> Migrations complete"

## migrate-down: Rollback all database migrations
migrate-down:
	@echo "=> Rolling back database migrations..."
	@migrate -path migrations/auth -database "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable" down || true
	@migrate -path migrations/restaurant -database "postgres://postgres:postgres@localhost:5433/restaurant_db?sslmode=disable" down || true
	@migrate -path migrations/booking -database "postgres://postgres:postgres@localhost:5434/booking_db?sslmode=disable" down || true
	@migrate -path migrations/spider -database "postgres://postgres:postgres@localhost:5435/spider_db?sslmode=disable" down || true
	@migrate -path migrations/mail -database "postgres://postgres:postgres@localhost:5436/mail_db?sslmode=disable" down || true
	@echo "=> Rollback complete"

## dev: Start development environment
dev: init up
	@echo "=> Development environment started!"

## test-unit: Run unit tests
test-unit:
	@echo "=> Running unit tests..."
	@go test -v -short ./internal/auth/application/...
	@echo "=> Unit tests complete"

## test-integration: Run integration tests
test-integration:
	@echo "=> Starting test environment..."
	@docker-compose -f docker-compose.test.yml up -d
	@echo "=> Waiting for services to be ready..."
	@sleep 5
	@echo "=> Running integration tests..."
	@TEST_DB_HOST=localhost TEST_DB_PORT=5433 TEST_REDIS_ADDR=localhost:6380 \
		go test -v ./tests/integration/...
	@echo "=> Stopping test environment..."
	@docker-compose -f docker-compose.test.yml down
	@echo "=> Integration tests complete"

## test-all: Run all tests
test-all: test-unit test-integration
	@echo "=> All tests complete"

##.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Mock generation
.PHONY: generate-mocks
generate-mocks:
	@echo "Generating mocks..."
	@mkdir -p internal/restaurant/mocks
	mockgen -source=internal/restaurant/domain/repository/restaurant_repository.go \
		-destination=internal/restaurant/mocks/mock_restaurant_repository.go \
		-package=mocks
	mockgen -source=internal/restaurant/domain/repository/favorite_repository.go \
		-destination=internal/restaurant/mocks/mock_favorite_repository.go \
		-package=mocks
	mockgen -source=internal/restaurant/application/map_client.go \
		-destination=internal/restaurant/mocks/mock_map_service_client.go \
		-package=mocks
	@echo "✅ Restaurant mocks generated"
	@echo "Generating spider service mocks..."
	@cd internal/spider && go generate ./domain/repositories/...
	@echo "✅ Spider mocks generated in internal/spider/testutil/mocks/"
	@echo "✅ All mocks generated"

.PHONY: generate-spider-mocks
generate-spider-mocks:
	@echo "Generating spider service mocks..."
	@cd internal/spider && go generate ./domain/repositories/...
	@echo "✅ Spider mocks generated"

.PHONY: test-with-mocks
test-with-mocks: generate-mocks
	@echo "Running tests with generated mocks..."
	cd internal/restaurant/application && go test -v -cover

.PHONY: clean-mocks
clean-mocks:
	@echo "Cleaning generated mocks..."
	rm -rf internal/restaurant/mocks
	@echo "✅ Mocks cleaned"

.PHONY: clean
clean: clean-mocks
	@echo "Cleaning build artifacts..."
	rm -f coverage.out coverage.html
	rm -rf bin/
	@echo "✅ Clean complete"

## auth-build: Build Auth Service Docker Image
auth-build:
	@echo "=> Building Auth Service Docker Image..."
	@docker build -f cmd/auth-service/Dockerfile -t tabelogo-auth-service:latest .
	@echo "=> Auth Service Image built"

## auth-rebuild: Rebuild and restart Auth Service (full rebuild with tests)
auth-rebuild:
	@echo "=> Rebuilding Auth Service..."
	@./scripts/rebuild-docker-auth.sh

## auth-up: Start Auth Service and dependencies (PostgreSQL, Redis)
auth-up:
	@echo "=> Starting Auth Service..."
	@docker-compose -f deployments/docker-compose/auth-service.yml up -d
	@echo "=> Auth Service started"
	@echo "=> HTTP API: http://localhost:8080"
	@echo "=> gRPC API: localhost:50051"
	@echo "=> View logs: make auth-logs"

## auth-down: Stop Auth Service
auth-down:
	@echo "=> Stopping Auth Service..."
	@docker-compose -f deployments/docker-compose/auth-service.yml down
	@echo "=> Auth Service stopped"

## auth-restart: Restart Auth Service
auth-restart: auth-down auth-up

## auth-logs: View Auth Service logs
auth-logs:
	@docker-compose -f deployments/docker-compose/auth-service.yml logs -f auth-service

## auth-ps: View Auth Service status
auth-ps:
	@docker-compose -f deployments/docker-compose/auth-service.yml ps

## auth-clean: Clean Auth Service container and data
auth-clean:
	@echo "=> Cleaning Auth Service..."
	@docker-compose -f deployments/docker-compose/auth-service.yml down -v
	@echo "=> Auth Service cleaned"

## auth-shell: Enter Auth Service container
auth-shell:
	@docker exec -it tabelogo-auth-service sh

## auth-db: Connect to Auth Service PostgreSQL
auth-db:
	@docker exec -it tabelogo-postgres-auth-dev psql -U postgres -d auth_db

## auth-redis: Connect to Auth Service Redis
auth-redis:
	@docker exec -it tabelogo-redis-auth-dev redis-cli

## swagger: Generate Swagger documentation for all services
swagger: swagger-auth
	@echo "=> All Swagger documentation generated"

## swagger-auth: Generate Swagger documentation for Auth Service
swagger-auth:
	@echo "=> Generating Swagger docs for Auth Service..."
	swag init --generalInfo cmd/auth-service/main.go --output internal/auth/docs --parseDependency --parseInternal
	@echo "=> Auth Service Swagger docs generated at internal/auth/docs/"
	@echo "=> Local Dev: http://localhost:8081/auth-service/swagger/index.html"
	@echo "=> Docker: http://localhost:8080/auth-service/swagger/index.html"

## auth-dev: Run Auth Service locally with auto Swagger generation
auth-dev:
	@echo "=> Starting Auth Service in development mode..."
	@./scripts/start-auth-service.sh
