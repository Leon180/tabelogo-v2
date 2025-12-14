package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCachedResult_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *CachedResult
		expected bool
	}{
		{
			name: "not expired - future expiry",
			setup: func() *CachedResult {
				return &CachedResult{
					PlaceID:   "test-place",
					Results:   []TabelogRestaurantDTO{},
					CachedAt:  time.Now(),
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}
			},
			expected: false,
		},
		{
			name: "expired - past expiry",
			setup: func() *CachedResult {
				return &CachedResult{
					PlaceID:   "test-place",
					Results:   []TabelogRestaurantDTO{},
					CachedAt:  time.Now().Add(-2 * time.Hour),
					ExpiresAt: time.Now().Add(-1 * time.Hour),
				}
			},
			expected: true,
		},
		{
			name: "just expired - exactly now",
			setup: func() *CachedResult {
				now := time.Now()
				return &CachedResult{
					PlaceID:   "test-place",
					Results:   []TabelogRestaurantDTO{},
					CachedAt:  now.Add(-1 * time.Hour),
					ExpiresAt: now.Add(-1 * time.Millisecond),
				}
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cached := tt.setup()

			// Act
			result := cached.IsExpired()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCachedResult_Structure(t *testing.T) {
	// Arrange
	placeID := "test-place-123"
	results := []TabelogRestaurantDTO{
		{
			Link:        "https://tabelog.com/1",
			Name:        "Restaurant 1",
			Rating:      3.5,
			RatingCount: 100,
			Bookmarks:   50,
			Phone:       "03-1234-5678",
			Types:       []string{"Japanese"},
			Photos:      []string{"photo1.jpg"},
		},
	}
	cachedAt := time.Now()
	expiresAt := cachedAt.Add(24 * time.Hour)

	// Act
	cached := &CachedResult{
		PlaceID:   placeID,
		Results:   results,
		CachedAt:  cachedAt,
		ExpiresAt: expiresAt,
	}

	// Assert
	assert.Equal(t, placeID, cached.PlaceID)
	assert.Len(t, cached.Results, 1)
	assert.Equal(t, "Restaurant 1", cached.Results[0].Name)
	assert.Equal(t, cachedAt, cached.CachedAt)
	assert.Equal(t, expiresAt, cached.ExpiresAt)
}

func TestCachedResult_EmptyResults(t *testing.T) {
	// Arrange
	cached := &CachedResult{
		PlaceID:   "test-place",
		Results:   []TabelogRestaurantDTO{},
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Assert
	assert.Empty(t, cached.Results)
	assert.False(t, cached.IsExpired())
}

func TestCachedResult_TTLCalculation(t *testing.T) {
	// Arrange
	cachedAt := time.Now()
	ttl := 2 * time.Hour
	expiresAt := cachedAt.Add(ttl)

	cached := &CachedResult{
		PlaceID:   "test-place",
		Results:   []TabelogRestaurantDTO{},
		CachedAt:  cachedAt,
		ExpiresAt: expiresAt,
	}

	// Act
	actualTTL := cached.ExpiresAt.Sub(cached.CachedAt)

	// Assert
	assert.Equal(t, ttl, actualTTL)
	assert.False(t, cached.IsExpired())
}
