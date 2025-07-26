#!/bin/bash

# AI Service API Examples
# Make sure the service is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "=== AI Service API Examples ==="
echo ""

echo "1. Health Check"
echo "curl $BASE_URL/health"
curl -s "$BASE_URL/health" | jq .
echo ""

echo "2. Get Provider Information"
echo "curl $BASE_URL/providers"
curl -s "$BASE_URL/providers" | jq .
echo ""

echo "3. Get AI Commands Examples"
echo "curl $BASE_URL/commands"
curl -s "$BASE_URL/commands" | jq .
echo ""

echo "4. Generate Content with OpenAI"
echo "curl -X POST $BASE_URL/generate ..."
curl -s -X POST "$BASE_URL/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "model": "gpt-3.5-turbo",
    "prompt": "Write a simple hello world function in Go",
    "system_message": "You are an expert Go developer. Write clean, well-commented code.",
    "max_tokens": 500,
    "temperature": 0.3
  }' | jq .
echo ""

echo "5. Generate Content with Gemini"
echo "curl -X POST $BASE_URL/generate ..."
curl -s -X POST "$BASE_URL/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "gemini",
    "prompt": "Explain the difference between concurrency and parallelism in simple terms",
    "max_tokens": 300,
    "temperature": 0.5
  }' | jq .
echo ""

echo "6. Compare Providers"
echo "curl -X POST $BASE_URL/compare ..."
curl -s -X POST "$BASE_URL/compare" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "What are the benefits of using Go for backend development?",
    "providers": ["openai", "gemini"],
    "max_tokens": 400,
    "temperature": 0.4
  }' | jq .
echo ""

echo "7. Get Generation History"
echo "curl $BASE_URL/history?limit=5"
curl -s "$BASE_URL/history?limit=5" | jq .
echo ""

echo "8. Get Statistics"
echo "curl $BASE_URL/stats"
curl -s "$BASE_URL/stats" | jq .
echo ""

echo "=== Advanced Examples ==="
echo ""

echo "9. Code Generation with Specific Instructions"
curl -s -X POST "$BASE_URL/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "model": "gpt-4",
    "prompt": "Create a middleware function in Go for rate limiting HTTP requests using a token bucket algorithm",
    "system_message": "You are a senior backend engineer. Write production-ready Go code with proper error handling, documentation, and tests.",
    "max_tokens": 1500,
    "temperature": 0.2
  }' | jq .
echo ""

echo "10. Data Analysis Request"
curl -s -X POST "$BASE_URL/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "gemini",
    "prompt": "Analyze this JSON data and provide insights: {\"sales\": [100, 150, 120, 200, 180], \"months\": [\"Jan\", \"Feb\", \"Mar\", \"Apr\", \"May\"]}",
    "system_message": "You are a data analyst. Provide clear insights with trends and recommendations.",
    "max_tokens": 600,
    "temperature": 0.3
  }' | jq .
echo ""

echo "=== Error Handling Examples ==="
echo ""

echo "11. Invalid Provider"
curl -s -X POST "$BASE_URL/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "invalid_provider",
    "prompt": "Test prompt"
  }' | jq .
echo ""

echo "12. Missing Required Fields"
curl -s -X POST "$BASE_URL/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai"
  }' | jq .
echo ""

echo "=== Done ==="