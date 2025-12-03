package repository

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/google/uuid"
)

// RestaurantRepository defines the interface for restaurant persistence
type RestaurantRepository interface {
	// Create creates a new restaurant
	Create(ctx context.Context, restaurant *model.Restaurant) error

	// FindByID finds a restaurant by ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Restaurant, error)

	// FindByExternalID finds a restaurant by source and external ID
	FindByExternalID(ctx context.Context, source model.RestaurantSource, externalID string) (*model.Restaurant, error)

	// Update updates an existing restaurant
	Update(ctx context.Context, restaurant *model.Restaurant) error

	// Delete soft-deletes a restaurant by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// Search searches restaurants by query string
	Search(ctx context.Context, query string, limit, offset int) ([]*model.Restaurant, error)

	// FindByLocation finds restaurants within a radius from a location
	FindByLocation(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*model.Restaurant, error)

	// List lists all restaurants with pagination
	List(ctx context.Context, limit, offset int) ([]*model.Restaurant, error)

	// Count returns the total count of restaurants
	Count(ctx context.Context) (int64, error)

	// FindByCuisineType finds restaurants by cuisine type
	FindByCuisineType(ctx context.Context, cuisineType string, limit, offset int) ([]*model.Restaurant, error)

	// FindBySource finds restaurants by source
	FindBySource(ctx context.Context, source model.RestaurantSource, limit, offset int) ([]*model.Restaurant, error)
}
