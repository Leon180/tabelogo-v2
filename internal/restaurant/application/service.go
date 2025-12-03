package application

import (
	"context"

	domainerrors "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/repository"
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
	logger         *zap.Logger
}

// NewRestaurantService creates a new restaurant service
func NewRestaurantService(
	restaurantRepo repository.RestaurantRepository,
	favoriteRepo repository.FavoriteRepository,
	logger *zap.Logger,
) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
		favoriteRepo:   favoriteRepo,
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
		req.Source,
		req.ExternalID,
		req.Address,
		location,
	)

	// Set optional fields
	if req.Rating > 0 {
		restaurant.UpdateRating(req.Rating)
	}
	if req.PriceRange != "" || req.CuisineType != "" || req.Phone != "" || req.Website != "" {
		restaurant.UpdateDetails("", "", req.PriceRange, req.CuisineType, req.Phone, req.Website)
	}

	// Save restaurant
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
		return nil, err
	}

	// Update details
	restaurant.UpdateDetails(
		req.Name,
		req.Address,
		req.PriceRange,
		req.CuisineType,
		req.Phone,
		req.Website,
	)

	// Update rating
	if req.Rating > 0 {
		restaurant.UpdateRating(req.Rating)
	}

	// Update location
	if req.Latitude != 0 && req.Longitude != 0 {
		location, err := model.NewLocation(req.Latitude, req.Longitude)
		if err != nil {
			return nil, domainerrors.ErrInvalidLocation
		}
		restaurant.UpdateLocation(location)
	}

	// Save changes
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
