package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// MockPlacesClient handles communication with Mock Map Service
type MockPlacesClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewMockPlacesClient creates a new Mock Places API client
func NewMockPlacesClient(baseURL string, logger *zap.Logger) *MockPlacesClient {
	return &MockPlacesClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger.With(zap.String("component", "mock_places_client")),
	}
}

// GetPlaceDetails retrieves place details from Mock service
func (c *MockPlacesClient) GetPlaceDetails(
	ctx context.Context,
	placeID string,
	languageCode string,
	fieldMask string,
) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/places/%s", c.baseURL, placeID)

	c.logger.Info("Calling Mock Places API",
		zap.String("place_id", placeID),
		zap.String("url", url),
	)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Mock API request failed", zap.Error(err))
		return nil, fmt.Errorf("mock api request failed: %w", err)
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
		c.logger.Warn("Mock API returned non-200 status",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return nil, fmt.Errorf("mock api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error("Failed to parse response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	c.logger.Info("Successfully retrieved place details from Mock",
		zap.String("place_id", placeID),
	)

	return result, nil
}

// TextSearch performs a text-based search for places using Mock service
func (c *MockPlacesClient) TextSearch(
	ctx context.Context,
	textQuery string,
	locationBias map[string]interface{},
	maxResultCount int,
	rankPreference string,
	languageCode string,
	fieldMask string,
) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/places:searchText", c.baseURL)

	c.logger.Info("Calling Mock Text Search API",
		zap.String("text_query", textQuery),
		zap.String("url", url),
	)

	// Build request body
	requestBody := map[string]interface{}{
		"textQuery": textQuery,
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
		c.logger.Warn("Mock Text Search API returned non-200 status",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return nil, fmt.Errorf("mock text search error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error("Failed to parse response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Log result count
	places, _ := result["places"].([]interface{})
	c.logger.Info("Successfully retrieved text search results from Mock",
		zap.String("text_query", textQuery),
		zap.Int("result_count", len(places)),
	)

	return result, nil
}
