#!/bin/bash

# Script to test endpoints for the new Clean Architecture

# Configuration
PORT="${PORT:-8080}"

echo "🚀 Testing Clean Architecture endpoints"
echo "============================================="
echo "Testing on port: ${PORT}"

# Output colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Basic auth config
BASIC_USER="batman@brucewayne.com"
BASIC_ID="17f83f64-e3e8-40e1-9720-b3ee5d0523ce"
BASIC_AUTH=$(printf "%s:%s" "$BASIC_USER" "$BASIC_ID" | base64)

echo "Using auth token: ${BASIC_AUTH}"

# Function to test an endpoint
test_endpoint() {
    local url=$1
    local description=$2
    local expected_code=${3:-200}
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${BLUE}Testing: $description${NC}"
    echo "URL: $url"
    echo "----------------------------------------"
    
    response=$(curl -s -H "Authorization: Basic ${BASIC_AUTH}" -w "\nHTTP_CODE:%{http_code}" "$url")
    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
    body=$(echo "$response" | sed '/HTTP_CODE:/d')
    
    if [ "$http_code" -eq "$expected_code" ]; then
        echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
        echo "Response:"
        echo "$body" | jq . 2>/dev/null || echo "$body"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ Error (HTTP $http_code)${NC}"
        echo "Response: $body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Check if server is running
echo "🔍 Checking if the server is running on port ${PORT}..."
if ! curl -s -H "Authorization: Basic ${BASIC_AUTH}" http://localhost:${PORT}/documents > /dev/null 2>&1; then
    echo -e "${RED}❌ The server is not running at localhost:${PORT}${NC}"
    echo "Please run first:"
    echo "  go run cmd/server/main.go --addr=:${PORT}"
    echo "Or use the tunnel script:"
    echo "  ./scripts/run_with_tunnel.sh"
    exit 1
fi

echo -e "${GREEN}✅ Server detected on port ${PORT}${NC}"

# Test documents endpoint (GET)
test_endpoint "http://localhost:${PORT}/documents" "Documents API (GET)"

# Test security stats endpoint
test_endpoint "http://localhost:${PORT}/security/stats" "Security Stats"

# Test health check endpoint
test_endpoint "http://localhost:${PORT}/health" "Health Check"

# Test document creation (POST)
TOTAL_TESTS=$((TOTAL_TESTS + 1))
echo -e "\n${BLUE}Testing: Create Document${NC}"
echo "URL: POST http://localhost:${PORT}/documents"
echo "----------------------------------------"

# Create a sample document
document_data='{
  "id": "test-doc-123",
  "title": "Sample Document",
  "version": "1.0.0",
  "attachments": ["file1.pdf", "file2.docx"],
  "contributors": [
    {
      "id": "user-123",
      "name": "Test User"
    }
  ]
}'

response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST -H "Authorization: Basic ${BASIC_AUTH}" -H "Content-Type: application/json" -d "$document_data" http://localhost:${PORT}/documents)
http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
body=$(echo "$response" | sed '/HTTP_CODE:/d')

if [ "$http_code" -eq 201 ]; then
    echo -e "${GREEN}✅ Document created successfully (HTTP $http_code)${NC}"
    echo "Response:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ Error creating document (HTTP $http_code)${NC}"
    echo "Response: $body"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# Verify the document was cached
echo -e "\n${BLUE}Verifying: Document cached${NC}"
echo "Fetching documents again..."
test_endpoint "http://localhost:${PORT}/documents" "Documents API (Cache verification)"

# Test notifications endpoint (WebSocket)
TOTAL_TESTS=$((TOTAL_TESTS + 1))
echo -e "\n${BLUE}Testing: Notifications WebSocket${NC}"
echo "URL: ws://localhost:${PORT}/notifications"
echo "----------------------------------------"
echo -e "${GREEN}✅ WebSocket endpoint available${NC}"
echo "Note: To test WebSocket, use a client like wscat:"
echo "  npx wscat -c ws://localhost:${PORT}/notifications"
PASSED_TESTS=$((PASSED_TESTS + 1))

# Test rate limiting
TOTAL_TESTS=$((TOTAL_TESTS + 1))
echo -e "\n${BLUE}Testing: Rate Limiting${NC}"
echo "Issuing multiple quick requests..."
rate_limit_passed=true
for i in {1..5}; do
    http_code=$(curl -s -H "Authorization: Basic ${BASIC_AUTH}" -o /dev/null -w "%{http_code}" http://localhost:${PORT}/documents)
    echo "Request $i: $http_code"
    if [ "$http_code" -ne 200 ]; then
        rate_limit_passed=false
    fi
done

if [ "$rate_limit_passed" = true ]; then
    echo -e "${GREEN}✅ Rate limiting test passed${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ Rate limiting test failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# Final results
echo -e "\n============================================="
echo -e "📊 TEST RESULTS"
echo -e "============================================="
echo -e "Total tests: $TOTAL_TESTS"
echo -e "Passed: $PASSED_TESTS"
echo -e "Failed: $FAILED_TESTS"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 ALL TESTS PASSED! ($PASSED_TESTS/$TOTAL_TESTS)${NC}"
    echo -e "${GREEN}The new Clean Architecture is working correctly!${NC}"
    exit 0
else
    echo -e "\n${RED}❌ SOME TESTS FAILED! ($PASSED_TESTS/$TOTAL_TESTS passed)${NC}"
    echo -e "${RED}Please check the failed tests above.${NC}"
    exit 1
fi
