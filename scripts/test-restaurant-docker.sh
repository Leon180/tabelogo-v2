#!/bin/bash

# Test Restaurant Service with Docker Compose
# This script builds and tests the Restaurant Service in Docker

set -e  # Exit on error

echo "============================================"
echo "Restaurant Service Docker Integration Test"
echo "============================================"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Navigate to docker-compose directory
cd "$(dirname "$0")/../deployments/docker-compose"

echo -e "\n${YELLOW}Step 1: Stopping any existing containers...${NC}"
docker-compose down -v

echo -e "\n${YELLOW}Step 2: Building Restaurant Service Docker image...${NC}"
docker-compose build restaurant-service

echo -e "\n${YELLOW}Step 3: Starting infrastructure (PostgreSQL, Redis)...${NC}"
docker-compose up -d postgres-restaurant redis

echo -e "\n${YELLOW}Step 4: Waiting for PostgreSQL to be ready...${NC}"
sleep 10
docker-compose exec -T postgres-restaurant pg_isready -U postgres

echo -e "\n${YELLOW}Step 5: Starting Restaurant Service...${NC}"
docker-compose up -d restaurant-service

echo -e "\n${YELLOW}Step 6: Waiting for Restaurant Service to start...${NC}"
sleep 15

echo -e "\n${YELLOW}Step 7: Checking service health...${NC}"
docker-compose ps restaurant-service

echo -e "\n${YELLOW}Step 8: Viewing service logs...${NC}"
docker-compose logs --tail=50 restaurant-service

echo -e "\n${YELLOW}Step 9: Testing health endpoint...${NC}"
if curl -f http://localhost:18082/health 2>/dev/null; then
    echo -e "${GREEN}✓ Health check passed!${NC}"
else
    echo -e "${RED}✗ Health check failed!${NC}"
    echo -e "\n${RED}Service logs:${NC}"
    docker-compose logs restaurant-service
    exit 1
fi

echo -e "\n${YELLOW}Step 10: Testing restaurant creation endpoint...${NC}"
RESPONSE=$(curl -s -X POST http://localhost:18082/api/v1/restaurants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Sushi Restaurant",
    "source": "google",
    "external_id": "test-docker-123",
    "address": "Tokyo, Japan",
    "latitude": 35.6762,
    "longitude": 139.6503,
    "rating": 4.5,
    "price_range": "$$",
    "cuisine_type": "Japanese",
    "phone": "03-1234-5678",
    "website": "https://example.com"
  }')

if echo "$RESPONSE" | grep -q "id"; then
    echo -e "${GREEN}✓ Restaurant creation succeeded!${NC}"
    echo "Response: $RESPONSE"
    RESTAURANT_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "Restaurant ID: $RESTAURANT_ID"
else
    echo -e "${RED}✗ Restaurant creation failed!${NC}"
    echo "Response: $RESPONSE"
    exit 1
fi

echo -e "\n${YELLOW}Step 11: Testing restaurant retrieval by ID...${NC}"
if [ -n "$RESTAURANT_ID" ]; then
    RESPONSE=$(curl -s http://localhost:18082/api/v1/restaurants/$RESTAURANT_ID)
    if echo "$RESPONSE" | grep -q "Test Sushi Restaurant"; then
        echo -e "${GREEN}✓ Restaurant retrieval succeeded!${NC}"
        echo "Response: $RESPONSE"
    else
        echo -e "${RED}✗ Restaurant retrieval failed!${NC}"
        echo "Response: $RESPONSE"
    fi
fi

echo -e "\n${YELLOW}Step 12: Testing restaurant retrieval by external ID...${NC}"
RESPONSE=$(curl -s http://localhost:18082/api/v1/restaurants/external/google/test-docker-123)
if echo "$RESPONSE" | grep -q "Test Sushi Restaurant"; then
    echo -e "${GREEN}✓ External ID retrieval succeeded!${NC}"
    echo "Response: $RESPONSE"
else
    echo -e "${RED}✗ External ID retrieval failed!${NC}"
    echo "Response: $RESPONSE"
fi

echo -e "\n${YELLOW}Step 13: Testing restaurant search...${NC}"
RESPONSE=$(curl -s "http://localhost:18082/api/v1/restaurants?limit=10&offset=0")
if echo "$RESPONSE" | grep -q "Test Sushi Restaurant"; then
    echo -e "${GREEN}✓ Restaurant search succeeded!${NC}"
    echo "Response: $RESPONSE"
else
    echo -e "${RED}✗ Restaurant search failed!${NC}"
    echo "Response: $RESPONSE"
fi

echo -e "\n${GREEN}============================================${NC}"
echo -e "${GREEN}All tests passed! ✓${NC}"
echo -e "${GREEN}============================================${NC}"

echo -e "\n${YELLOW}Service is running at:${NC}"
echo "  - Health Check: http://localhost:18082/health"
echo "  - Restaurants API: http://localhost:18082/api/v1/restaurants"
echo ""
echo -e "${YELLOW}To view logs:${NC}"
echo "  docker-compose logs -f restaurant-service"
echo ""
echo -e "${YELLOW}To stop services:${NC}"
echo "  docker-compose down"
echo ""
echo -e "${YELLOW}To stop and remove volumes:${NC}"
echo "  docker-compose down -v"
