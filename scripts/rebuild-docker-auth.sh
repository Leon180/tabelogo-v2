#!/bin/bash

# Auth Service Docker Rebuild Script
# This script rebuilds and restarts the Auth Service with Swagger UI fix

set -e

echo "🔧 Auth Service Docker Rebuild & Test"
echo "======================================"
echo ""

# Navigate to project root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

echo "📍 Project root: $PROJECT_ROOT"
echo ""

# Step 1: Stop existing containers
echo "1️⃣  Stopping existing Auth Service..."
docker-compose -f deployments/docker-compose/auth-service.yml down
echo "✅ Stopped"
echo ""

# Step 2: Build new image
echo "2️⃣  Building new Auth Service image..."
docker build -f cmd/auth-service/Dockerfile -t tabelogo-auth-service:latest .
echo "✅ Built"
echo ""

# Step 3: Start services
echo "3️⃣  Starting Auth Service..."
docker-compose -f deployments/docker-compose/auth-service.yml up -d
echo "✅ Started"
echo ""

# Step 4: Wait for services to be ready
echo "4️⃣  Waiting for services to be ready (30 seconds)..."
sleep 30
echo "✅ Ready"
echo ""

# Step 5: Check container status
echo "5️⃣  Checking container status..."
docker-compose -f deployments/docker-compose/auth-service.yml ps
echo ""

# Step 6: Test endpoints
echo "6️⃣  Testing endpoints..."
echo ""

# Test Health Check
echo "   Testing Health Check..."
if curl -sf http://localhost:18080/health > /dev/null; then
    echo "   ✅ Health Check: OK"
else
    echo "   ❌ Health Check: FAILED"
fi

# Test Swagger JSON
echo "   Testing Swagger JSON..."
if curl -sf http://localhost:18080/auth-service/swagger/doc.json > /dev/null; then
    echo "   ✅ Swagger JSON: OK"
else
    echo "   ❌ Swagger JSON: FAILED"
fi

# Test Swagger UI
echo "   Testing Swagger UI..."
if curl -sf http://localhost:18080/auth-service/swagger/index.html > /dev/null; then
    echo "   ✅ Swagger UI: OK"
else
    echo "   ❌ Swagger UI: FAILED"
fi

# Test Redirect
echo "   Testing Redirect..."
if curl -sf -L http://localhost:18080/swagger > /dev/null; then
    echo "   ✅ Redirect: OK"
else
    echo "   ❌ Redirect: FAILED"
fi

echo ""
echo "======================================"
echo "🎉 Rebuild Complete!"
echo ""
echo "📚 Access Swagger UI at:"
echo "   http://localhost:18080/auth-service/swagger/index.html"
echo ""
echo "🔍 View logs:"
echo "   docker-compose -f deployments/docker-compose/auth-service.yml logs -f auth-service"
echo ""
echo "🐚 Enter container:"
echo "   docker exec -it tabelogo-auth-service sh"
echo ""
