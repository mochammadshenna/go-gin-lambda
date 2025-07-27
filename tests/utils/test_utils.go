package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TestDB holds test database connection
type TestDB struct {
	*sql.DB
	DSN string
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	// Use SQLite for testing
	dsn := getTestDSN()

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Configure connection pool for tests
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return &TestDB{
		DB:  db,
		DSN: dsn,
	}
}

// getTestDSN returns test database connection string
func getTestDSN() string {
	dbname := getEnv("TEST_DB_NAME", ":memory:")
	return dbname
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// SetupTestDatabase sets up test database with migrations
func (tdb *TestDB) SetupTestDatabase(t *testing.T) {
	// SQLite-compatible schema
	schema := `
	CREATE TABLE IF NOT EXISTS generations (
		id TEXT PRIMARY KEY,
		provider TEXT NOT NULL,
		model TEXT NOT NULL,
		prompt TEXT NOT NULL,
		response TEXT NOT NULL,
		tokens_used INTEGER NOT NULL DEFAULT 0,
		duration_ms INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'success',
		error_message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS providers (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		api_key_hash TEXT,
		is_active BOOLEAN DEFAULT 1,
		config TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS stats (
		id TEXT PRIMARY KEY,
		provider TEXT NOT NULL,
		date DATE NOT NULL,
		total_generations INTEGER DEFAULT 0,
		total_tokens INTEGER DEFAULT 0,
		avg_duration_ms INTEGER DEFAULT 0,
		error_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(provider, date)
	);

	CREATE TABLE IF NOT EXISTS api_keys (
		id TEXT PRIMARY KEY,
		provider TEXT NOT NULL,
		key_hash TEXT NOT NULL,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_generations_provider ON generations(provider);
	CREATE INDEX IF NOT EXISTS idx_generations_created_at ON generations(created_at);
	CREATE INDEX IF NOT EXISTS idx_generations_status ON generations(status);
	CREATE INDEX IF NOT EXISTS idx_stats_provider_date ON stats(provider, date);
	CREATE INDEX IF NOT EXISTS idx_api_keys_provider ON api_keys(provider);
	`

	_, err := tdb.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to setup test database schema: %v", err)
	}

	log.Println("Test database setup completed")
}

// CleanupTestDatabase cleans up test data
func (tdb *TestDB) CleanupTestDatabase(t *testing.T) {
	tables := []string{"generations", "providers", "stats", "api_keys"}

	for _, table := range tables {
		_, err := tdb.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Logf("Warning: Failed to delete from table %s: %v", table, err)
		}
	}
}

// Close closes the test database connection
func (tdb *TestDB) Close() error {
	return tdb.DB.Close()
}

// MockAIProvider is a mock AI provider for testing
type MockAIProvider struct {
	ShouldFail bool
	Response   string
	Error      error
	Delay      time.Duration
}

// Generate simulates AI generation for testing
func (m *MockAIProvider) Generate(ctx context.Context, req interface{}) (interface{}, error) {
	if m.Delay > 0 {
		time.Sleep(m.Delay)
	}

	if m.ShouldFail {
		return nil, m.Error
	}

	return m.Response, nil
}

// IsAvailable returns availability status for testing
func (m *MockAIProvider) IsAvailable() bool {
	return !m.ShouldFail
}

// TestContext creates a test context with timeout
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// AssertNoError is a helper to assert no error occurred
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError is a helper to assert an error occurred
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got none", message)
	}
}

// AssertEqual is a helper to assert equality
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertNotNil is a helper to assert not nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Fatalf("%s: expected non-nil value", message)
	}
}

// CreateTestGeneration creates a test generation record
func CreateTestGeneration(t *testing.T, db *sql.DB, provider, model, prompt, response string) string {
	query := `
		INSERT INTO generations (
			id, provider, model, prompt, response, tokens_used, duration_ms, status
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	id := fmt.Sprintf("test-%d", time.Now().UnixNano())
	_, err := db.Exec(query, id, provider, model, prompt, response, 100, 1000, "success")
	AssertNoError(t, err, "Failed to create test generation")

	return id
}

// CreateTestProvider creates a test provider record
func CreateTestProvider(t *testing.T, db *sql.DB, name string, isActive bool) string {
	query := `
		INSERT INTO providers (
			id, name, is_active, config
		) VALUES (
			?, ?, ?, ?
		)
	`

	id := fmt.Sprintf("test-provider-%d", time.Now().UnixNano())
	config := fmt.Sprintf(`{"default_model": "test-model", "max_tokens": 1000}`)
	_, err := db.Exec(query, id, name, isActive, config)
	AssertNoError(t, err, "Failed to create test provider")

	return id
}
