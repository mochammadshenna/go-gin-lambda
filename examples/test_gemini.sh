#!/bin/bash

echo "🧪 Testing Gemini Provider Functionality"
echo "========================================"
echo

BASE_URL="http://localhost:8080"

echo "1. Testing Gemini with valid model (gemini-1.5-flash):"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"gemini","model":"gemini-1.5-flash","prompt":"Hello, how are you?"}' | head -200
echo -e "\n"

echo "2. Testing Gemini with valid model (gemini-1.5-pro):"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"gemini","model":"gemini-1.5-pro","prompt":"What is 2+2?"}' | head -200
echo -e "\n"

echo "3. Testing Gemini with valid model (gemini-2.0-flash):"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"gemini","model":"gemini-2.0-flash","prompt":"Tell me a joke"}' | head -200
echo -e "\n"

echo "4. Testing Gemini with invalid model:"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"gemini","model":"invalid-model","prompt":"test"}' | head -200
echo -e "\n"

echo "5. Testing Gemini with empty model (should fail validation):"
curl -s -X POST "$BASE_URL/api/generate" \
  -H "Content-Type: application/json" \
  -d '{"provider":"gemini","model":"","prompt":"test"}' | head -200
echo -e "\n"

echo "6. Testing provider models endpoint:"
curl -s "$BASE_URL/api/providers" | head -200
echo -e "\n"

echo "✅ Gemini provider tests completed!"
echo "🌐 Open http://localhost:8080 in your browser to test the UI"
echo "💡 Try switching to Gemini provider and submitting a request" 