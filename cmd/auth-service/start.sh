#!/bin/bash

# Auth Service Quick Start Script

set -e

echo "🚀 Starting Auth Service..."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop first."
    exit 1
fi

# Navigate to auth-service directory
cd "$(dirname "$0")"

# Check if .env exists, if not copy from example
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please update JWT_SECRET in .env before production use!"
    echo ""
fi

# Start services
echo "🐳 Starting Docker containers..."
docker-compose up -d

echo ""
echo "✅ Auth Service is starting up!"
echo ""
echo "📊 Service Information:"
echo "  - HTTP API:  http://localhost:8080"
echo "  - gRPC API:  localhost:9090"
echo "  - Health:    http://localhost:8080/health"
echo ""
echo "📝 Useful Commands:"
echo "  - View logs:     docker-compose logs -f auth-service"
echo "  - Stop service:  docker-compose down"
echo "  - Restart:       docker-compose restart auth-service"
echo ""
echo "🔍 Checking service health..."
sleep 5

# Wait for service to be healthy
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Auth Service is healthy and ready!"
        echo ""
        echo "🎉 You can now use the Auth Service!"
        echo ""
        echo "📚 API Documentation:"
        echo "  - POST /api/v1/auth/register - Register new user"
        echo "  - POST /api/v1/auth/login    - Login"
        echo "  - POST /api/v1/auth/refresh  - Refresh token"
        echo "  - GET  /api/v1/auth/validate - Validate token"
        exit 0
    fi
    attempt=$((attempt + 1))
    echo "⏳ Waiting for service to be ready... ($attempt/$max_attempts)"
    sleep 2
done

echo "⚠️  Service is taking longer than expected to start."
echo "   Check logs with: docker-compose logs -f auth-service"
