package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `json:"server"`

	// Database configuration
	Database DatabaseConfig `json:"database"`

	// AI Providers configuration
	AIProviders AIProvidersConfig `json:"ai_providers"`

	// Security configuration
	Security SecurityConfig `json:"security"`

	// Observability configuration
	Observability ObservabilityConfig `json:"observability"`

	// Rate limiting configuration
	RateLimit RateLimitConfig `json:"rate_limit"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string        `json:"port"`
	Host         string        `json:"host"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	Environment  string        `json:"environment"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

// AIProvidersConfig represents AI providers configuration
type AIProvidersConfig struct {
	OpenAI    OpenAIConfig    `json:"openai"`
	Gemini    GeminiConfig    `json:"gemini"`
	Anthropic AnthropicConfig `json:"anthropic"`
}

// OpenAIConfig represents OpenAI configuration
type OpenAIConfig struct {
	APIKey       string `json:"api_key"`
	BaseURL      string `json:"base_url"`
	DefaultModel string `json:"default_model"`
	MaxTokens    int    `json:"max_tokens"`
}

// GeminiConfig represents Google Gemini configuration
type GeminiConfig struct {
	APIKey       string `json:"api_key"`
	DefaultModel string `json:"default_model"`
	MaxTokens    int    `json:"max_tokens"`
}

// AnthropicConfig represents Anthropic configuration
type AnthropicConfig struct {
	APIKey       string `json:"api_key"`
	DefaultModel string `json:"default_model"`
	MaxTokens    int    `json:"max_tokens"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	JWTSecret     string        `json:"jwt_secret"`
	JWTExpiration time.Duration `json:"jwt_expiration"`
	CORSOrigins   []string      `json:"cors_origins"`
}

// ObservabilityConfig represents observability configuration
type ObservabilityConfig struct {
	LogLevel    string `json:"log_level"`
	SentryDSN   string `json:"sentry_dsn"`
	Environment string `json:"environment"`
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int           `json:"requests_per_minute"`
	BurstSize         int           `json:"burst_size"`
	WindowSize        time.Duration `json:"window_size"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist
	}

	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "./ai_service.db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		AIProviders: AIProvidersConfig{
			OpenAI: OpenAIConfig{
				APIKey:       getEnv("OPENAI_API_KEY", ""),
				BaseURL:      getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
				DefaultModel: getEnv("OPENAI_DEFAULT_MODEL", "gpt-3.5-turbo"),
				MaxTokens:    getIntEnv("OPENAI_MAX_TOKENS", 4096),
			},
			Gemini: GeminiConfig{
				APIKey:       getEnv("GEMINI_API_KEY", ""),
				DefaultModel: getEnv("GEMINI_DEFAULT_MODEL", "gemini-1.5-flash"),
				MaxTokens:    getIntEnv("GEMINI_MAX_TOKENS", 8192),
			},
			Anthropic: AnthropicConfig{
				APIKey:       getEnv("ANTHROPIC_API_KEY", ""),
				DefaultModel: getEnv("ANTHROPIC_DEFAULT_MODEL", "claude-3-sonnet-20240229"),
				MaxTokens:    getIntEnv("ANTHROPIC_MAX_TOKENS", 8192),
			},
		},
		Security: SecurityConfig{
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
			JWTExpiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
			CORSOrigins:   getStringSliceEnv("CORS_ORIGINS", []string{"*"}),
		},
		Observability: ObservabilityConfig{
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			SentryDSN:   getEnv("SENTRY_DSN", ""),
			Environment: getEnv("ENVIRONMENT", "development"),
			ServiceName: getEnv("SERVICE_NAME", "ai-service"),
			Version:     getEnv("VERSION", "1.0.0"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			BurstSize:         getIntEnv("RATE_LIMIT_BURST_SIZE", 10),
			WindowSize:        getDurationEnv("RATE_LIMIT_WINDOW_SIZE", time.Minute),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}

	// Validate at least one AI provider is configured
	if c.AIProviders.OpenAI.APIKey == "" &&
		c.AIProviders.Gemini.APIKey == "" &&
		c.AIProviders.Anthropic.APIKey == "" {
		return fmt.Errorf("at least one AI provider API key is required")
	}

	return nil
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated values
		// In production, you might want more sophisticated parsing
		return []string{value}
	}
	return defaultValue
}
