package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/map/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/map/infrastructure/external"
	"github.com/Leon180/tabelogo-v2/pkg/metrics"
	"go.uber.org/zap"
)

// AdvanceSearchUseCase handles advance search business logic
type AdvanceSearchUseCase struct {
	placesClient external.PlacesClient
	logger       *zap.Logger
}

// NewAdvanceSearchUseCase creates a new AdvanceSearchUseCase
func NewAdvanceSearchUseCase(
	placesClient external.PlacesClient,
	logger *zap.Logger,
) *AdvanceSearchUseCase {
	return &AdvanceSearchUseCase{
		placesClient: placesClient,
		logger:       logger.With(zap.String("usecase", "advance_search")),
	}
}

// Execute performs advance search
func (uc *AdvanceSearchUseCase) Execute(
	ctx context.Context,
	req *models.AdvanceSearchRequest,
) (*models.AdvanceSearchResponse, error) {
	startTime := time.Now()

	uc.logger.Info("Executing advance search",
		zap.String("text_query", req.TextQuery),
		zap.String("language", req.LanguageCode),
		zap.String("rank_preference", req.RankPreference),
	)

	// Build location bias
	locationBias := map[string]interface{}{
		"rectangle": map[string]interface{}{
			"low": map[string]interface{}{
				"latitude":  req.LocationBias.Rectangle.Low.Latitude,
				"longitude": req.LocationBias.Rectangle.Low.Longitude,
			},
			"high": map[string]interface{}{
				"latitude":  req.LocationBias.Rectangle.High.Latitude,
				"longitude": req.LocationBias.Rectangle.High.Longitude,
			},
		},
	}

	// Call Google Text Search API
	apiStart := time.Now()
	result, err := uc.placesClient.TextSearch(
		ctx,
		req.TextQuery,
		locationBias,
		req.MaxResultCount,
		req.RankPreference,
		req.LanguageCode,
		req.APIMask,
	)
	apiDuration := time.Since(apiStart).Seconds()

	if err != nil {
		metrics.GoogleAPICallsTotal.WithLabelValues("text_search", "error").Inc()
		metrics.ErrorsTotal.WithLabelValues("google_api").Inc()
		uc.logger.Error("Text search failed",
			zap.Error(err),
			zap.String("text_query", req.TextQuery),
		)
		return nil, fmt.Errorf("text search failed: %w", err)
	}

	metrics.GoogleAPICallsTotal.WithLabelValues("text_search", "success").Inc()
	metrics.GoogleAPIDuration.WithLabelValues("text_search").Observe(apiDuration)

	// Extract places from result
	places, ok := result["places"].([]interface{})
	if !ok {
		places = []interface{}{}
	}

	// Filter results
	filteredPlaces := uc.filterResults(places, req)

	// Calculate search time
	searchTimeMs := time.Since(startTime).Milliseconds()

	// Build response
	response := &models.AdvanceSearchResponse{
		Places:     filteredPlaces,
		TotalCount: len(filteredPlaces),
		SearchMetadata: models.SearchMetadata{
			TextQuery:    req.TextQuery,
			SearchTimeMs: searchTimeMs,
		},
	}

	uc.logger.Info("Advance search completed",
		zap.Int("total_results", len(filteredPlaces)),
		zap.Int64("search_time_ms", searchTimeMs),
	)

	return response, nil
}

// filterResults filters search results based on criteria
func (uc *AdvanceSearchUseCase) filterResults(
	places []interface{},
	req *models.AdvanceSearchRequest,
) []interface{} {
	filtered := []interface{}{}

	for _, place := range places {
		placeMap, ok := place.(map[string]interface{})
		if !ok {
			continue
		}

		// Filter by minimum rating
		if req.MinRating > 0 {
			rating, ok := placeMap["rating"].(float64)
			if !ok || rating < req.MinRating {
				uc.logger.Debug("Filtered out by min rating",
					zap.Float64("rating", rating),
					zap.Float64("min_rating", req.MinRating),
				)
				continue
			}
		}

		// Filter by open now
		if req.OpenNow {
			// Check if place has currentOpeningHours
			openingHours, ok := placeMap["currentOpeningHours"].(map[string]interface{})
			if !ok {
				uc.logger.Debug("Filtered out: no opening hours data")
				continue
			}
			openNow, ok := openingHours["openNow"].(bool)
			if !ok || !openNow {
				uc.logger.Debug("Filtered out: not open now")
				continue
			}
		}

		filtered = append(filtered, place)
	}

	uc.logger.Info("Filtering completed",
		zap.Int("original_count", len(places)),
		zap.Int("filtered_count", len(filtered)),
	)

	return filtered
}
