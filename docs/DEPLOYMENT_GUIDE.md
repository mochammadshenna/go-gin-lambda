# AI Service Deployment Guide

## 🚀 Quick Start

### 1. Local Development Setup

```bash
# Clone and setup
git clone <repository>
cd ai-service
make setup

# Configure API keys
nano .env  # Add your API keys

# Run the service
make run
```

Service will be available at: `http://localhost:8080`
API Documentation: `http://localhost:8080/swagger/index.html`

### 2. Docker Deployment

```bash
# Build and run with Docker
make docker-build
make docker-run

# Or use docker-compose
docker-compose up -d
```

### 3. Production Deployment

```bash
# Build for production
make build-linux

# Copy binary and .env to server
# Configure environment variables
# Run with systemd or supervisor
```

## 🛠️ Available Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make setup` | Initial development setup |
| `make run` | Run service locally |
| `make build` | Build binary |
| `make test` | Run tests |
| `make docker-build` | Build Docker image |
| `make test-api` | Test API endpoints |

## 📋 API Endpoints

### Core Endpoints

| Method | Endpoint | Description | Example |
|--------|----------|-------------|---------|
| `GET` | `/` | Service information | `curl localhost:8080/` |
| `GET` | `/api/v1/health` | Health check | `curl localhost:8080/api/v1/health` |
| `GET` | `/api/v1/providers` | Provider comparison | `curl localhost:8080/api/v1/providers` |
| `POST` | `/api/v1/generate` | Generate content | See examples below |
| `POST` | `/api/v1/compare` | Compare providers | See examples below |
| `GET` | `/api/v1/commands` | Command examples | `curl localhost:8080/api/v1/commands` |
| `GET` | `/api/v1/history` | Generation history | `curl localhost:8080/api/v1/history` |
| `GET` | `/api/v1/stats` | Usage statistics | `curl localhost:8080/api/v1/stats` |

### Example API Calls

#### 1. Generate Code with OpenAI

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "model": "gpt-3.5-turbo",
    "prompt": "Write a REST API endpoint in Go for user registration",
    "system_message": "You are an expert Go developer",
    "max_tokens": 1000,
    "temperature": 0.3
  }'
```

#### 2. Compare Multiple Providers

```bash
curl -X POST http://localhost:8080/api/v1/compare \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Explain microservices architecture",
    "providers": ["openai", "gemini"],
    "max_tokens": 500,
    "temperature": 0.5
  }'
```

#### 3. Generate with Specific Instructions

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "gemini",
    "prompt": "Create a Docker Compose file for a web application with Redis and PostgreSQL",
    "system_message": "You are a DevOps engineer. Provide production-ready configurations.",
    "max_tokens": 800,
    "temperature": 0.2
  }'
```

## 🔧 Configuration Options

### Environment Variables

```env
# Service Configuration
PORT=8080                              # Server port
GIN_MODE=debug                         # Gin mode (debug/release)
LOG_LEVEL=info                         # Logging level

# AI Provider API Keys
OPENAI_API_KEY=your_key_here          # OpenAI API key
GEMINI_API_KEY=your_key_here          # Google Gemini API key
ANTHROPIC_API_KEY=your_key_here       # Anthropic API key

# Database
DB_HOST=localhost                     # PostgreSQL host
DB_PORT=5432                          # PostgreSQL port
DB_USER=postgres                      # PostgreSQL user
DB_PASSWORD=your_password             # PostgreSQL password
DB_NAME=ai_service                    # PostgreSQL database name
DB_SSLMODE=disable                    # SSL mode

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=60     # Rate limit per IP

# Default Provider
DEFAULT_AI_PROVIDER=openai            # Default AI provider
```

### Provider-Specific Configuration

#### OpenAI

- **Models**: `gpt-3.5-turbo`, `gpt-4`, `gpt-4o`, `gpt-4o-mini`
- **Max Tokens**: Up to 4096 (varies by model)
- **Temperature**: 0.0 - 2.0
- **Cost**: ~$0.002/1K tokens (GPT-3.5)

#### Google Gemini

- **Models**: `gemini-1.5-flash`, `gemini-1.5-pro`, `gemini-1.0-pro`
- **Max Tokens**: Up to 8192
- **Temperature**: 0.0 - 1.0
- **Cost**: Free tier available

#### Anthropic Claude

- **Models**: `claude-3-opus`, `claude-3-sonnet`, `claude-3-haiku`
- **Max Tokens**: Up to 8192
- **Temperature**: 0.0 - 1.0
- **Cost**: Premium pricing

## 📊 Monitoring & Metrics

### Health Checks

```bash
# Basic health check
curl http://localhost:8080/api/v1/health

# Expected response
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z",
  "providers": {
    "openai": "available",
    "gemini": "available",
    "anthropic": "not_configured"
  }
}
```

