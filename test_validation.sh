#!/bin/bash

echo "🧪 Testing TanyAI Validation and Error Handling"
echo "================================================"
echo

BASE_URL="http://localhost:8080"

echo "1. Testing missing provider:"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"","model":"gpt-3.5-turbo","prompt":"test"}' | jq .
echo

echo "2. Testing missing model:"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"openai","model":"","prompt":"test"}' | jq .
echo

echo "3. Testing missing prompt:"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"openai","model":"gpt-3.5-turbo","prompt":""}' | jq .
echo

echo "4. Testing invalid model for provider:"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"openai","model":"invalid-model","prompt":"test"}' | jq .
echo

echo "5. Testing unsupported provider:"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"unsupported","model":"gpt-3.5-turbo","prompt":"test"}' | jq .
echo

echo "6. Testing valid request (will fail due to API quota, but validation passes):"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"openai","model":"gpt-3.5-turbo","prompt":"Hello, this is a test"}' | jq .
echo

echo "✅ Validation tests completed!"
echo "🌐 Open http://localhost:8080 in your browser to test the modern UI with modal error handling" 