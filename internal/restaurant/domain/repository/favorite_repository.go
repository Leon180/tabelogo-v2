package repository

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/google/uuid"
)

// FavoriteRepository defines the interface for favorite persistence
type FavoriteRepository interface {
	// Create creates a new favorite
	Create(ctx context.Context, favorite *model.Favorite) error

	// FindByID finds a favorite by ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Favorite, error)

	// FindByUserAndRestaurant finds a favorite by user ID and restaurant ID
	FindByUserAndRestaurant(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error)

	// FindByUserID finds all favorites for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Favorite, error)

	// FindByRestaurantID finds all favorites for a restaurant
	FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]*model.Favorite, error)

	// Update updates an existing favorite
	Update(ctx context.Context, favorite *model.Favorite) error

	// Delete soft-deletes a favorite by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// Exists checks if a favorite exists for a user and restaurant
	Exists(ctx context.Context, userID, restaurantID uuid.UUID) (bool, error)

	// Count returns the total count of favorites for a user
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

	// FindByTag finds favorites by tag for a user
	FindByTag(ctx context.Context, userID uuid.UUID, tag string) ([]*model.Favorite, error)
}
