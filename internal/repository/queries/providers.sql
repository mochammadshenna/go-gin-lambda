-- name: CreateProvider :one
INSERT INTO providers (
    name, api_key_hash, is_active, config
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetProviderByID :one
SELECT * FROM providers WHERE id = $1;

-- name: GetProviderByName :one
SELECT * FROM providers WHERE name = $1;

-- name: GetActiveProviders :many
SELECT * FROM providers WHERE is_active = true ORDER BY name;

-- name: GetAllProviders :many
SELECT * FROM providers ORDER BY name;

-- name: UpdateProvider :exec
UPDATE providers 
SET name = $2, api_key_hash = $3, is_active = $4, config = $5, updated_at = NOW()
WHERE id = $1;

-- name: UpdateProviderConfig :exec
UPDATE providers 
SET config = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateProviderStatus :exec
UPDATE providers 
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteProvider :exec
DELETE FROM providers WHERE id = $1;

-- name: GetProviderStats :many
SELECT 
    p.name,
    p.is_active,
    COUNT(g.id) as total_generations,
    COALESCE(SUM(g.tokens_used), 0) as total_tokens,
    COALESCE(AVG(g.duration_ms), 0) as avg_duration_ms,
    COUNT(CASE WHEN g.status = 'error' THEN 1 END) as error_count
FROM providers p
LEFT JOIN generations g ON p.name = g.provider
WHERE g.created_at >= $1 AND g.created_at <= $2 OR g.created_at IS NULL
GROUP BY p.id, p.name, p.is_active
ORDER BY total_generations DESC; 