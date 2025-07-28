# AI Service Setup Guide

## Getting Started

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Set up API Keys

Create a `.env` file in the root directory with your API keys:

```env
# AI Service Configuration
PORT=8080
GIN_MODE=debug
LOG_LEVEL=info

# Database
DB_PATH=./ai_service.db

# Service Info
SERVICE_NAME=ai-service
SERVICE_VERSION=1.0.0

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=60

# Default AI Provider
DEFAULT_AI_PROVIDER=openai

# API Keys - Get these from the respective providers
# OpenAI: https://platform.openai.com/api-keys
OPENAI_API_KEY=your_openai_api_key_here

# Google Gemini: https://makersuite.google.com/app/apikey
GEMINI_API_KEY=your_gemini_api_key_here

# Anthropic: https://console.anthropic.com/
ANTHROPIC_API_KEY=your_anthropic_api_key_here
```

### 3. How to Get Gemini API Key

1. Go to [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with your Google account
3. Click "Create API Key"
4. Copy the generated API key
5. Add it to your `.env` file as `GEMINI_API_KEY=your_key_here`

### 4. Run the Service

```bash
# Build and run
make run

# Or run directly
go run main.go
```

### 5. Access the UI

Once the service is running, you can access:

- **Web UI**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html
- **API Health Check**: http://localhost:8080/api/v1/health

## Features

### Web Interface
- **Home Page**: AI content generation with multiple providers
- **History Page**: View generation history
- **Stats Page**: View usage statistics

### API Endpoints
- `POST /api/v1/generate` - Generate content with a single provider
- `POST /api/v1/compare` - Compare multiple providers
- `GET /api/v1/providers` - Get available providers
- `GET /api/v1/history` - Get generation history
- `GET /api/v1/stats` - Get usage statistics
- `GET /api/v1/health` - Health check

## Troubleshooting

### Common Issues

1. **"Gemini API key not configured"**
   - Make sure you've added your Gemini API key to the `.env` file
   - Verify the key is correct and active

2. **"Provider not found"**
   - Check that you have at least one API key configured
   - Verify the provider name in your request

3. **"Invalid API key"**
   - Regenerate your API key from the provider's dashboard
   - Make sure you're using the correct key format

### Testing the API

You can test the API using curl:

```bash
# Test generation
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "gemini",
    "model": "gemini-1.5-flash",
    "prompt": "Hello, how are you?",
    "maxTokens": 100
  }'

# Test health check
curl http://localhost:8080/api/v1/health
```

## Supported Providers

- **OpenAI**: GPT-3.5, GPT-4 models
- **Google Gemini**: Gemini 1.5 Flash, Gemini 1.5 Pro
- **Anthropic**: Claude (coming soon) 