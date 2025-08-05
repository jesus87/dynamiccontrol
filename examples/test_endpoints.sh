#!/bin/bash

# Dynamic Control Plane - Test Endpoints
# Ejemplos de curl para probar todos los endpoints

BASE_URL="http://localhost:8080"

echo "🚀 Testing Dynamic Control Plane Endpoints"
echo "=========================================="
echo ""

# 1. Health Check
echo "1. Testing Health Check..."
curl -X GET "$BASE_URL/health" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 2. Service Info
echo "2. Testing Service Info..."
curl -X GET "$BASE_URL/info" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 3. Status Endpoint
echo "3. Testing Status Endpoint..."
curl -X GET "$BASE_URL/v1/status" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 4. Traffic Endpoint - Valid Request
echo "4. Testing Traffic Endpoint - Valid Request..."
curl -X POST "$BASE_URL/v1/services/service123/traffic" \
  -H "Content-Type: application/json" \
  -d '{
    "trafficType": "incoming",
    "volume": 100.5,
    "priority": "medium",
    "metadata": {
      "source": "service-a",
      "destination": "service-b",
      "protocol": "http"
    }
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 5. Traffic Endpoint - High Priority
echo "5. Testing Traffic Endpoint - High Priority..."
curl -X POST "$BASE_URL/v1/services/service456/traffic" \
  -H "Content-Type: application/json" \
  -d '{
    "trafficType": "outgoing",
    "volume": 250.0,
    "priority": "high",
    "metadata": {
      "source": "service-c",
      "destination": "service-d",
      "protocol": "grpc"
    }
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 6. Traffic Endpoint - Critical Priority
echo "6. Testing Traffic Endpoint - Critical Priority..."
curl -X POST "$BASE_URL/v1/services/service789/traffic" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token123" \
  -d '{
    "trafficType": "internal",
    "volume": 500.0,
    "priority": "critical",
    "metadata": {
      "source": "service-e",
      "destination": "service-f",
      "protocol": "tcp"
    }
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 7. Traffic Endpoint - Invalid Request (should fail validation)
echo "7. Testing Traffic Endpoint - Invalid Request (should fail)..."
curl -X POST "$BASE_URL/v1/services/service123/traffic" \
  -H "Content-Type: application/json" \
  -d '{
    "trafficType": "invalid",
    "volume": -10,
    "priority": "wrong"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 8. Traffic Endpoint - Missing Required Fields
echo "8. Testing Traffic Endpoint - Missing Required Fields..."
curl -X POST "$BASE_URL/v1/services/service123/traffic" \
  -H "Content-Type: application/json" \
  -d '{
    "trafficType": "incoming"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 9. Non-existent Endpoint (should return 404)
echo "9. Testing Non-existent Endpoint (should return 404)..."
curl -X GET "$BASE_URL/v1/nonexistent" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

# 10. Wrong Method (should return 405)
echo "10. Testing Wrong Method (should return 405)..."
curl -X PUT "$BASE_URL/v1/status" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n"
echo ""
echo ""

echo "✅ All tests completed!"
echo ""
echo "Expected Results:"
echo "- Tests 1-6: Should return 200 OK"
echo "- Test 7: Should return 400 Bad Request (validation error)"
echo "- Test 8: Should return 400 Bad Request (missing fields)"
echo "- Test 9: Should return 404 Not Found"
echo "- Test 10: Should return 405 Method Not Allowed"