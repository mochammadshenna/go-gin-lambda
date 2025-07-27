package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DBType      string
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	DBPath      string
	SSLMode     string
	MaxOpen     int
	MaxIdle     int
	MaxLifetime time.Duration
}

// NewDB creates a new database connection (SQLite or PostgreSQL)
func NewDB() *DB {
	config := getDatabaseConfig()

	var db *sql.DB
	var err error

	switch config.DBType {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
		db, err = sql.Open("postgres", dsn)
		log.Println("Connecting to PostgreSQL database...")
	case "sqlite", "":
		db, err = sql.Open("sqlite3", config.DBPath)
		log.Println("Connecting to SQLite database...")
	default:
		log.Fatalf("Unsupported database type: %s", config.DBType)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpen)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to %s database", config.DBType)
	return &DB{db}
}

// getDatabaseConfig returns database configuration from environment variables
func getDatabaseConfig() DatabaseConfig {
	dbType := getEnv("DB_TYPE", "sqlite")

	config := DatabaseConfig{
		DBType:      dbType,
		MaxOpen:     getEnvAsInt("DB_MAX_OPEN", 25),
		MaxIdle:     getEnvAsInt("DB_MAX_IDLE", 5),
		MaxLifetime: getEnvAsDuration("DB_MAX_LIFETIME", 5*time.Minute),
	}

	if dbType == "postgres" {
		config.Host = getEnv("DB_HOST", "localhost")
		config.Port = getEnv("DB_PORT", "5432")
		config.User = getEnv("DB_USER", "postgres")
		config.Password = getEnv("DB_PASSWORD", "password")
		config.DBName = getEnv("DB_NAME", "ai_service")
		config.SSLMode = getEnv("DB_SSLMODE", "disable")
	} else {
		config.DBPath = getEnv("DB_PATH", "ai_service.db")
	}

	return config
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets environment variable as int with fallback
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d", &fallback); err == nil && intValue == 1 {
			return fallback
		}
	}
	return fallback
}

// getEnvAsDuration gets environment variable as duration with fallback
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Ping checks if the database is accessible
func (db *DB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

// NewRedisDB creates a new Redis connection (placeholder for future implementation)
func NewRedisDB() interface{} {
	// Placeholder for Redis connection
	// TODO: Implement Redis connection when needed
	return nil
}
