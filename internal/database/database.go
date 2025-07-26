package database

import (
	"ai-service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func New(dbPath string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.GenerationHistory{})
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (db *DB) SaveGeneration(history *models.GenerationHistory) error {
	return db.Create(history).Error
}

func (db *DB) GetGenerationHistory(limit int) ([]models.GenerationHistory, error) {
	var histories []models.GenerationHistory
	err := db.Order("created_at desc").Limit(limit).Find(&histories).Error
	return histories, err
}

func (db *DB) GetGenerationsByProvider(provider string, limit int) ([]models.GenerationHistory, error) {
	var histories []models.GenerationHistory
	err := db.Where("provider = ?", provider).Order("created_at desc").Limit(limit).Find(&histories).Error
	return histories, err
}

func (db *DB) GetGenerationStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total generations
	var totalCount int64
	if err := db.Model(&models.GenerationHistory{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_generations"] = totalCount
	
	// Generations by provider
	var providerStats []struct {
		Provider string
		Count    int64
	}
	if err := db.Model(&models.GenerationHistory{}).
		Select("provider, count(*) as count").
		Group("provider").
		Scan(&providerStats).Error; err != nil {
		return nil, err
	}
	stats["by_provider"] = providerStats
	
	// Average tokens used
	var avgTokens float64
	if err := db.Model(&models.GenerationHistory{}).
		Select("avg(tokens_used) as avg_tokens").
		Scan(&avgTokens).Error; err != nil {
		return nil, err
	}
	stats["avg_tokens_used"] = avgTokens
	
	// Average duration
	var avgDuration float64
	if err := db.Model(&models.GenerationHistory{}).
		Select("avg(duration) as avg_duration").
		Scan(&avgDuration).Error; err != nil {
		return nil, err
	}
	stats["avg_duration_ms"] = avgDuration
	
	return stats, nil
}