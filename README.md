# AI Service - Multi-Provider AI Content Generation

A comprehensive AI service built with Go and Gin that supports multiple AI providers (OpenAI, Google Gemini, Anthropic) for content generation. Features a modern web interface, REST API, and comprehensive statistics tracking.

## 🚀 Features

### Core Functionality
- **Multi-Provider Support**: OpenAI GPT, Google Gemini, Anthropic Claude
- **Content Generation**: Text generation with customizable parameters
- **Provider Comparison**: Compare responses from multiple AI providers
- **Model Selection**: Dynamic model loading based on provider
- **Parameter Control**: Temperature, max tokens, system messages

### Web Interface
- **Modern UI**: Bootstrap 5 with Font Awesome icons
- **Real-time Generation**: AJAX-based content generation
- **Provider Information**: Detailed provider capabilities and pricing
- **Quick Examples**: Pre-filled examples for different use cases
- **Generation History**: View and filter past generations
- **Statistics Dashboard**: Usage analytics and performance insights

### API Features
- **RESTful API**: Complete REST API with Swagger documentation
- **Rate Limiting**: Configurable request rate limiting
- **Error Handling**: Comprehensive error handling and logging
- **Health Monitoring**: Service health checks and provider status
- **Database Storage**: SQLite database for generation history

## 📋 Prerequisites

- Go 1.21 or higher
- SQLite (included)
- API keys for desired providers

## 🛠️ Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd go-gin-lambda
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Set Up Environment Variables
Create a `.env` file in the root directory:

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

### 4. Get API Keys

#### Google Gemini
1. Go to [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with your Google account
3. Click "Create API Key"
4. Copy the generated API key
5. Add it to your `.env` file as `GEMINI_API_KEY=your_key_here`

#### OpenAI
1. Go to [OpenAI Platform](https://platform.openai.com/api-keys)
2. Sign in to your account
3. Click "Create new secret key"
4. Copy the generated API key
5. Add it to your `.env` file as `OPENAI_API_KEY=your_key_here`

#### Anthropic (Optional)
1. Go to [Anthropic Console](https://console.anthropic.com/)
2. Sign in to your account
3. Navigate to API Keys
4. Create a new API key
5. Add it to your `.env` file as `ANTHROPIC_API_KEY=your_key_here`

## 🏃‍♂️ Running the Service

### Development Mode
```bash
# Build and run
make run

# Or run directly
go run main.go
```

### Production Mode
```bash
# Build the binary
go build -o ai-service main.go

# Run with production settings
GIN_MODE=release ./ai-service
```

## 🌐 Accessing the Service

Once the service is running, you can access:

- **Web Interface**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/api/v1/health

## 📚 API Documentation

### Core Endpoints

#### Generate Content
```bash
POST /api/v1/generate
Content-Type: application/json

{
  "provider": "gemini",
  "model": "gemini-1.5-flash",
  "prompt": "Hello, how are you?",
  "maxTokens": 100,
  "temperature": 0.7,
  "systemMessage": "You are a helpful assistant."
}
```

#### Compare Providers
```bash
POST /api/v1/compare
Content-Type: application/json

{
  "prompt": "Explain quantum computing",
  "providers": ["openai", "gemini"],
  "maxTokens": 500,
  "temperature": 0.7
}
```

#### Get Providers
```bash
GET /api/v1/providers
```

#### Get History
```bash
GET /api/v1/history?limit=50&provider=gemini
```

#### Get Statistics
```bash
GET /api/v1/stats
```

### Example Usage

#### Using curl
```bash
# Generate content with Gemini
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "gemini",
    "model": "gemini-1.5-flash",
    "prompt": "Write a function in Go that validates an email address",
    "maxTokens": 200
  }'

# Get service health
curl http://localhost:8080/api/v1/health
```

## 🏗️ Project Structure

```
go-gin-lambda/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── .env.example           # Environment variables example
├── SETUP.md               # Detailed setup instructions
├── README.md              # This file
├── templates/             # HTML templates
│   ├── home_standalone.html
│   ├── history_standalone.html
│   ├── stats_standalone.html
│   ├── error.html
│   └── test.html
└── internal/              # Internal application code
    ├── ai/                # AI provider implementations
    │   ├── interface.go   # Provider interface
    │   ├── manager.go     # AI manager
    │   ├── openai.go      # OpenAI provider
    │   └── gemini.go      # Gemini provider
    ├── config/            # Configuration management
    │   └── config.go
    ├── database/          # Database operations
    │   └── database.go
    ├── handlers/          # HTTP handlers
    │   ├── handlers.go    # API handlers
    │   └── web.go         # Web interface handlers
    ├── middleware/        # HTTP middleware
    │   └── middleware.go
    └── models/            # Data models
        └── models.go
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin framework mode | `debug` |
| `LOG_LEVEL` | Logging level | `info` |
| `DB_PATH` | SQLite database path | `./ai_service.db` |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Rate limiting | `60` |
| `DEFAULT_AI_PROVIDER` | Default provider | `openai` |

### Supported AI Providers

#### OpenAI
- **Models**: GPT-3.5-turbo, GPT-4, GPT-4-turbo
- **Max Tokens**: 4096
- **Pricing**: Pay per token (~$0.002/1K tokens for GPT-3.5)
- **Best For**: General text generation, code completion, creative writing

#### Google Gemini
- **Models**: gemini-1.5-flash, gemini-1.5-pro, gemini-1.0-pro
- **Max Tokens**: 8192
- **Pricing**: Free tier available, pay per token for pro usage
- **Best For**: Multimodal tasks, long document analysis, cost-effective solutions

#### Anthropic Claude (Coming Soon)
- **Models**: Claude-3-Sonnet, Claude-3-Haiku, Claude-3-Opus
- **Max Tokens**: 8192
- **Pricing**: Pay per token (premium pricing)
- **Best For**: Safety-critical applications, research and analysis

## 🐛 Troubleshooting

### Common Issues

#### "Provider not found" Error
- Check that you have at least one API key configured
- Verify the provider name in your request
- Ensure the provider is available in the service

#### "Invalid API key" Error
- Regenerate your API key from the provider's dashboard
- Make sure you're using the correct key format
- Check that the key has the necessary permissions

#### "Rate limit exceeded" Error
- The service has built-in rate limiting (60 requests per minute by default)
- Wait a minute before making more requests
- Adjust the rate limit in your `.env` file if needed

#### Template Errors
- Ensure all template files are present in the `templates/` directory
- Check that the service has read permissions for template files
- Restart the service after making template changes

### Debug Mode

Enable debug logging by setting:
```env
LOG_LEVEL=debug
GIN_MODE=debug
```

### Health Check

Monitor service health:
```bash
curl http://localhost:8080/api/v1/health
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Bootstrap](https://getbootstrap.com/)
- [Font Awesome](https://fontawesome.com/)
- [OpenAI API](https://openai.com/api/)
- [Google Gemini API](https://ai.google.dev/)
- [Anthropic Claude API](https://www.anthropic.com/)

## 📞 Support

For support and questions:
- Create an issue in the repository
- Check the [SETUP.md](SETUP.md) file for detailed setup instructions
- Review the API documentation at http://localhost:8080/swagger/index.html

---

**Note**: This service requires valid API keys from the respective AI providers to function. Make sure to follow the provider's terms of service and usage guidelines.
