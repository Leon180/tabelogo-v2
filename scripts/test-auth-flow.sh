#!/bin/bash

# Authentication Flow Test Script
# Tests the complete authentication flow including login, token refresh, and validation

set -e

echo "🧪 Authentication Flow Test"
echo "============================"
echo ""

AUTH_URL="http://localhost:8080"
TEST_EMAIL="test@example.com"
TEST_PASSWORD="password123"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo "1️⃣  Testing Health Check..."
HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$AUTH_URL/health")
if [ "$HEALTH_RESPONSE" -eq 200 ]; then
    echo -e "${GREEN}✓${NC} Health check passed"
else
    echo -e "${RED}✗${NC} Health check failed (HTTP $HEALTH_RESPONSE)"
    exit 1
fi
echo ""

# Test 2: CORS Preflight (OPTIONS)
echo "2️⃣  Testing CORS Preflight..."
OPTIONS_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X OPTIONS "$AUTH_URL/api/v1/auth/login" \
    -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: POST" \
    -H "Access-Control-Request-Headers: content-type")

if [ "$OPTIONS_RESPONSE" -eq 204 ]; then
    echo -e "${GREEN}✓${NC} CORS preflight passed (HTTP 204)"
else
    echo -e "${RED}✗${NC} CORS preflight failed (HTTP $OPTIONS_RESPONSE)"
    exit 1
fi
echo ""

# Test 3: Login
echo "3️⃣  Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -H "Origin: http://localhost:3000" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token // empty')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token // empty')
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.user.id // empty')
USERNAME=$(echo "$LOGIN_RESPONSE" | jq -r '.user.username // empty')

if [ -n "$ACCESS_TOKEN" ] && [ -n "$REFRESH_TOKEN" ] && [ -n "$USER_ID" ]; then
    echo -e "${GREEN}✓${NC} Login successful"
    echo "  User: $USERNAME"
    echo "  User ID: $USER_ID"
    echo "  Access Token: ${ACCESS_TOKEN:0:20}..."
    echo "  Refresh Token: ${REFRESH_TOKEN:0:20}..."
else
    echo -e "${RED}✗${NC} Login failed"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi
echo ""

# Test 4: Validate Token
echo "4️⃣  Testing Token Validation..."
VALIDATE_RESPONSE=$(curl -s -X GET "$AUTH_URL/api/v1/auth/validate" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

VALID=$(echo "$VALIDATE_RESPONSE" | jq -r '.valid // false')
VALIDATE_USER=$(echo "$VALIDATE_RESPONSE" | jq -r '.user.username // empty')

if [ "$VALID" = "true" ] && [ -n "$VALIDATE_USER" ]; then
    echo -e "${GREEN}✓${NC} Token validation passed"
    echo "  Valid: $VALID"
    echo "  User: $VALIDATE_USER"
else
    echo -e "${RED}✗${NC} Token validation failed"
    echo "Response: $VALIDATE_RESPONSE"
    exit 1
fi
echo ""

# Test 5: Token Refresh
echo "5️⃣  Testing Token Refresh..."
REFRESH_RESPONSE=$(curl -s -X POST "$AUTH_URL/api/v1/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")

NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.access_token // empty')
NEW_REFRESH_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.refresh_token // empty')

if [ -n "$NEW_ACCESS_TOKEN" ] && [ -n "$NEW_REFRESH_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} Token refresh successful"
    echo "  New Access Token: ${NEW_ACCESS_TOKEN:0:20}..."
    echo "  New Refresh Token: ${NEW_REFRESH_TOKEN:0:20}..."

    # Update tokens for next test
    ACCESS_TOKEN="$NEW_ACCESS_TOKEN"
    REFRESH_TOKEN="$NEW_REFRESH_TOKEN"
else
    echo -e "${RED}✗${NC} Token refresh failed"
    echo "Response: $REFRESH_RESPONSE"
    exit 1
fi
echo ""

# Test 6: Use New Token
echo "6️⃣  Testing New Token..."
REVALIDATE_RESPONSE=$(curl -s -X GET "$AUTH_URL/api/v1/auth/validate" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

REVALIDATE_VALID=$(echo "$REVALIDATE_RESPONSE" | jq -r '.valid // false')

if [ "$REVALIDATE_VALID" = "true" ]; then
    echo -e "${GREEN}✓${NC} New token is valid"
else
    echo -e "${RED}✗${NC} New token validation failed"
    echo "Response: $REVALIDATE_RESPONSE"
    exit 1
fi
echo ""

# Test 7: Invalid Token
echo "7️⃣  Testing Invalid Token Handling..."
INVALID_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X GET "$AUTH_URL/api/v1/auth/validate" \
    -H "Authorization: Bearer invalid_token_123")

if [ "$INVALID_RESPONSE" -eq 401 ]; then
    echo -e "${GREEN}✓${NC} Invalid token correctly rejected (HTTP 401)"
else
    echo -e "${RED}✗${NC} Invalid token not rejected (HTTP $INVALID_RESPONSE)"
    exit 1
fi
echo ""

# Test 8: CORS Headers Check
echo "8️⃣  Testing CORS Headers..."
CORS_HEADERS=$(curl -s -X POST "$AUTH_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -H "Origin: http://localhost:3000" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}" \
    -D - -o /dev/null)

if echo "$CORS_HEADERS" | grep -q "Access-Control-Allow-Origin"; then
    echo -e "${GREEN}✓${NC} CORS headers present"
    echo "$CORS_HEADERS" | grep "Access-Control" | sed 's/^/  /'
else
    echo -e "${RED}✗${NC} CORS headers missing"
    echo "$CORS_HEADERS"
    exit 1
fi
echo ""

# Summary
echo "============================"
echo -e "${GREEN}✅ All tests passed!${NC}"
echo ""
echo "📋 Test Summary:"
echo "  1. Health Check ✓"
echo "  2. CORS Preflight ✓"
echo "  3. Login ✓"
echo "  4. Token Validation ✓"
echo "  5. Token Refresh ✓"
echo "  6. New Token Usage ✓"
echo "  7. Invalid Token Handling ✓"
echo "  8. CORS Headers ✓"
echo ""
echo "🎉 Authentication system is working correctly!"
echo ""
echo "💡 Next steps:"
echo "  1. Start frontend: cd web && npm run dev"
echo "  2. Visit: http://localhost:3000/auth/login"
echo "  3. Login with: $TEST_EMAIL / $TEST_PASSWORD"
echo "  4. Check user display: Should see 'Hi, $USERNAME'"
