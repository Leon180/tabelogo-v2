package models

import "time"

// CachedResult represents cached scrape results
type CachedResult struct {
	PlaceID   string
	Results   []TabelogRestaurant
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired checks if the cached result has expired
func (c *CachedResult) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}
