package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Leon180/tabelogo-v2/pkg/config"
	"go.uber.org/zap"
)

// GooglePlacesClient handles communication with Google Places API (New)
type GooglePlacesClient struct {
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewGooglePlacesClient creates a new Google Places API client
func NewGooglePlacesClient(cfg *config.Config, logger *zap.Logger) *GooglePlacesClient {
	return &GooglePlacesClient{
		apiKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger.With(zap.String("component", "google_places_client")),
	}
}

// GetPlaceDetails retrieves place details from Google Places API
func (c *GooglePlacesClient) GetPlaceDetails(
	ctx context.Context,
	placeID string,
	languageCode string,
	fieldMask string,
) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://places.googleapis.com/v1/places/%s", placeID)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	// IMPORTANT: Always use English for addressComponents to ensure area is in English
	// Even if client requests Japanese, we need English for area extraction
	q.Add("languageCode", "en")
	req.URL.RawQuery = q.Encode()

	// Add headers
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	if fieldMask != "" {
		req.Header.Set("X-Goog-FieldMask", fieldMask)
	} else {
		// Default field mask - includes addressComponents for area extraction
		req.Header.Set("X-Goog-FieldMask", "id,displayName,formattedAddress,location,rating,priceLevel,photos,currentOpeningHours,addressComponents")
	}

	c.logger.Info("Calling Google Places API",
		zap.String("place_id", placeID),
		zap.String("language", languageCode),
	)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Google API request failed", zap.Error(err))
		return nil, fmt.Errorf("google api request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		c.logger.Warn("Google API returned non-200 status",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return nil, fmt.Errorf("google api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error("Failed to parse response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	c.logger.Info("Successfully retrieved place details",
		zap.String("place_id", placeID),
	)

	return result, nil
}

// TextSearch performs a text-based search for places
func (c *GooglePlacesClient) TextSearch(
	ctx context.Context,
	textQuery string,
	locationBias map[string]interface{},
	maxResultCount int,
	rankPreference string,
	languageCode string,
	fieldMask string,
) (map[string]interface{}, error) {
	url := "https://places.googleapis.com/v1/places:searchText"

	// Build request body
	requestBody := map[string]interface{}{
		"textQuery":      textQuery,
		"locationBias":   locationBias,
		"maxResultCount": maxResultCount,
		"rankPreference": rankPreference,
		"languageCode":   languageCode,
	}

	// Marshal request body
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.logger.Error("Failed to marshal request", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	if fieldMask != "" {
		req.Header.Set("X-Goog-FieldMask", fieldMask)
	} else {
		// Default field mask - includes addressComponents for area extraction
		req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.formattedAddress,places.location,places.rating,places.priceLevel,places.currentOpeningHours,places.addressComponents")
	}

	c.logger.Info("Calling Google Text Search API",
		zap.String("text_query", textQuery),
		zap.String("language", languageCode),
	)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Text search request failed", zap.Error(err))
		return nil, fmt.Errorf("text search request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		c.logger.Warn("Google Text Search API returned non-200 status",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return nil, fmt.Errorf("google text search error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error("Failed to parse response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Log result count
	places, _ := result["places"].([]interface{})
	c.logger.Info("Successfully retrieved text search results",
		zap.String("text_query", textQuery),
		zap.Int("result_count", len(places)),
	)

	return result, nil
}
