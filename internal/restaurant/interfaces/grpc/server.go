package grpc

import (
	"context"

	restaurantv1 "github.com/Leon180/tabelogo-v2/api/gen/restaurant/v1"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/application"
	domainerrors "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RestaurantServer struct {
	restaurantv1.UnimplementedRestaurantServiceServer
	service application.RestaurantService
	logger  *zap.Logger
}

func NewRestaurantServer(
	service application.RestaurantService,
	logger *zap.Logger,
) *RestaurantServer {
	return &RestaurantServer{
		service: service,
		logger:  logger,
	}
}

// CreateRestaurant creates a new restaurant
func (s *RestaurantServer) CreateRestaurant(
	ctx context.Context,
	req *restaurantv1.CreateRestaurantRequest,
) (*restaurantv1.CreateRestaurantResponse, error) {
	location, err := fromProtoLocation(req.Location)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid location")
	}

	appReq := application.CreateRestaurantRequest{
		Name:        req.Name,
		Source:      model.RestaurantSource(req.Source),
		ExternalID:  req.ExternalId,
		Address:     req.Address,
		Latitude:    location.Latitude(),
		Longitude:   location.Longitude(),
		Rating:      req.Rating,
		PriceRange:  req.PriceRange,
		CuisineType: req.CuisineType,
		Phone:       req.Phone,
		Website:     req.Website,
	}

	restaurant, err := s.service.CreateRestaurant(ctx, appReq)
	if err != nil {
		if err == domainerrors.ErrRestaurantAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "restaurant already exists")
		}
		s.logger.Error("Failed to create restaurant", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create restaurant")
	}

	return &restaurantv1.CreateRestaurantResponse{
		Restaurant: toProtoRestaurant(restaurant),
	}, nil
}

// GetRestaurant retrieves a restaurant by ID
func (s *RestaurantServer) GetRestaurant(
	ctx context.Context,
	req *restaurantv1.GetRestaurantRequest,
) (*restaurantv1.GetRestaurantResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid restaurant ID")
	}

	restaurant, err := s.service.GetRestaurant(ctx, id)
	if err != nil {
		if err == domainerrors.ErrRestaurantNotFound {
			return nil, status.Error(codes.NotFound, "restaurant not found")
		}
		s.logger.Error("Failed to get restaurant", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get restaurant")
	}

	return &restaurantv1.GetRestaurantResponse{
		Restaurant: toProtoRestaurant(restaurant),
	}, nil
}

