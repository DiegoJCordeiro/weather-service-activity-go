#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="${1:-https://weather-service-activity-go-617034962015.us-central1.run.app}"

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  Weather Service API Tests${NC}"
echo -e "${YELLOW}  Base URL: $BASE_URL${NC}"
echo -e "${YELLOW}========================================${NC}\n"

# Test counter
PASSED=0
FAILED=0

# Function to test endpoint
test_endpoint() {
    local test_name="$1"
    local url="$2"
    local expected_status="$3"
    local expected_message="$4"

    echo -e "${YELLOW}Test:${NC} $test_name"
    echo -e "${YELLOW}URL:${NC} $url"

    response=$(curl -s -k -w "\n%{http_code}" "$url")
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    echo -e "${YELLOW}Response Status:${NC} $status_code"
    echo -e "${YELLOW}Response Body:${NC} $body"

    if [ "$status_code" = "$expected_status" ]; then
        if [ -z "$expected_message" ] || echo "$body" | grep -q "$expected_message"; then
            echo -e "${GREEN}✓ PASSED${NC}\n"
            ((PASSED++))
        else
            echo -e "${RED}✗ FAILED - Message mismatch${NC}\n"
            ((FAILED++))
        fi
    else
        echo -e "${RED}✗ FAILED - Expected status $expected_status, got $status_code${NC}\n"
        ((FAILED++))
    fi
}

# Test 1: Health Check
test_endpoint \
    "Health Check" \
    "$BASE_URL/health" \
    "200" \
    "OK"

# Test 2: Valid CEP - Av. Paulista, São Paulo
test_endpoint \
    "Valid CEP - 01310100 (Av. Paulista, SP)" \
    "$BASE_URL/weather/01310100" \
    "200" \
    "temp_C"

# Test 3: Valid CEP with hyphen
test_endpoint \
    "Valid CEP with hyphen - 01310-100" \
    "$BASE_URL/weather/01310-100" \
    "200" \
    "temp_C"

# Test 4: Valid CEP - Centro, Rio de Janeiro
test_endpoint \
    "Valid CEP - 20040020 (Centro, RJ)" \
    "$BASE_URL/weather/20040020" \
    "200" \
    "temp_C"

# Test 5: Invalid CEP format - too short
test_endpoint \
    "Invalid CEP - Too short (123)" \
    "$BASE_URL/weather/123" \
    "422" \
    "invalid zipcode"

# Test 6: Invalid CEP format - too long
test_endpoint \
    "Invalid CEP - Too long (123456789)" \
    "$BASE_URL/weather/123456789" \
    "422" \
    "invalid zipcode"

# Test 7: Invalid CEP format - with letters
test_endpoint \
    "Invalid CEP - With letters (0131010a)" \
    "$BASE_URL/weather/0131010a" \
    "422" \
    "invalid zipcode"

# Test 9: Empty CEP
test_endpoint \
    "Empty CEP" \
    "$BASE_URL/weather/" \
    "422" \
    "invalid zipcode"

# Test 10: Valid CEP - Brasília
test_endpoint \
    "Valid CEP - 70040902 (Brasília, DF)" \
    "$BASE_URL/weather/70040902" \
    "200" \
    "temp_C"

# Test 11: Cannot Find CEP
test_endpoint \
    "Invalid CEP - With letters (09980491)" \
    "$BASE_URL/weather/09980491" \
    "404" \
    "can not find zipcode"

# Summary
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  Test Summary${NC}"
echo -e "${YELLOW}========================================${NC}"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo -e "${YELLOW}Total: $((PASSED + FAILED))${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed! 🎉${NC}\n"
    exit 0
else
    echo -e "\n${RED}Some tests failed! 😞${NC}\n"
    exit 1
fi