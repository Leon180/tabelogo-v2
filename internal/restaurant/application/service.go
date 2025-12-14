package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/application/converters"
	domainerrors "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/repository"
	"github.com/Leon180/tabelogo-v2/pkg/metrics"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RestaurantService defines the application service interface
type RestaurantService interface {
	// Restaurant operations
	CreateRestaurant(ctx context.Context, req CreateRestaurantRequest) (*model.Restaurant, error)
	GetRestaurant(ctx context.Context, id uuid.UUID) (*model.Restaurant, error)
	GetRestaurantByExternalID(ctx context.Context, source model.RestaurantSource, externalID string) (*model.Restaurant, error)
	UpdateRestaurant(ctx context.Context, id uuid.UUID, req UpdateRestaurantRequest) (*model.Restaurant, error)
	DeleteRestaurant(ctx context.Context, id uuid.UUID) error
	SearchRestaurants(ctx context.Context, query string, limit, offset int) ([]*model.Restaurant, error)
	ListRestaurants(ctx context.Context, limit, offset int) ([]*model.Restaurant, error)
	FindRestaurantsByLocation(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*model.Restaurant, error)
	FindRestaurantsByCuisineType(ctx context.Context, cuisineType string, limit, offset int) ([]*model.Restaurant, error)
	IncrementRestaurantViewCount(ctx context.Context, id uuid.UUID) error

	// Map Service integration - Quick search by Google Place ID
	QuickSearchByPlaceID(ctx context.Context, placeID string) (*model.Restaurant, error)

	// Favorite operations
	AddToFavorites(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error)
	RemoveFromFavorites(ctx context.Context, userID, restaurantID uuid.UUID) error
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]*model.Favorite, error)
	GetFavoriteByUserAndRestaurant(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error)
	UpdateFavoriteNotes(ctx context.Context, userID, restaurantID uuid.UUID, notes string) (*model.Favorite, error)
	AddFavoriteTag(ctx context.Context, userID, restaurantID uuid.UUID, tag string) (*model.Favorite, error)
	RemoveFavoriteTag(ctx context.Context, userID, restaurantID uuid.UUID, tag string) (*model.Favorite, error)
	AddFavoriteVisit(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error)
	IsFavorite(ctx context.Context, userID, restaurantID uuid.UUID) (bool, error)
}

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
	favoriteRepo   repository.FavoriteRepository
	mapClient      MapServiceClient // Map Service gRPC client
	config         *Config          // Configuration for TTL, etc.
	logger         *zap.Logger
}

// MapServiceClient defines the interface for Map Service gRPC client
// Config holds service configuration
type Config struct {
	DataFreshnessTTL time.Duration
}

// NewRestaurantService creates a new restaurant service
func NewRestaurantService(
	restaurantRepo repository.RestaurantRepository,
	favoriteRepo repository.FavoriteRepository,
	mapClient MapServiceClient,
	config *Config,
	logger *zap.Logger,
) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
		favoriteRepo:   favoriteRepo,
		mapClient:      mapClient,
		config:         config,
		logger:         logger,
	}
}

// Restaurant operations