### Usage Statistics

```bash
# Get usage stats
curl http://localhost:8080/api/v1/stats

# Response includes:
# - Total generations
# - Generations by provider
# - Average tokens used
# - Average response time
```

### Generation History

```bash
# Get recent history
curl http://localhost:8080/api/v1/history?limit=10

# Filter by provider
curl http://localhost:8080/api/v1/history?provider=openai&limit=5
```

## 🔒 Security Best Practices

### 1. API Key Management

- Store API keys in environment variables
- Use separate keys for development/production
- Rotate keys regularly
- Monitor API key usage

### 2. Rate Limiting

- Default: 60 requests per minute per IP
- Adjust based on your needs
- Monitor for abuse patterns

### 3. CORS Configuration

- Current setting allows all origins (`*`)
- Restrict to your domains in production:

```go
c.Header("Access-Control-Allow-Origin", "https://yourdomain.com")
```

### 4. Input Validation

- All requests are validated
- Sanitize user inputs
- Implement content filtering if needed

## 🚀 Production Deployment

### 1. System Requirements

- Go 1.21+ (for building)
- Linux/macOS/Windows
- 512MB RAM minimum
- PostgreSQL database storage

### 2. Systemd Service (Linux)

Create `/etc/systemd/system/ai-service.service`:

```ini
[Unit]
Description=AI Service
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/opt/ai-service
ExecStart=/opt/ai-service/ai-service
Restart=always
RestartSec=10
Environment=PORT=8080
Environment=GIN_MODE=release
EnvironmentFile=/opt/ai-service/.env

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable ai-service
sudo systemctl start ai-service
sudo systemctl status ai-service
```

### 3. Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 4. Docker Production Setup

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  ai-service:
    build: .
    restart: unless-stopped
    environment:
      - GIN_MODE=release
      - PORT=8080
    ports:
      - "8080:8080"
    volumes:
      - ai_data:/data
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  ai_data:
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Service Won't Start

```bash
# Check logs
systemctl status ai-service
journalctl -u ai-service -f

# Common causes:
# - Missing API keys
# - Port already in use
# - Database permissions
```

#### 2. API Key Issues

```bash
# Test API keys manually
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models

# Check provider availability
curl http://localhost:8080/api/v1/providers
```

#### 3. Database Issues

```bash
# Reset database
rm ai_service.db

# Check permissions
ls -la ai_service.db
```

#### 4. Rate Limiting

```bash
# Check rate limit response
curl -v http://localhost:8080/api/v1/health

# HTTP 429 = Rate limited
# Adjust RATE_LIMIT_REQUESTS_PER_MINUTE
```

### Performance Tuning

#### 1. Concurrent Requests

- Service handles concurrent requests well
- PostgreSQL handles high concurrency well
- Consider PostgreSQL for high-traffic deployments

#### 2. Response Time Optimization

- Use faster models for time-sensitive requests
- Implement caching for common requests
- Monitor provider response times

#### 3. Memory Usage

- Service uses minimal memory (~50MB base)
- Memory usage scales with concurrent requests
- Monitor for memory leaks in long-running deployments

## 📝 API Integration Examples

### JavaScript/Node.js

```javascript
const response = await fetch('http://localhost:8080/api/v1/generate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    provider: 'openai',
    prompt: 'Generate a hello world function',
    max_tokens: 500
  })
});
const result = await response.json();
console.log(result.content);
```

### Python

```python
import requests

response = requests.post('http://localhost:8080/api/v1/generate', json={
    'provider': 'gemini',
    'prompt': 'Explain REST APIs',
    'max_tokens': 300
})
result = response.json()
print(result['content'])
```

### Go

```go
type GenerationRequest struct {
    Provider   string `json:"provider"`
    Prompt     string `json:"prompt"`
    MaxTokens  int    `json:"max_tokens"`
}

req := GenerationRequest{
    Provider:  "openai",
    Prompt:    "Write Go code",
    MaxTokens: 500,
}

// Make HTTP request...
```

## 🎯 Use Case Examples

### 1. Code Generation Service

- Generate boilerplate code
- Create documentation
- Debug and optimize code
- Generate test cases

### 2. Content Creation Platform

- Blog post generation
- Marketing copy
- Technical documentation
- Creative writing assistance

### 3. Data Analysis Tool

- Analyze datasets
- Generate insights
- Create reports
- Explain complex data

### 4. Educational Assistant

- Explain concepts
- Create learning materials
- Generate quizzes
- Provide tutoring

---

## 📞 Support

- **Documentation**: Check `/swagger/index.html` for API docs
- **Health Check**: Monitor `/api/v1/health` endpoint
- **Logs**: Check application logs for debugging
- **Performance**: Monitor `/api/v1/stats` for usage metrics

**Happy AI Generating! 🤖✨**
