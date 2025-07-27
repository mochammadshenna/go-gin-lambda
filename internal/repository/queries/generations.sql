-- name: CreateGeneration :one
INSERT INTO generations (
    provider, model, prompt, response, tokens_used, duration_ms, status, error_message
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetGenerationByID :one
SELECT * FROM generations WHERE id = $1;

-- name: GetGenerationsByProvider :many
SELECT * FROM generations 
WHERE provider = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetRecentGenerations :many
SELECT * FROM generations 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2;

-- name: GetGenerationsByStatus :many
SELECT * FROM generations 
WHERE status = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetGenerationStats :many
SELECT 
    provider,
    COUNT(*) as total_generations,
    SUM(tokens_used) as total_tokens,
    AVG(duration_ms) as avg_duration_ms,
    COUNT(CASE WHEN status = 'error' THEN 1 END) as error_count
FROM generations 
WHERE created_at >= $1 AND created_at <= $2
GROUP BY provider
ORDER BY total_generations DESC;

-- name: GetProviderStats :one
SELECT 
    COUNT(*) as total_generations,
    SUM(tokens_used) as total_tokens,
    AVG(duration_ms) as avg_duration_ms,
    COUNT(CASE WHEN status = 'error' THEN 1 END) as error_count
FROM generations 
WHERE provider = $1 AND created_at >= $2 AND created_at <= $3;

-- name: DeleteGeneration :exec
DELETE FROM generations WHERE id = $1;

-- name: UpdateGenerationStatus :exec
UPDATE generations 
SET status = $2, error_message = $3, updated_at = NOW()
WHERE id = $1; 