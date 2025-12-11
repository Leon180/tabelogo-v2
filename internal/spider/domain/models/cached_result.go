package models

import "time"

// CachedResult represents cached scrape results
// Uses DTO for JSON serialization to Redis
type CachedResult struct {
	PlaceID   string                 `json:"place_id"`
	Results   []TabelogRestaurantDTO `json:"results"`
	CachedAt  time.Time              `json:"cached_at"`
	ExpiresAt time.Time              `json:"expires_at"`
}

// IsExpired checks if the cached result has expired
func (c *CachedResult) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}
