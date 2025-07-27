-- name: CreateStats :one
INSERT INTO stats (
    provider, date, total_generations, total_tokens, avg_duration_ms, error_count
) VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT (provider, date) 
DO UPDATE SET 
    total_generations = stats.total_generations + EXCLUDED.total_generations,
    total_tokens = stats.total_tokens + EXCLUDED.total_tokens,
    avg_duration_ms = (stats.avg_duration_ms + EXCLUDED.avg_duration_ms) / 2,
    error_count = stats.error_count + EXCLUDED.error_count,
    updated_at = NOW()
RETURNING *;

-- name: GetStatsByProvider :many
SELECT * FROM stats 
WHERE provider = $1 
ORDER BY date DESC 
LIMIT $2 OFFSET $3;

-- name: GetStatsByDateRange :many
SELECT * FROM stats 
WHERE date >= $1 AND date <= $2 
ORDER BY date DESC, provider;

-- name: GetDailyStats :many
SELECT 
    date,
    SUM(total_generations) as total_generations,
    SUM(total_tokens) as total_tokens,
    AVG(avg_duration_ms) as avg_duration_ms,
    SUM(error_count) as error_count
FROM stats 
WHERE date >= $1 AND date <= $2
GROUP BY date
ORDER BY date DESC;

-- name: GetProviderDailyStats :many
SELECT 
    provider,
    date,
    total_generations,
    total_tokens,
    avg_duration_ms,
    error_count
FROM stats 
WHERE provider = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: GetTopProviders :many
SELECT 
    provider,
    SUM(total_generations) as total_generations,
    SUM(total_tokens) as total_tokens,
    AVG(avg_duration_ms) as avg_duration_ms,
    SUM(error_count) as error_count
FROM stats 
WHERE date >= $1 AND date <= $2
GROUP BY provider
ORDER BY total_generations DESC
LIMIT $3;

-- name: DeleteStatsByDate :exec
DELETE FROM stats WHERE date = $1;

-- name: DeleteStatsByProvider :exec
DELETE FROM stats WHERE provider = $1; 