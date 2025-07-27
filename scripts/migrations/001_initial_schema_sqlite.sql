-- SQLite-compatible initial schema
-- Note: SQLite doesn't support UUID extension, so we use TEXT PRIMARY KEY

-- Create generations table
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

-- Create providers table
CREATE TABLE IF NOT EXISTS providers (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    api_key_hash TEXT,
    is_active BOOLEAN DEFAULT 1,
    config TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create stats table for aggregated metrics
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

-- Create api_keys table for secure key management
CREATE TABLE IF NOT EXISTS api_keys (
    id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    key_hash TEXT NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_generations_provider ON generations(provider);
CREATE INDEX IF NOT EXISTS idx_generations_created_at ON generations(created_at);
CREATE INDEX IF NOT EXISTS idx_generations_status ON generations(status);
CREATE INDEX IF NOT EXISTS idx_stats_provider_date ON stats(provider, date);
CREATE INDEX IF NOT EXISTS idx_api_keys_provider ON api_keys(provider);

-- Insert default providers
INSERT OR IGNORE INTO providers (id, name, is_active, config) VALUES
    ('prov-openai', 'openai', 1, '{"default_model": "gpt-3.5-turbo", "max_tokens": 4096}'),
    ('prov-gemini', 'gemini', 1, '{"default_model": "gemini-1.5-flash", "max_tokens": 8192}'),
    ('prov-anthropic', 'anthropic', 0, '{"default_model": "claude-3-sonnet", "max_tokens": 8192}'); 