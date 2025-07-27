package unit

import (
	"testing"
	"time"

	"ai-service/internal/model"
	"ai-service/internal/repository"
	"ai-service/tests/utils"
)

func TestGenerationRepository_Create(t *testing.T) {
	// Setup
	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	testDB.SetupTestDatabase(t)
	defer testDB.CleanupTestDatabase(t)

	repo := repository.NewGenerationRepository(testDB.DB)
	ctx := utils.TestContext(t)

	// Test data
	generation := &model.GenerationHistory{
		Provider:   "openai",
		Model:      "gpt-3.5-turbo",
		Prompt:     "Test prompt",
		Response:   "Test response",
		TokensUsed: 100,
		Duration:   1000,
		Status:     "success",
	}

	// Execute
	err := repo.Create(ctx, generation)

	// Assert
	utils.AssertNoError(t, err, "Failed to create generation")
	utils.AssertNotNil(t, generation.ID, "Generation ID should not be nil")
	utils.AssertNotNil(t, generation.CreatedAt, "CreatedAt should not be nil")
	utils.AssertNotNil(t, generation.UpdatedAt, "UpdatedAt should not be nil")
}

func TestGenerationRepository_GetByID(t *testing.T) {
	// Setup
	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	testDB.SetupTestDatabase(t)
	defer testDB.CleanupTestDatabase(t)

	repo := repository.NewGenerationRepository(testDB.DB)
	ctx := utils.TestContext(t)

	// Create test data
	generationID := utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt", "Test response")

	// Execute
	generation, err := repo.GetByID(ctx, generationID)

	// Assert
	utils.AssertNoError(t, err, "Failed to get generation by ID")
	utils.AssertNotNil(t, generation, "Generation should not be nil")
	utils.AssertEqual(t, generationID, generation.ID, "Generation ID should match")
	utils.AssertEqual(t, "openai", generation.Provider, "Provider should match")
	utils.AssertEqual(t, "gpt-3.5-turbo", generation.Model, "Model should match")
}

func TestGenerationRepository_GetByProvider(t *testing.T) {
	// Setup
	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	testDB.SetupTestDatabase(t)
	defer testDB.CleanupTestDatabase(t)

	repo := repository.NewGenerationRepository(testDB.DB)
	ctx := utils.TestContext(t)

	// Create test data
	utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt 1", "Test response 1")
	utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt 2", "Test response 2")
	utils.CreateTestGeneration(t, testDB.DB, "gemini", "gemini-1.5-flash", "Test prompt 3", "Test response 3")

	// Execute
	generations, err := repo.GetByProvider(ctx, "openai", 10, 0)

	// Assert
	utils.AssertNoError(t, err, "Failed to get generations by provider")
	utils.AssertEqual(t, 2, len(generations), "Should return 2 generations for openai")

	for _, gen := range generations {
		utils.AssertEqual(t, "openai", gen.Provider, "All generations should be from openai provider")
	}
}

func TestGenerationRepository_GetRecent(t *testing.T) {
	// Setup
	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	testDB.SetupTestDatabase(t)
	defer testDB.CleanupTestDatabase(t)

	repo := repository.NewGenerationRepository(testDB.DB)
	ctx := utils.TestContext(t)

	// Create test data
	utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt 1", "Test response 1")
	utils.CreateTestGeneration(t, testDB.DB, "gemini", "gemini-1.5-flash", "Test prompt 2", "Test response 2")

	// Execute
	generations, err := repo.GetRecent(ctx, 10, 0)

	// Assert
	utils.AssertNoError(t, err, "Failed to get recent generations")
	utils.AssertEqual(t, 2, len(generations), "Should return 2 generations")

	// Check that they are ordered by created_at DESC (most recent first)
	if len(generations) >= 2 {
		utils.AssertEqual(t, true, generations[0].CreatedAt.After(generations[1].CreatedAt) ||
			generations[0].CreatedAt.Equal(generations[1].CreatedAt),
			"Generations should be ordered by created_at DESC")
	}
}

func TestGenerationRepository_GetStats(t *testing.T) {
	// Setup
	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	testDB.SetupTestDatabase(t)
	defer testDB.CleanupTestDatabase(t)

	repo := repository.NewGenerationRepository(testDB.DB)
	ctx := utils.TestContext(t)

	// Create test data
	utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt 1", "Test response 1")
	utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt 2", "Test response 2")
	utils.CreateTestGeneration(t, testDB.DB, "gemini", "gemini-1.5-flash", "Test prompt 3", "Test response 3")

	// Execute
	startDate := time.Now().AddDate(0, 0, -1) // Yesterday
	endDate := time.Now().AddDate(0, 0, 1)    // Tomorrow
	stats, err := repo.GetStats(ctx, startDate, endDate)

	// Assert
	utils.AssertNoError(t, err, "Failed to get generation stats")
	utils.AssertEqual(t, 2, len(stats), "Should return stats for 2 providers")

	// Find openai stats
	var openaiStats *model.ProviderStats
	for _, stat := range stats {
		if stat.Provider == "openai" {
			openaiStats = stat
			break
		}
	}

	utils.AssertNotNil(t, openaiStats, "Should have stats for openai")
	utils.AssertEqual(t, 2, openaiStats.TotalGenerations, "OpenAI should have 2 generations")
}

func TestGenerationRepository_UpdateStatus(t *testing.T) {
	// Setup
	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	testDB.SetupTestDatabase(t)
	defer testDB.CleanupTestDatabase(t)

	repo := repository.NewGenerationRepository(testDB.DB)
	ctx := utils.TestContext(t)

	// Create test data
	generationID := utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt", "Test response")

	// Execute
	err := repo.UpdateStatus(ctx, generationID, "error", "Test error message")

	// Assert
	utils.AssertNoError(t, err, "Failed to update generation status")

	// Verify the update
	generation, err := repo.GetByID(ctx, generationID)
	utils.AssertNoError(t, err, "Failed to get updated generation")
	utils.AssertEqual(t, "error", generation.Status, "Status should be updated to error")
	utils.AssertEqual(t, "Test error message", generation.ErrorMessage, "Error message should be updated")
}

func TestGenerationRepository_Delete(t *testing.T) {
	// Setup
	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	testDB.SetupTestDatabase(t)
	defer testDB.CleanupTestDatabase(t)

	repo := repository.NewGenerationRepository(testDB.DB)
	ctx := utils.TestContext(t)

	// Create test data
	generationID := utils.CreateTestGeneration(t, testDB.DB, "openai", "gpt-3.5-turbo", "Test prompt", "Test response")

	// Execute
	err := repo.Delete(ctx, generationID)

	// Assert
	utils.AssertNoError(t, err, "Failed to delete generation")

	// Verify deletion
	_, err = repo.GetByID(ctx, generationID)
	utils.AssertError(t, err, "Generation should not exist after deletion")
}