// GetRestaurantByExternalID retrieves a restaurant by external ID
func (s *RestaurantServer) GetRestaurantByExternalID(
	ctx context.Context,
	req *restaurantv1.GetRestaurantByExternalIDRequest,
) (*restaurantv1.GetRestaurantResponse, error) {
	restaurant, err := s.service.GetRestaurantByExternalID(ctx, model.RestaurantSource(req.Source), req.ExternalId)
	if err != nil {
		if err == domainerrors.ErrRestaurantNotFound {
			return nil, status.Error(codes.NotFound, "restaurant not found")
		}
		s.logger.Error("Failed to get restaurant by external ID", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get restaurant")
	}

	return &restaurantv1.GetRestaurantResponse{
		Restaurant: toProtoRestaurant(restaurant),
	}, nil
}

// UpdateRestaurant updates a restaurant
func (s *RestaurantServer) UpdateRestaurant(
	ctx context.Context,
	req *restaurantv1.UpdateRestaurantRequest,
) (*restaurantv1.RestaurantResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid restaurant ID: %v", err)
	}

	// Currently only supports updating Japanese name
	// Future: Add more fields as needed
	appReq := application.UpdateRestaurantRequest{
		NameJa: req.NameJa,
	}

	restaurant, err := s.service.UpdateRestaurant(ctx, id, appReq)
	if err != nil {
		if err == domainerrors.ErrRestaurantNotFound {
			return nil, status.Error(codes.NotFound, "restaurant not found")
		}
		s.logger.Error("Failed to update restaurant", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update restaurant: %v", err)
	}

	return &restaurantv1.RestaurantResponse{
		Restaurant: toProtoRestaurant(restaurant),
	}, nil
}

// DeleteRestaurant deletes a restaurant
func (s *RestaurantServer) DeleteRestaurant(
	ctx context.Context,
	req *restaurantv1.DeleteRestaurantRequest,
) (*restaurantv1.DeleteRestaurantResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid restaurant ID")
	}

	err = s.service.DeleteRestaurant(ctx, id)
	if err != nil {
		if err == domainerrors.ErrRestaurantNotFound {
			return nil, status.Error(codes.NotFound, "restaurant not found")
		}
		s.logger.Error("Failed to delete restaurant", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete restaurant")
	}

	return &restaurantv1.DeleteRestaurantResponse{
		Success: true,
	}, nil
}

// SearchRestaurants searches for restaurants
func (s *RestaurantServer) SearchRestaurants(
	ctx context.Context,
	req *restaurantv1.SearchRestaurantsRequest,
) (*restaurantv1.SearchRestaurantsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	restaurants, err := s.service.SearchRestaurants(ctx, req.Query, limit, offset)
	if err != nil {
		s.logger.Error("Failed to search restaurants", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to search restaurants")
	}

	return &restaurantv1.SearchRestaurantsResponse{
		Restaurants: toProtoRestaurants(restaurants),
		Total:       int32(len(restaurants)),
	}, nil
}

// ListRestaurants lists restaurants with pagination
func (s *RestaurantServer) ListRestaurants(
	ctx context.Context,
	req *restaurantv1.ListRestaurantsRequest,
) (*restaurantv1.ListRestaurantsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	restaurants, err := s.service.ListRestaurants(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list restaurants", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list restaurants")
	}

	return &restaurantv1.ListRestaurantsResponse{
		Restaurants: toProtoRestaurants(restaurants),
		Total:       int32(len(restaurants)),
	}, nil
}

// AddToFavorites adds a restaurant to user's favorites
func (s *RestaurantServer) AddToFavorites(
	ctx context.Context,
	req *restaurantv1.AddToFavoritesRequest,
) (*restaurantv1.AddToFavoritesResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	restaurantID, err := uuid.Parse(req.RestaurantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid restaurant ID")
	}

	favorite, err := s.service.AddToFavorites(ctx, userID, restaurantID)
	if err != nil {
		if err == domainerrors.ErrFavoriteAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "already in favorites")
		}
		s.logger.Error("Failed to add to favorites", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to add to favorites")
	}

	return &restaurantv1.AddToFavoritesResponse{
		Favorite: toProtoFavorite(favorite),
	}, nil
}

// RemoveFromFavorites removes a restaurant from user's favorites
func (s *RestaurantServer) RemoveFromFavorites(
	ctx context.Context,
	req *restaurantv1.RemoveFromFavoritesRequest,
) (*restaurantv1.RemoveFromFavoritesResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	restaurantID, err := uuid.Parse(req.RestaurantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid restaurant ID")
	}

	err = s.service.RemoveFromFavorites(ctx, userID, restaurantID)
	if err != nil {
		if err == domainerrors.ErrFavoriteNotFound {
			return nil, status.Error(codes.NotFound, "favorite not found")
		}
		s.logger.Error("Failed to remove from favorites", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to remove from favorites")
	}

	return &restaurantv1.RemoveFromFavoritesResponse{
		Success: true,
	}, nil
}

// GetUserFavorites retrieves all favorites for a user
func (s *RestaurantServer) GetUserFavorites(
	ctx context.Context,
	req *restaurantv1.GetUserFavoritesRequest,
) (*restaurantv1.GetUserFavoritesResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	favorites, err := s.service.GetUserFavorites(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user favorites", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get favorites")
	}

	return &restaurantv1.GetUserFavoritesResponse{
		Favorites: toProtoFavorites(favorites),
		Total:     int32(len(favorites)),
	}, nil
}
