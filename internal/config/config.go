package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	GinMode                 string
	OpenAIAPIKey            string
	GeminiAPIKey            string
	AnthropicAPIKey         string
	DBPath                  string
	ServiceName             string
	ServiceVersion          string
	LogLevel                string
	RateLimitRequestsPerMin int
	DefaultAIProvider       string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	rateLimitStr := getEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", "60")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		rateLimit = 60
	}

	return &Config{
		Port:                    getEnv("PORT", "8080"),
		GinMode:                 getEnv("GIN_MODE", "debug"),
		OpenAIAPIKey:            getEnv("OPENAI_API_KEY", ""),
		GeminiAPIKey:            getEnv("GEMINI_API_KEY", ""),
		AnthropicAPIKey:         getEnv("ANTHROPIC_API_KEY", ""),
		DBPath:                  getEnv("DB_PATH", "./ai_service.db"),
		ServiceName:             getEnv("SERVICE_NAME", "ai-service"),
		ServiceVersion:          getEnv("SERVICE_VERSION", "1.0.0"),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		RateLimitRequestsPerMin: rateLimit,
		DefaultAIProvider:       getEnv("DEFAULT_AI_PROVIDER", "openai"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}