#!/bin/bash

# Build Restaurant Service Docker Image
# This script builds the Restaurant Service Docker image for testing

set -e  # Exit on error

echo "============================================"
echo "Building Restaurant Service Docker Image"
echo "============================================"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Navigate to project root
cd "$(dirname "$0")/.."

echo -e "\n${YELLOW}Step 1: Checking Dockerfile...${NC}"
if [ -f "cmd/restaurant-service/Dockerfile" ]; then
    echo -e "${GREEN}✓ Dockerfile found${NC}"
else
    echo -e "${RED}✗ Dockerfile not found!${NC}"
    exit 1
fi

echo -e "\n${YELLOW}Step 2: Checking main.go...${NC}"
if [ -f "cmd/restaurant-service/main.go" ]; then
    echo -e "${GREEN}✓ main.go found${NC}"
else
    echo -e "${RED}✗ main.go not found!${NC}"
    exit 1
fi

echo -e "\n${YELLOW}Step 3: Checking go.mod...${NC}"
if [ -f "cmd/restaurant-service/go.mod" ]; then
    echo -e "${GREEN}✓ go.mod found${NC}"
else
    echo -e "${RED}✗ go.mod not found!${NC}"
    exit 1
fi

echo -e "\n${YELLOW}Step 4: Building Docker image...${NC}"
docker build -f cmd/restaurant-service/Dockerfile -t tabelogo-restaurant-service:test .

echo -e "\n${GREEN}============================================${NC}"
echo -e "${GREEN}Build completed successfully! ✓${NC}"
echo -e "${GREEN}============================================${NC}"

echo -e "\n${YELLOW}Image details:${NC}"
docker images tabelogo-restaurant-service:test

echo -e "\n${YELLOW}To run the image:${NC}"
echo "  docker run -p 18082:18082 --env-file cmd/restaurant-service/.env tabelogo-restaurant-service:test"

echo -e "\n${YELLOW}To use with docker-compose:${NC}"
echo "  cd deployments/docker-compose"
echo "  docker-compose up restaurant-service"
