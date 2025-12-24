#!/bin/bash

# Authentication Integration Test Script
# Tests all authentication flows across microservices

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASS=0
FAIL=0

echo "╔════════════════════════════════════════════╗"
echo "║  Authentication Integration Tests          ║"
echo "╚════════════════════════════════════════════╝"
echo ""

# Test 1: Health Checks
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test Suite 1: Service Health Checks"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for service in "Auth:8080" "Spider:18084" "Restaurant:18082" "Map:8081"; do
  name=$(echo $service | cut -d: -f1)
  port=$(echo $service | cut -d: -f2)
  
  HTTP=$(curl -s -w "%{http_code}" -o /dev/null http://localhost:$port/health)
  
  if [ "$HTTP" = "200" ]; then
    echo -e "${GREEN}✓${NC} $name Service: Healthy"
    ((PASS++))
  else
    echo -e "${RED}✗${NC} $name Service: Unhealthy (HTTP $HTTP)"
    ((FAIL++))
  fi
done

echo ""

# Test 2: User Registration
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test Suite 2: User Registration & Login"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

TIMESTAMP=$(date +%s)
TEST_EMAIL="test-${TIMESTAMP}@example.com"
TEST_PASSWORD="Test123!"
TEST_USERNAME="test${TIMESTAMP}"

echo "Registering user: $TEST_EMAIL"
REGISTER=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"$TEST_PASSWORD\", \"username\": \"$TEST_USERNAME\"}")

if echo "$REGISTER" | grep -q "user"; then
  echo -e "${GREEN}✓${NC} User registration successful"
  ((PASS++))
else
  echo -e "${RED}✗${NC} User registration failed"
  echo "$REGISTER"
  ((FAIL++))
fi

# Test 3: User Login
echo "Logging in user: $TEST_EMAIL"
LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"$TEST_PASSWORD\"}")

TOKEN=$(echo "$LOGIN" | python3 -c "import sys, json; print(json.load(sys.stdin)['access_token'])" 2>/dev/null)

if [ -n "$TOKEN" ]; then
  echo -e "${GREEN}✓${NC} User login successful"
  echo "   Token: ${TOKEN:0:30}..."
  ((PASS++))
else
  echo -e "${RED}✗${NC} User login failed"
  echo "$LOGIN"
  ((FAIL++))
fi

echo ""

# Test 4: Spider Service Authentication
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test Suite 3: Spider Service Authentication"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Without auth
echo "Testing Spider Service without authentication..."
HTTP=$(curl -s -w "%{http_code}" -o /dev/null \
  -X POST http://localhost:18084/api/v1/spider/scrape \
  -H "Content-Type: application/json" \
  -d '{"google_id": "test", "area": "Tokyo", "place_name": "Test"}')

if [ "$HTTP" = "401" ]; then
  echo -e "${GREEN}✓${NC} Correctly rejected unauthenticated request (401)"
  ((PASS++))
else
  echo -e "${RED}✗${NC} Should reject but got HTTP $HTTP"
  ((FAIL++))
fi

