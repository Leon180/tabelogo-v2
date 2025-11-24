#!/bin/bash

# Swagger Auto-Generation Script for Auth Service
# This script generates Swagger documentation before starting the service

set -e

echo "🔄 Generating Swagger documentation..."

# Navigate to project root (parent of scripts directory)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

# Generate Swagger docs
swag init \
  --generalInfo cmd/auth-service/main.go \
  --output internal/auth/docs \
  --parseDependency \
  --parseInternal

echo "✅ Swagger documentation generated successfully!"
echo "📚 Docs location: internal/auth/docs/"
echo ""
echo "Starting Auth Service..."
echo ""

# Start the service
go run cmd/auth-service/main.go
