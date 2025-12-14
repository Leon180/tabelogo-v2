package http

import "time"

// QuickSearchRequest represents the request for quick search
type QuickSearchRequest struct {
	PlaceID      string `json:"place_id" binding:"required"`
	LanguageCode string `json:"language_code" binding:"required,oneof=en ja zh-TW"`
	APIMask      string `json:"api_mask"`
}

// QuickSearchResponse represents the response for quick search
type QuickSearchResponse struct {
	Source   string      `json:"source"`
	CachedAt *time.Time  `json:"cached_at,omitempty"`
	Result   interface{} `json:"result"`
}

// AdvanceSearchRequest represents the request for advance search
type AdvanceSearchRequest struct {
	TextQuery      string       `json:"text_query" binding:"required"`
	LocationBias   LocationBias `json:"location_bias" binding:"required"`
	MaxResultCount int          `json:"max_result_count" binding:"required,min=1,max=20"`
	MinRating      float64      `json:"min_rating" binding:"omitempty,min=0,max=5"`
	OpenNow        bool         `json:"open_now"`
	RankPreference string       `json:"rank_preference" binding:"required,oneof=DISTANCE RELEVANCE"`
	LanguageCode   string       `json:"language_code" binding:"required,oneof=en ja zh-TW"`
	APIMask        string       `json:"api_mask"`
}

// LocationBias represents location bias for search
type LocationBias struct {
	Rectangle Rectangle `json:"rectangle" binding:"required"`
}

// Rectangle represents a rectangular area
type Rectangle struct {
	Low  Coordinates `json:"low" binding:"required"`
	High Coordinates `json:"high" binding:"required"`
}

// Coordinates represents geographic coordinates
type Coordinates struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// AdvanceSearchResponse represents the response for advance search
type AdvanceSearchResponse struct {
	Places         []interface{}  `json:"places"`
	TotalCount     int            `json:"total_count"`
	SearchMetadata SearchMetadata `json:"search_metadata"`
}

// SearchMetadata represents search metadata
type SearchMetadata struct {
	TextQuery    string `json:"text_query"`
	SearchTimeMs int64  `json:"search_time_ms"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status       string            `json:"status"`
	Timestamp    time.Time         `json:"timestamp"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}
