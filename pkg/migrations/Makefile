# Migration Management Makefile

.PHONY: help test test-verbose test-coverage install-tools example-up example-down

# Default target
help:
	@echo "Available targets:"
	@echo "  test             - Run all tests"
	@echo "  test-verbose     - Run tests with verbose output"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  install-tools    - Install required tools"
	@echo "  tidy             - Run go mod tidy"
	@echo ""
	@echo "Migration CLI examples (requires postgres running):"
	@echo "  example-up       - Run example: migrate up"
	@echo "  example-down     - Run example: migrate down"
	@echo "  example-version  - Run example: show version"

# Install required tools
install-tools:
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run tests
test:
	go test -v -race ./...

# Run tests with verbose output
test-verbose:
	go test -v -race -count=1 ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Tidy dependencies (在 pkg 目錄執行)
tidy:
	cd .. && go mod tidy

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html

# Example: Migrate up (requires environment variables)
example-up:
	@if [ -z "$$DB_DSN" ]; then \
		echo "Error: DB_DSN environment variable is required"; \
		echo "Example: export DB_DSN='postgres://user:pass@localhost/db?sslmode=disable'"; \
		exit 1; \
	fi
	@if [ -z "$$MIGRATIONS_PATH" ]; then \
		echo "Error: MIGRATIONS_PATH environment variable is required"; \
		echo "Example: export MIGRATIONS_PATH='../../migrations/auth'"; \
		exit 1; \
	fi
	@if [ -z "$$SERVICE_NAME" ]; then \
		echo "Error: SERVICE_NAME environment variable is required"; \
		echo "Example: export SERVICE_NAME='auth'"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go \
		-dsn "$$DB_DSN" \
		-path "$$MIGRATIONS_PATH" \
		-service "$$SERVICE_NAME" \
		-command up

# Example: Migrate down
example-down:
	@if [ -z "$$DB_DSN" ] || [ -z "$$MIGRATIONS_PATH" ] || [ -z "$$SERVICE_NAME" ]; then \
		echo "Error: DB_DSN, MIGRATIONS_PATH, and SERVICE_NAME environment variables are required"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go \
		-dsn "$$DB_DSN" \
		-path "$$MIGRATIONS_PATH" \
		-service "$$SERVICE_NAME" \
		-command down

# Example: Show version
example-version:
	@if [ -z "$$DB_DSN" ] || [ -z "$$MIGRATIONS_PATH" ] || [ -z "$$SERVICE_NAME" ]; then \
		echo "Error: DB_DSN, MIGRATIONS_PATH, and SERVICE_NAME environment variables are required"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go \
		-dsn "$$DB_DSN" \
		-path "$$MIGRATIONS_PATH" \
		-service "$$SERVICE_NAME" \
		-command version

# Build migration CLI tool
build-cli:
	go build -o bin/migrate cmd/migrate/main.go
	@echo "Binary built: bin/migrate"
