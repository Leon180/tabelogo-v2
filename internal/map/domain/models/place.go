package models

import "time"

// QuickSearchRequest - 快速搜索請求
type QuickSearchRequest struct {
	PlaceID      string `json:"place_id" binding:"required"`
	LanguageCode string `json:"language_code" binding:"required,oneof=en ja zh-TW"`
	APIMask      string `json:"api_mask"`
}

// QuickSearchResponse - 快速搜索響應
type QuickSearchResponse struct {
	Source   string      `json:"source"`
	CachedAt *time.Time  `json:"cached_at,omitempty"`
	Result   interface{} `json:"result"`
}

// AdvanceSearchRequest - 高級搜索請求
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

// LocationBias - 位置偏好
type LocationBias struct {
	Rectangle Rectangle `json:"rectangle" binding:"required"`
}

// Rectangle - 矩形範圍
type Rectangle struct {
	Low  Coordinates `json:"low" binding:"required"`
	High Coordinates `json:"high" binding:"required"`
}

// Coordinates - 座標
type Coordinates struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// AdvanceSearchResponse - 高級搜索響應
type AdvanceSearchResponse struct {
	Places         []interface{}  `json:"places"`
	TotalCount     int            `json:"total_count"`
	SearchMetadata SearchMetadata `json:"search_metadata"`
}

// SearchMetadata - 搜索元數據
type SearchMetadata struct {
	TextQuery    string `json:"text_query"`
	SearchTimeMs int64  `json:"search_time_ms"`
}

// ErrorResponse - 錯誤響應
type ErrorResponse struct {
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthCheckResponse - 健康檢查響應
type HealthCheckResponse struct {
	Status       string            `json:"status"`
	Timestamp    time.Time         `json:"timestamp"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}
