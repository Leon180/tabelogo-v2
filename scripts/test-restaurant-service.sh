#!/bin/bash

# Restaurant Service Integration Test Script
# Tests basic CRUD operations and external ID deduplication

set -e

BASE_URL="${BASE_URL:-http://localhost:18082}"
API_VERSION="v1"
API_BASE="$BASE_URL/api/$API_VERSION"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
print_test() {
    echo -e "\n${YELLOW}TEST: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

print_summary() {
    echo -e "\n=========================================="
    echo -e "Test Summary:"
    echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    echo -e "==========================================\n"
}

# Check if service is running
print_test "Health Check"
if curl -s "$BASE_URL/health" | grep -q "healthy"; then
    print_success "Service is healthy"
else
    print_error "Service is not healthy or not running"
    exit 1
fi

# Test 1: Create Restaurant from Google Maps
print_test "Create Restaurant (Google Maps Source)"
GOOGLE_RESPONSE=$(curl -s -X POST "$API_BASE/restaurants" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sushi Dai",
    "source": "google",
    "external_id": "ChIJTest123GooglePlaceID",
    "address": "Tokyo, Chuo-ku",
    "latitude": 35.6654,
    "longitude": 139.7707,
    "rating": 4.5,
    "price_range": "$$",
    "cuisine_type": "Sushi",
    "phone": "03-1234-5678",
    "website": "https://example.com"
  }')

GOOGLE_RESTAURANT_ID=$(echo "$GOOGLE_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -n "$GOOGLE_RESTAURANT_ID" ]; then
    print_success "Created restaurant from Google Maps (ID: $GOOGLE_RESTAURANT_ID)"
else
    print_error "Failed to create restaurant from Google Maps"
    echo "Response: $GOOGLE_RESPONSE"
fi

# Test 2: Try to create duplicate (should fail with 409)
print_test "Prevent Duplicate Google Maps Restaurant"
DUPLICATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_BASE/restaurants" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sushi Dai Duplicate",
    "source": "google",
    "external_id": "ChIJTest123GooglePlaceID",
    "address": "Different Address",
    "latitude": 35.6654,
    "longitude": 139.7707
  }')

HTTP_CODE=$(echo "$DUPLICATE_RESPONSE" | tail -n1)
if [ "$HTTP_CODE" = "409" ]; then
    print_success "Correctly prevented duplicate Google Maps restaurant (409 Conflict)"
else
    print_error "Should have returned 409 Conflict, got: $HTTP_CODE"
fi

# Test 3: Create same restaurant from Tabelog (should succeed)
print_test "Create Same Restaurant from Tabelog (Different Source)"
TABELOG_RESPONSE=$(curl -s -X POST "$API_BASE/restaurants" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "鮨 大",
    "source": "tabelog",
    "external_id": "https://tabelog.com/tokyo/A1301/A130101/13012345/",
    "address": "東京都中央区築地",
    "latitude": 35.6654,
    "longitude": 139.7707,
    "rating": 3.85,
    "cuisine_type": "寿司"
  }')

TABELOG_RESTAURANT_ID=$(echo "$TABELOG_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -n "$TABELOG_RESTAURANT_ID" ]; then
    print_success "Created same restaurant from Tabelog with different source (ID: $TABELOG_RESTAURANT_ID)"
else
    print_error "Failed to create restaurant from Tabelog"
    echo "Response: $TABELOG_RESPONSE"
fi

# Test 4: Get Restaurant by ID
print_test "Get Restaurant by ID"
if [ -n "$GOOGLE_RESTAURANT_ID" ]; then
    GET_RESPONSE=$(curl -s "$API_BASE/restaurants/$GOOGLE_RESTAURANT_ID")
    if echo "$GET_RESPONSE" | grep -q "Sushi Dai"; then
        print_success "Successfully retrieved restaurant by ID"
    else
        print_error "Failed to retrieve restaurant by ID"
        echo "Response: $GET_RESPONSE"
    fi
fi

# Test 5: Search Restaurants
print_test "Search Restaurants"
SEARCH_RESPONSE=$(curl -s "$API_BASE/restaurants/search?q=Sushi&limit=10")
RESULT_COUNT=$(echo "$SEARCH_RESPONSE" | grep -o '"total":[0-9]*' | cut -d':' -f2)

if [ "$RESULT_COUNT" -ge 1 ]; then
    print_success "Search found $RESULT_COUNT restaurant(s)"
else
    print_error "Search should have found at least 1 restaurant"
    echo "Response: $SEARCH_RESPONSE"
fi

# Test 6: Add to Favorites
print_test "Add Restaurant to Favorites"
USER_ID="550e8400-e29b-41d4-a716-446655440000"  # Test UUID

if [ -n "$GOOGLE_RESTAURANT_ID" ]; then
    FAVORITE_RESPONSE=$(curl -s -X POST "$API_BASE/favorites" \
      -H "Content-Type: application/json" \
      -d "{
        \"user_id\": \"$USER_ID\",
        \"restaurant_id\": \"$GOOGLE_RESTAURANT_ID\"
      }")

    FAVORITE_ID=$(echo "$FAVORITE_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

    if [ -n "$FAVORITE_ID" ]; then
        print_success "Added restaurant to favorites (Favorite ID: $FAVORITE_ID)"
    else
        print_error "Failed to add restaurant to favorites"
        echo "Response: $FAVORITE_RESPONSE"
    fi
fi

# Test 7: Get User Favorites
print_test "Get User Favorites"
FAVORITES_RESPONSE=$(curl -s "$API_BASE/users/$USER_ID/favorites")
FAVORITES_COUNT=$(echo "$FAVORITES_RESPONSE" | grep -o '"total":[0-9]*' | cut -d':' -f2)

if [ "$FAVORITES_COUNT" -ge 1 ]; then
    print_success "User has $FAVORITES_COUNT favorite(s)"
else
    print_error "User should have at least 1 favorite"
    echo "Response: $FAVORITES_RESPONSE"
fi

# Test 8: Try to add duplicate favorite (should fail with 409)
print_test "Prevent Duplicate Favorite"
if [ -n "$GOOGLE_RESTAURANT_ID" ]; then
    DUPLICATE_FAV_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_BASE/favorites" \
      -H "Content-Type: application/json" \
      -d "{
        \"user_id\": \"$USER_ID\",
        \"restaurant_id\": \"$GOOGLE_RESTAURANT_ID\"
      }")

    FAV_HTTP_CODE=$(echo "$DUPLICATE_FAV_RESPONSE" | tail -n1)
    if [ "$FAV_HTTP_CODE" = "409" ]; then
        print_success "Correctly prevented duplicate favorite (409 Conflict)"
    else
        print_error "Should have returned 409 Conflict for duplicate favorite, got: $FAV_HTTP_CODE"
    fi
fi

# Print final summary
print_summary

# Exit with appropriate code
if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
