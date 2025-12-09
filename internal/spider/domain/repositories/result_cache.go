package repositories

import (
	"context"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/spider/domain/models"
)

// ResultCacheRepository defines the interface for caching scrape results
type ResultCacheRepository interface {
	// Get retrieves cached results for a place
	Get(ctx context.Context, placeID string) (*models.CachedResult, error)

	// Set stores results in cache with TTL
	Set(ctx context.Context, placeID string, results []models.TabelogRestaurant, ttl time.Duration) error

	// Delete removes cached results
	Delete(ctx context.Context, placeID string) error
}