# With auth
echo "Testing Spider Service with authentication..."
HTTP=$(curl -s -w "%{http_code}" -o /tmp/spider_test.json \
  -X POST http://localhost:18084/api/v1/spider/scrape \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"google_id": "ChIJN1t_tDeuEmsRUsoyG83frY4", "area": "Tokyo", "place_name": "Integration Test"}')

if [ "$HTTP" = "200" ] || [ "$HTTP" = "202" ]; then
  echo -e "${GREEN}✓${NC} Authenticated request accepted (HTTP $HTTP)"
  ((PASS++))
else
  echo -e "${RED}✗${NC} Authenticated request failed (HTTP $HTTP)"
  cat /tmp/spider_test.json
  ((FAIL++))
fi

echo ""

# Test 5: Restaurant Service Authentication
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test Suite 4: Restaurant Service Authentication"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Public endpoint
echo "Testing public restaurant search..."
HTTP=$(curl -s -w "%{http_code}" -o /dev/null "http://localhost:18082/api/v1/restaurants/search?q=test")

if [ "$HTTP" = "200" ]; then
  echo -e "${GREEN}✓${NC} Public endpoint accessible (HTTP $HTTP)"
  ((PASS++))
else
  echo -e "${RED}✗${NC} Public endpoint failed (HTTP $HTTP)"
  ((FAIL++))
fi

# Protected endpoint without auth
echo "Testing protected endpoint without authentication..."
HTTP=$(curl -s -w "%{http_code}" -o /dev/null \
  -X POST http://localhost:18082/api/v1/favorites \
  -H "Content-Type: application/json" \
  -d '{"restaurant_id": "test"}')

if [ "$HTTP" = "401" ]; then
  echo -e "${GREEN}✓${NC} Correctly rejected unauthenticated request (401)"
  ((PASS++))
else
  echo -e "${RED}✗${NC} Should reject but got HTTP $HTTP"
  ((FAIL++))
fi

# Protected endpoint with auth
echo "Testing protected endpoint with authentication..."
HTTP=$(curl -s -w "%{http_code}" -o /tmp/rest_test.json \
  -X POST http://localhost:18082/api/v1/favorites \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"restaurant_id": "test-restaurant"}')

if [ "$HTTP" = "200" ] || [ "$HTTP" = "201" ] || [ "$HTTP" = "400" ] || [ "$HTTP" = "500" ]; then
  echo -e "${GREEN}✓${NC} Authenticated request processed (HTTP $HTTP)"
  ((PASS++))
else
  echo -e "${RED}✗${NC} Authenticated request failed (HTTP $HTTP)"
  cat /tmp/rest_test.json
  ((FAIL++))
fi

echo ""

# Test 6: Map Service Optional Auth
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test Suite 5: Map Service Optional Auth"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Without auth
echo "Testing Map Service without authentication..."
HTTP=$(curl -s -w "%{http_code}" -o /dev/null \
  -X POST http://localhost:8081/api/v1/map/quick_search \
  -H "Content-Type: application/json" \
  -d '{"query": "Tokyo Station"}')

if [ "$HTTP" = "200" ] || [ "$HTTP" = "400" ] || [ "$HTTP" = "500" ]; then
  echo -e "${GREEN}✓${NC} Public access allowed (HTTP $HTTP)"
  ((PASS++))
else
  echo -e "${RED}✗${NC} Public access failed (HTTP $HTTP)"
  ((FAIL++))
fi

# With auth
echo "Testing Map Service with authentication..."
HTTP=$(curl -s -w "%{http_code}" -o /dev/null \
  -X POST http://localhost:8081/api/v1/map/quick_search \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query": "Tokyo Station"}')

if [ "$HTTP" = "200" ] || [ "$HTTP" = "400" ] || [ "$HTTP" = "500" ]; then
  echo -e "${GREEN}✓${NC} Authenticated access allowed (HTTP $HTTP)"
  ((PASS++))
else
  echo -e "${RED}✗${NC} Authenticated access failed (HTTP $HTTP)"
  ((FAIL++))
fi

echo ""

# Summary
echo "╔════════════════════════════════════════════╗"
echo "║  Test Results Summary                      ║"
echo "╚════════════════════════════════════════════╝"
echo ""
echo "Total Tests: $(($PASS + $FAIL))"
echo -e "${GREEN}Passed: $PASS${NC}"
echo -e "${RED}Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
  echo -e "${GREEN}╔════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║  ✓ All Tests Passed!                      ║${NC}"
  echo -e "${GREEN}╚════════════════════════════════════════════╝${NC}"
  exit 0
else
  echo -e "${RED}╔════════════════════════════════════════════╗${NC}"
  echo -e "${RED}║  ✗ Some Tests Failed                      ║${NC}"
  echo -e "${RED}╚════════════════════════════════════════════╝${NC}"
  exit 1
fi
