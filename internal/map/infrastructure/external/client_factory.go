package external

import (
	"context"
	"os"
	"strings"

	"github.com/Leon180/tabelogo-v2/pkg/config"
	"go.uber.org/zap"
)

// PlacesClient is an interface for Places API clients (Google or Mock)
type PlacesClient interface {
	GetPlaceDetails(
		ctx context.Context,
		placeID string,
		languageCode string,
		fieldMask string,
	) (map[string]interface{}, error)

	TextSearch(
		ctx context.Context,
		textQuery string,
		locationBias map[string]interface{},
		maxResultCount int,
		rankPreference string,
		languageCode string,
		fieldMask string,
	) (map[string]interface{}, error)
}

// NewPlacesClient creates a Places API client based on configuration
// If USE_MOCK_API=true, returns MockPlacesClient
// Otherwise, returns GooglePlacesClient
func NewPlacesClient(cfg *config.Config, logger *zap.Logger) PlacesClient {
	useMockAPI := strings.ToLower(os.Getenv("USE_MOCK_API")) == "true"
	mockBaseURL := os.Getenv("MOCK_API_BASE_URL")

	if useMockAPI && mockBaseURL != "" {
		logger.Info("Using Mock Places API",
			zap.String("base_url", mockBaseURL),
		)
		return NewMockPlacesClient(mockBaseURL, logger)
	}

	logger.Info("Using Google Places API")
	return NewGooglePlacesClient(cfg, logger)
}
