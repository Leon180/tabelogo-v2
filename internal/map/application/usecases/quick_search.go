package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/map/domain/models"
	"github.com/Leon180/tabelogo-v2/internal/map/infrastructure/cache"
	"github.com/Leon180/tabelogo-v2/internal/map/infrastructure/external"
	"github.com/Leon180/tabelogo-v2/pkg/metrics"
	"go.uber.org/zap"
)

// QuickSearchUseCase handles quick search business logic
type QuickSearchUseCase struct {
	placesClient *external.GooglePlacesClient
	cache        *cache.PlaceCache
	logger       *zap.Logger
	cacheTTL     time.Duration
}

// NewQuickSearchUseCase creates a new QuickSearchUseCase
func NewQuickSearchUseCase(
	placesClient *external.GooglePlacesClient,
	cache *cache.PlaceCache,
	logger *zap.Logger,
) *QuickSearchUseCase {
	return &QuickSearchUseCase{
		placesClient: placesClient,
		cache:        cache,
		logger:       logger.With(zap.String("usecase", "quick_search")),
		cacheTTL:     1 * time.Hour, // Default 1 hour TTL
	}
}

// Execute performs quick search with cache-first strategy
func (uc *QuickSearchUseCase) Execute(
	ctx context.Context,
	req *models.QuickSearchRequest,
) (*models.QuickSearchResponse, error) {
	uc.logger.Info("Executing quick search",
		zap.String("place_id", req.PlaceID),
		zap.String("language", req.LanguageCode),
	)

	// 1. Try to get from cache first
	cached, err := uc.cache.GetPlace(ctx, req.PlaceID, req.LanguageCode)
	if err != nil {
		// Log cache error but continue to API
		uc.logger.Warn("Cache get failed, falling back to API",
			zap.Error(err),
		)
	}

	if cached != nil {
		// Cache hit
		metrics.CacheHitsTotal.Inc()
		uc.logger.Info("Returning cached result",
			zap.String("place_id", req.PlaceID),
			zap.Time("cached_at", cached.CachedAt),
		)
		return &models.QuickSearchResponse{
			Source:   "redis",
			CachedAt: &cached.CachedAt,
			Result:   cached.Data,
		}, nil
	}

	// Cache miss
	metrics.CacheMissesTotal.Inc()

	// 2. Cache miss - call Google API
	uc.logger.Info("Cache miss, calling Google API",
		zap.String("place_id", req.PlaceID),
	)

	// Record API call metrics
	apiStart := time.Now()
	placeData, err := uc.placesClient.GetPlaceDetails(
		ctx,
		req.PlaceID,
		req.LanguageCode,
		req.APIMask,
	)
	apiDuration := time.Since(apiStart).Seconds()

	if err != nil {
		metrics.GoogleAPICallsTotal.WithLabelValues("place_details", "error").Inc()
		metrics.ErrorsTotal.WithLabelValues("google_api").Inc()
		uc.logger.Error("Google API call failed",
			zap.Error(err),
			zap.String("place_id", req.PlaceID),
		)
		return nil, fmt.Errorf("failed to get place details: %w", err)
	}

	metrics.GoogleAPICallsTotal.WithLabelValues("place_details", "success").Inc()
	metrics.GoogleAPIDuration.WithLabelValues("place_details").Observe(apiDuration)

	// 3. Store in cache (best effort - don't fail if cache write fails)
	if err := uc.cache.SetPlace(ctx, req.PlaceID, req.LanguageCode, placeData, uc.cacheTTL); err != nil {
		uc.logger.Warn("Failed to cache result",
			zap.Error(err),
			zap.String("place_id", req.PlaceID),
		)
		// Continue anyway - we have the data
	}

	// 4. Return result
	return &models.QuickSearchResponse{
		Source:   "google",
		CachedAt: nil,
		Result:   placeData,
	}, nil
}