func (s *restaurantService) CreateRestaurant(ctx context.Context, req CreateRestaurantRequest) (*model.Restaurant, error) {
	// Check if restaurant already exists
	existing, err := s.restaurantRepo.FindByExternalID(ctx, req.Source, req.ExternalID)
	if err == nil && existing != nil {
		return nil, domainerrors.ErrRestaurantAlreadyExists
	}

	// Create location
	location, err := model.NewLocation(req.Latitude, req.Longitude)
	if err != nil {
		return nil, domainerrors.ErrInvalidLocation
	}

	// Create restaurant
	restaurant := model.NewRestaurant(
		req.Name,
		req.Area,
		req.Source,
		req.ExternalID,
		req.Address,
		location,
	)

	// Set Japanese name if provided
	if req.NameJa != "" {
		restaurant.UpdateNameJa(req.NameJa)
	}

	// Set additional details if provided
	if req.Rating > 0 {
		restaurant.UpdateRating(req.Rating)
	}
	if req.PriceRange != "" {
		restaurant.UpdatePriceRange(req.PriceRange)
	}
	if req.CuisineType != "" {
		restaurant.UpdateCuisineType(req.CuisineType)
	}
	if req.Phone != "" {
		restaurant.UpdatePhone(req.Phone)
	}
	if req.Website != "" {
		restaurant.UpdateWebsite(req.Website)
	}

	// Save to repository
	if err := s.restaurantRepo.Create(ctx, restaurant); err != nil {
		s.logger.Error("Failed to create restaurant", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Restaurant created successfully", zap.String("id", restaurant.ID().String()))
	return restaurant, nil
}

func (s *restaurantService) GetRestaurant(ctx context.Context, id uuid.UUID) (*model.Restaurant, error) {
	restaurant, err := s.restaurantRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get restaurant", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return restaurant, nil
}

func (s *restaurantService) GetRestaurantByExternalID(ctx context.Context, source model.RestaurantSource, externalID string) (*model.Restaurant, error) {
	restaurant, err := s.restaurantRepo.FindByExternalID(ctx, source, externalID)
	if err != nil {
		s.logger.Error("Failed to get restaurant by external ID",
			zap.String("source", string(source)),
			zap.String("externalID", externalID),
			zap.Error(err))
		return nil, err
	}

	return restaurant, nil
}

func (s *restaurantService) UpdateRestaurant(ctx context.Context, id uuid.UUID, req UpdateRestaurantRequest) (*model.Restaurant, error) {
	// Get existing restaurant
	restaurant, err := s.restaurantRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to find restaurant for update", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	// Update area if provided
	if req.Area != "" {
		restaurant.UpdateArea(req.Area)
		s.logger.Info("Updating restaurant area",
			zap.String("id", id.String()),
			zap.String("area", req.Area))
	}

	// Update Japanese name if provided
	if req.NameJa != "" {
		restaurant.UpdateNameJa(req.NameJa)
		s.logger.Info("Updating restaurant Japanese name",
			zap.String("id", id.String()),
			zap.String("name_ja", req.NameJa))
	}
	if req.Address != "" {
		restaurant.UpdateAddress(req.Address)
		s.logger.Info("Updating restaurant address",
			zap.String("id", id.String()),
			zap.String("address", req.Address))
	}
	if req.Rating > 0 {
		restaurant.UpdateRating(req.Rating)
		s.logger.Info("Updating restaurant rating",
			zap.String("id", id.String()),
			zap.Float64("rating", req.Rating))
	}
	if req.PriceRange != "" {
		restaurant.UpdatePriceRange(req.PriceRange)
		s.logger.Info("Updating restaurant price range",
			zap.String("id", id.String()),
			zap.String("price_range", req.PriceRange))
	}
	if req.CuisineType != "" {
		restaurant.UpdateCuisineType(req.CuisineType)
		s.logger.Info("Updating restaurant cuisine type",
			zap.String("id", id.String()),
			zap.String("cuisine_type", req.CuisineType))
	}
	if req.Phone != "" {
		restaurant.UpdatePhone(req.Phone)
		s.logger.Info("Updating restaurant phone",
			zap.String("id", id.String()),
			zap.String("phone", req.Phone))
	}
	if req.Website != "" {
		restaurant.UpdateWebsite(req.Website)
		s.logger.Info("Updating restaurant website",
			zap.String("id", id.String()),
			zap.String("website", req.Website))
	}
	if req.Latitude != 0 && req.Longitude != 0 {
		location, err := model.NewLocation(req.Latitude, req.Longitude)
		if err != nil {
			return nil, err
		}
		restaurant.UpdateLocation(location)
		s.logger.Info("Updating restaurant location",
			zap.String("id", id.String()),
			zap.Float64("latitude", req.Latitude),
			zap.Float64("longitude", req.Longitude))
	}

	// Save updated restaurant
	if err := s.restaurantRepo.Update(ctx, restaurant); err != nil {
		s.logger.Error("Failed to update restaurant", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Restaurant updated successfully", zap.String("id", id.String()))
	return restaurant, nil
}

func (s *restaurantService) DeleteRestaurant(ctx context.Context, id uuid.UUID) error {
	if err := s.restaurantRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete restaurant", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.logger.Info("Restaurant deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *restaurantService) SearchRestaurants(ctx context.Context, query string, limit, offset int) ([]*model.Restaurant, error) {
	restaurants, err := s.restaurantRepo.Search(ctx, query, limit, offset)
	if err != nil {
		s.logger.Error("Failed to search restaurants", zap.String("query", query), zap.Error(err))
		return nil, err
	}

	return restaurants, nil
}

func (s *restaurantService) ListRestaurants(ctx context.Context, limit, offset int) ([]*model.Restaurant, error) {
	restaurants, err := s.restaurantRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list restaurants", zap.Error(err))
		return nil, err
	}

	return restaurants, nil
}

func (s *restaurantService) FindRestaurantsByLocation(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*model.Restaurant, error) {
	restaurants, err := s.restaurantRepo.FindByLocation(ctx, lat, lng, radiusKm, limit)
	if err != nil {
		s.logger.Error("Failed to find restaurants by location", zap.Error(err))
		return nil, err
	}

	return restaurants, nil
}

func (s *restaurantService) FindRestaurantsByCuisineType(ctx context.Context, cuisineType string, limit, offset int) ([]*model.Restaurant, error) {
	restaurants, err := s.restaurantRepo.FindByCuisineType(ctx, cuisineType, limit, offset)
	if err != nil {
		s.logger.Error("Failed to find restaurants by cuisine type", zap.String("cuisineType", cuisineType), zap.Error(err))
		return nil, err
	}

	return restaurants, nil
}

func (s *restaurantService) IncrementRestaurantViewCount(ctx context.Context, id uuid.UUID) error {
	restaurant, err := s.restaurantRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	restaurant.IncrementViewCount()

	if err := s.restaurantRepo.Update(ctx, restaurant); err != nil {
		s.logger.Error("Failed to increment view count", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	return nil
}

// QuickSearchByPlaceID implements cache-first search by Google Place ID
// 1. Check local DB first (cache hit)
// 2. If not found or stale, call Map Service
// 3. Save/update result in local DB
// 4. Return restaurant
// QuickSearchByPlaceID searches for a restaurant by Google Place ID with cache-first strategy
func (s *restaurantService) QuickSearchByPlaceID(ctx context.Context, placeID string) (*model.Restaurant, error) {
	s.logger.Info("[QuickSearch] START - Request received",
		zap.String("place_id", placeID),
	)

	// Step 1: Try cache first
	s.logger.Info("[QuickSearch] STEP 1 - Checking cache",
		zap.String("place_id", placeID),
	)
	restaurant, err := s.restaurantRepo.FindByExternalID(ctx, model.SourceGoogle, placeID)
	s.logger.Info("[QuickSearch] Cache query result",
		zap.String("place_id", placeID),
		zap.Bool("found", restaurant != nil),
		zap.Bool("has_error", err != nil),
	)
	if err != nil && err != domainerrors.ErrRestaurantNotFound {
		s.logger.Error("[QuickSearch] ERROR - Failed to query cache",
			zap.String("place_id", placeID),
			zap.Error(err),
		)
		return nil, err
	}

	// Check if we have fresh data
	if err == nil && restaurant != nil {
		// Check data freshness
		if time.Since(restaurant.UpdatedAt()) < s.config.DataFreshnessTTL {
			metrics.RestaurantCacheHitsTotal.Inc() // METRIC: Cache hit
			s.logger.Info("[QuickSearch] Cache hit - returning fresh data",
				zap.String("place_id", placeID),
				zap.Duration("age", time.Since(restaurant.UpdatedAt())),
			)
			return restaurant, nil
		}
		s.logger.Info("[QuickSearch] Cache hit but data is stale, refreshing from Map Service",
			zap.String("place_id", placeID),
			zap.Duration("age", time.Since(restaurant.UpdatedAt())),
		)
	} else {
		metrics.RestaurantCacheMissesTotal.Inc() // METRIC: Cache miss
		s.logger.Info("[QuickSearch] Cache miss - fetching from Map Service",
			zap.String("place_id", placeID),
		)
	}

	// Step 2: Cache miss or stale data - call Map Service
	s.logger.Info("[QuickSearch] STEP 2 - Calling Map Service",
		zap.String("place_id", placeID),
	)
	start := time.Now()
	place, err := s.mapClient.QuickSearch(ctx, placeID)
	s.logger.Info("[QuickSearch] Map Service response",
		zap.String("place_id", placeID),
		zap.Bool("success", err == nil),
		zap.Bool("has_place", place != nil),
		zap.Duration("duration", time.Since(start)),
	)
	if err != nil {
		metrics.RestaurantMapServiceCallsTotal.WithLabelValues("error").Inc() // METRIC: Map Service error
		s.logger.Error("[QuickSearch] ERROR - Map Service call failed",
			zap.String("place_id", placeID),
			zap.Error(err),
		)

		// If Map Service fails and we have stale data, return it with a warning
		if restaurant != nil {
			metrics.RestaurantStaleDataReturnsTotal.Inc() // METRIC: Stale data fallback
			s.logger.Warn("[QuickSearch] Map Service failed, returning stale data",
				zap.String("place_id", placeID),
				zap.Error(err),
			)
			return restaurant, nil
		}
		// No cached data and Map Service failed
		s.logger.Error("[QuickSearch] FATAL - Map Service failed and no cached data available",
			zap.String("place_id", placeID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("map service quick search failed: %w", err)
	}

	metrics.RestaurantMapServiceCallsTotal.WithLabelValues("success").Inc() // METRIC: Map Service success
	metrics.RestaurantSyncDuration.Observe(time.Since(start).Seconds())     // METRIC: Sync duration

	// Step 3: Convert Map Service proto to domain model
	s.logger.Info("[QuickSearch] STEP 3 - Converting Map Service response to domain model",
		zap.String("place_id", placeID),
	)
	newRestaurant := converters.MapPlaceToRestaurant(place)
	if newRestaurant == nil {
		s.logger.Warn("[QuickSearch] Conversion from MapPlace to Restaurant returned nil",
			zap.String("place_id", placeID),
		)
		return nil, domainerrors.ErrRestaurantNotFound
	}

	// Step 4: Save or update in local DB
	if restaurant == nil {
		// Create new restaurant
		if err := s.restaurantRepo.Create(ctx, newRestaurant); err != nil {
			s.logger.Error("Failed to save restaurant from Map Service",
				zap.String("place_id", placeID),
				zap.Error(err),
			)
			// Return the data anyway, even if save failed
			return newRestaurant, nil
		}
		s.logger.Info("Saved new restaurant from Map Service",
			zap.String("place_id", placeID),
			zap.String("name", newRestaurant.Name()),
		)
	} else {
		// Update existing restaurant with fresh data
		restaurant.UpdateDetails(
			newRestaurant.Name(),
			newRestaurant.Address(),
			newRestaurant.PriceRange(),
			newRestaurant.CuisineType(),
			newRestaurant.Phone(),
			newRestaurant.Website(),
		)
		restaurant.UpdateRating(newRestaurant.Rating())
		if newRestaurant.Location() != nil {
			restaurant.UpdateLocation(newRestaurant.Location())
		}

		if err := s.restaurantRepo.Update(ctx, restaurant); err != nil {
			s.logger.Error("Failed to update restaurant from Map Service",
				zap.String("place_id", placeID),
				zap.Error(err),
			)
		} else {
			s.logger.Info("Updated restaurant from Map Service",
				zap.String("place_id", placeID),
				zap.String("name", restaurant.Name()),
			)
		}
		return restaurant, nil
	}

	return newRestaurant, nil
}

// Favorite operations

func (s *restaurantService) AddToFavorites(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error) {
	// Check if restaurant exists
	_, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	// Check if already in favorites
	exists, err := s.favoriteRepo.Exists(ctx, userID, restaurantID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrFavoriteAlreadyExists
	}

	// Create favorite
	favorite := model.NewFavorite(userID, restaurantID)

	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		s.logger.Error("Failed to add to favorites", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Added to favorites", zap.String("favoriteID", favorite.ID().String()))
	return favorite, nil
}

func (s *restaurantService) RemoveFromFavorites(ctx context.Context, userID, restaurantID uuid.UUID) error {
	favorite, err := s.favoriteRepo.FindByUserAndRestaurant(ctx, userID, restaurantID)
	if err != nil {
		return err
	}

	if err := s.favoriteRepo.Delete(ctx, favorite.ID()); err != nil {
		s.logger.Error("Failed to remove from favorites", zap.Error(err))
		return err
	}

	s.logger.Info("Removed from favorites", zap.String("favoriteID", favorite.ID().String()))
	return nil
}

func (s *restaurantService) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]*model.Favorite, error) {
	favorites, err := s.favoriteRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user favorites", zap.Error(err))
		return nil, err
	}

	return favorites, nil
}

func (s *restaurantService) GetFavoriteByUserAndRestaurant(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error) {
	favorite, err := s.favoriteRepo.FindByUserAndRestaurant(ctx, userID, restaurantID)
	if err != nil {
		return nil, err
	}

	return favorite, nil
}

func (s *restaurantService) UpdateFavoriteNotes(ctx context.Context, userID, restaurantID uuid.UUID, notes string) (*model.Favorite, error) {
	favorite, err := s.favoriteRepo.FindByUserAndRestaurant(ctx, userID, restaurantID)
	if err != nil {
		return nil, err
	}

	favorite.UpdateNotes(notes)

	if err := s.favoriteRepo.Update(ctx, favorite); err != nil {
		s.logger.Error("Failed to update favorite notes", zap.Error(err))
		return nil, err
	}

	return favorite, nil
}

func (s *restaurantService) AddFavoriteTag(ctx context.Context, userID, restaurantID uuid.UUID, tag string) (*model.Favorite, error) {
	favorite, err := s.favoriteRepo.FindByUserAndRestaurant(ctx, userID, restaurantID)
	if err != nil {
		return nil, err
	}

	favorite.AddTag(tag)

	if err := s.favoriteRepo.Update(ctx, favorite); err != nil {
		s.logger.Error("Failed to add favorite tag", zap.Error(err))
		return nil, err
	}

	return favorite, nil
}

func (s *restaurantService) RemoveFavoriteTag(ctx context.Context, userID, restaurantID uuid.UUID, tag string) (*model.Favorite, error) {
	favorite, err := s.favoriteRepo.FindByUserAndRestaurant(ctx, userID, restaurantID)
	if err != nil {
		return nil, err
	}

	favorite.RemoveTag(tag)

	if err := s.favoriteRepo.Update(ctx, favorite); err != nil {
		s.logger.Error("Failed to remove favorite tag", zap.Error(err))
		return nil, err
	}

	return favorite, nil
}

func (s *restaurantService) AddFavoriteVisit(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error) {
	favorite, err := s.favoriteRepo.FindByUserAndRestaurant(ctx, userID, restaurantID)
	if err != nil {
		return nil, err
	}

	favorite.AddVisit()

	if err := s.favoriteRepo.Update(ctx, favorite); err != nil {
		s.logger.Error("Failed to add favorite visit", zap.Error(err))
		return nil, err
	}

	return favorite, nil
}

func (s *restaurantService) IsFavorite(ctx context.Context, userID, restaurantID uuid.UUID) (bool, error) {
	exists, err := s.favoriteRepo.Exists(ctx, userID, restaurantID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
