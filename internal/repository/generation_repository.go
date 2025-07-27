package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ai-service/internal/model"
)

// GenerationRepository defines the interface for generation data access
type GenerationRepository interface {
	Create(ctx context.Context, generation *model.GenerationHistory) error
	GetByID(ctx context.Context, id string) (*model.GenerationHistory, error)
	GetByProvider(ctx context.Context, provider string, limit, offset int) ([]*model.GenerationHistory, error)
	GetRecent(ctx context.Context, limit, offset int) ([]*model.GenerationHistory, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*model.GenerationHistory, error)
	GetStats(ctx context.Context, startDate, endDate time.Time) ([]*model.ProviderStats, error)
	GetProviderStats(ctx context.Context, provider string, startDate, endDate time.Time) (*model.ProviderStats, error)
	UpdateStatus(ctx context.Context, id string, status string, errorMessage string) error
	Delete(ctx context.Context, id string) error
}

// generationRepository implements GenerationRepository
type generationRepository struct {
	db *sql.DB
	// TODO: Add SQLC generated querier when available
}

// NewGenerationRepository creates a new generation repository
func NewGenerationRepository(db *sql.DB) GenerationRepository {
	return &generationRepository{
		db: db,
	}
}

// Create saves a new generation record
func (r *generationRepository) Create(ctx context.Context, generation *model.GenerationHistory) error {
	query := `
		INSERT INTO generations (
			id, provider, model, prompt, response, tokens_used, duration_ms, status, error_message
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	// Generate UUID for SQLite
	if generation.ID == "" {
		generation.ID = fmt.Sprintf("gen-%d", time.Now().UnixNano())
	}

	_, err := r.db.ExecContext(ctx, query,
		generation.ID,
		generation.Provider,
		generation.Model,
		generation.Prompt,
		generation.Response,
		generation.TokensUsed,
		generation.Duration,
		generation.Status,
		generation.ErrorMessage,
	)

	if err != nil {
		return err
	}

	// Set timestamps
	generation.CreatedAt = time.Now()
	generation.UpdatedAt = time.Now()

	return nil
}

// GetByID retrieves a generation by ID
func (r *generationRepository) GetByID(ctx context.Context, id string) (*model.GenerationHistory, error) {
	query := `SELECT id, provider, model, prompt, response, tokens_used, duration_ms, status, error_message, created_at, updated_at FROM generations WHERE id = ?`

	var generation model.GenerationHistory
	var errorMessage sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&generation.ID,
		&generation.Provider,
		&generation.Model,
		&generation.Prompt,
		&generation.Response,
		&generation.TokensUsed,
		&generation.Duration,
		&generation.Status,
		&errorMessage,
		&generation.CreatedAt,
		&generation.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if errorMessage.Valid {
		generation.ErrorMessage = errorMessage.String
	}

	return &generation, nil
}

// GetByProvider retrieves generations by provider with pagination
func (r *generationRepository) GetByProvider(ctx context.Context, provider string, limit, offset int) ([]*model.GenerationHistory, error) {
	query := `
		SELECT id, provider, model, prompt, response, tokens_used, duration_ms, status, error_message, created_at, updated_at FROM generations 
		WHERE provider = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, provider, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var generations []*model.GenerationHistory
	for rows.Next() {
		var generation model.GenerationHistory
		var errorMessage sql.NullString
		err := rows.Scan(
			&generation.ID,
			&generation.Provider,
			&generation.Model,
			&generation.Prompt,
			&generation.Response,
			&generation.TokensUsed,
			&generation.Duration,
			&generation.Status,
			&errorMessage,
			&generation.CreatedAt,
			&generation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if errorMessage.Valid {
			generation.ErrorMessage = errorMessage.String
		}
		generations = append(generations, &generation)
	}

	return generations, nil
}

// GetRecent retrieves recent generations with pagination
func (r *generationRepository) GetRecent(ctx context.Context, limit, offset int) ([]*model.GenerationHistory, error) {
	query := `
		SELECT id, provider, model, prompt, response, tokens_used, duration_ms, status, error_message, created_at, updated_at FROM generations 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var generations []*model.GenerationHistory
	for rows.Next() {
		var generation model.GenerationHistory
		var errorMessage sql.NullString
		err := rows.Scan(
			&generation.ID,
			&generation.Provider,
			&generation.Model,
			&generation.Prompt,
			&generation.Response,
			&generation.TokensUsed,
			&generation.Duration,
			&generation.Status,
			&errorMessage,
			&generation.CreatedAt,
			&generation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if errorMessage.Valid {
			generation.ErrorMessage = errorMessage.String
		}
		generations = append(generations, &generation)
	}

	return generations, nil
}

// GetByStatus retrieves generations by status with pagination
func (r *generationRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*model.GenerationHistory, error) {
	query := `
		SELECT id, provider, model, prompt, response, tokens_used, duration_ms, status, error_message, created_at, updated_at FROM generations 
		WHERE status = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var generations []*model.GenerationHistory
	for rows.Next() {
		var generation model.GenerationHistory
		var errorMessage sql.NullString
		err := rows.Scan(
			&generation.ID,
			&generation.Provider,
			&generation.Model,
			&generation.Prompt,
			&generation.Response,
			&generation.TokensUsed,
			&generation.Duration,
			&generation.Status,
			&errorMessage,
			&generation.CreatedAt,
			&generation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if errorMessage.Valid {
			generation.ErrorMessage = errorMessage.String
		}
		generations = append(generations, &generation)
	}

	return generations, nil
}

// GetStats retrieves generation statistics for all providers
func (r *generationRepository) GetStats(ctx context.Context, startDate, endDate time.Time) ([]*model.ProviderStats, error) {
	query := `
		SELECT 
			provider,
			COUNT(*) as total_generations,
			SUM(tokens_used) as total_tokens,
			AVG(duration_ms) as avg_duration_ms,
			COUNT(CASE WHEN status = 'error' THEN 1 END) as error_count
		FROM generations 
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY provider
		ORDER BY total_generations DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*model.ProviderStats
	for rows.Next() {
		var stat model.ProviderStats
		err := rows.Scan(
			&stat.Provider,
			&stat.TotalGenerations,
			&stat.TotalTokens,
			&stat.AvgDuration,
			&stat.ErrorCount,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}

// GetProviderStats retrieves statistics for a specific provider
func (r *generationRepository) GetProviderStats(ctx context.Context, provider string, startDate, endDate time.Time) (*model.ProviderStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_generations,
			SUM(tokens_used) as total_tokens,
			AVG(duration_ms) as avg_duration_ms,
			COUNT(CASE WHEN status = 'error' THEN 1 END) as error_count
		FROM generations 
		WHERE provider = ? AND created_at >= ? AND created_at <= ?
	`

	var stat model.ProviderStats
	stat.Provider = provider

	err := r.db.QueryRowContext(ctx, query, provider, startDate, endDate).Scan(
		&stat.TotalGenerations,
		&stat.TotalTokens,
		&stat.AvgDuration,
		&stat.ErrorCount,
	)

	if err != nil {
		return nil, err
	}

	return &stat, nil
}

// UpdateStatus updates the status of a generation
func (r *generationRepository) UpdateStatus(ctx context.Context, id string, status string, errorMessage string) error {
	query := `
		UPDATE generations 
		SET status = ?, error_message = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, status, errorMessage, time.Now(), id)
	return err
}

// Delete removes a generation record
func (r *generationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM generations WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
