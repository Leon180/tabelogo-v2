package testutil

import (
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
	"github.com/google/uuid"
)

// CreateTestJob creates a test scraping job with default values
func CreateTestJob() *models.ScrapingJob {
	return models.NewScrapingJob("test-google-id", "Tokyo", "Test Restaurant")
}

// CreateTestJobWithID creates a test job with a specific ID
func CreateTestJobWithID(id uuid.UUID) *models.ScrapingJob {
	job := CreateTestJob()
	// Note: We can't directly set the ID, so we create a new job
	// This is a limitation of the current design
	return job
}

// CreateTestRestaurant creates a test Tabelog restaurant
func CreateTestRestaurant() *models.TabelogRestaurant {
	return models.NewTabelogRestaurant(
		"https://tabelog.com/test",
		"Test Restaurant",
		3.5,
		100,
		50,
		"03-1234-5678",
		[]string{"Japanese", "Sushi"},
		[]string{"https://example.com/photo1.jpg"},
	)
}

// CreateTestRestaurants creates multiple test restaurants
func CreateTestRestaurants(count int) []*models.TabelogRestaurant {
	restaurants := make([]*models.TabelogRestaurant, count)
	for i := 0; i < count; i++ {
		restaurants[i] = CreateTestRestaurant()
	}
	return restaurants
}

// CreateTestCachedResult creates a test cached result
func CreateTestCachedResult(placeID string, results []*models.TabelogRestaurant) *models.CachedResult {
	dtos := make([]models.TabelogRestaurantDTO, len(results))
	for i, r := range results {
		dtos[i] = r.ToDTO()
	}

	return &models.CachedResult{
		PlaceID:   placeID,
		Results:   dtos,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

// CreateExpiredCachedResult creates an expired cached result for testing
func CreateExpiredCachedResult(placeID string) *models.CachedResult {
	return &models.CachedResult{
		PlaceID:   placeID,
		Results:   []models.TabelogRestaurantDTO{},
		CachedAt:  time.Now().Add(-48 * time.Hour),
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}
}
