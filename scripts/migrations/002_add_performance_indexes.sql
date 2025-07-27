-- Add composite indexes for common query patterns
CREATE INDEX idx_generations_provider_created_at ON generations(provider, created_at DESC);
CREATE INDEX idx_generations_status_created_at ON generations(status, created_at DESC);

-- Add partial indexes for active records
CREATE INDEX idx_providers_active ON providers(name) WHERE is_active = true;
CREATE INDEX idx_api_keys_active ON api_keys(provider) WHERE is_active = true;

-- Add index for stats aggregation queries
CREATE INDEX idx_stats_date_range ON stats(date DESC, provider);

-- Add index for text search on prompts (if needed for future features)
CREATE INDEX idx_generations_prompt_gin ON generations USING gin(to_tsvector('english', prompt));

-- Add index for JSONB queries on provider config
CREATE INDEX idx_providers_config_gin ON providers USING gin(config); 