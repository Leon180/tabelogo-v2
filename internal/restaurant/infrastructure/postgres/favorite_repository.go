package postgres

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/repository"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// FavoriteORM is the database model for Favorite
type FavoriteORM struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	RestaurantID  uuid.UUID      `gorm:"type:uuid;not null;index"`
	Notes         string         `gorm:"type:text"`
	Tags          pq.StringArray `gorm:"type:varchar(255)[]"`
	VisitCount    int            `gorm:"type:int;default:0"`
	LastVisitedAt *time.Time     `gorm:"type:timestamp"`
	CreatedAt     time.Time      `gorm:"not null"`
	UpdatedAt     time.Time      `gorm:"not null"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name
func (FavoriteORM) TableName() string {
	return "user_favorites"
}

// ToDomainFavorite converts ORM model to Domain entity
func (f *FavoriteORM) ToDomain() *model.Favorite {
	// Convert pq.StringArray to []string
	tags := make([]string, len(f.Tags))
	copy(tags, f.Tags)

	// Handle DeletedAt
	var deletedAt *time.Time
	if f.DeletedAt.Valid {
		deletedAt = &f.DeletedAt.Time
	}

	return model.ReconstructFavorite(
		f.ID,
		f.UserID,
		f.RestaurantID,
		f.Notes,
		tags,
		f.VisitCount,
		f.LastVisitedAt,
		f.CreatedAt,
		f.UpdatedAt,
		deletedAt,
	)
}

// FromDomainFavorite converts Domain entity to ORM model
func FromDomainFavorite(f *model.Favorite) *FavoriteORM {
	// Convert []string to pq.StringArray
	tags := pq.StringArray(f.Tags())

	// Handle DeletedAt
	var deletedAt gorm.DeletedAt
	if f.DeletedAt() != nil {
		deletedAt = gorm.DeletedAt{Time: *f.DeletedAt(), Valid: true}
	}

	return &FavoriteORM{
		ID:            f.ID(),
		UserID:        f.UserID(),
		RestaurantID:  f.RestaurantID(),
		Notes:         f.Notes(),
		Tags:          tags,
		VisitCount:    f.VisitCount(),
		LastVisitedAt: f.LastVisitedAt(),
		CreatedAt:     f.CreatedAt(),
		UpdatedAt:     f.UpdatedAt(),
		DeletedAt:     deletedAt,
	}
}

type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository creates a new postgres favorite repository
func NewFavoriteRepository(db *gorm.DB) repository.FavoriteRepository {
	return &favoriteRepository{db: db}
}

// Create creates a new favorite
func (r *favoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	orm := FromDomainFavorite(favorite)

	if err := r.db.WithContext(ctx).Create(orm).Error; err != nil {
		return err
	}

	return nil
}

// FindByID finds a favorite by ID
func (r *favoriteRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Favorite, error) {
	var orm FavoriteORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrFavoriteNotFound
		}
		return nil, err
	}

	return orm.ToDomain(), nil
}

// FindByUserAndRestaurant finds a favorite by user ID and restaurant ID
func (r *favoriteRepository) FindByUserAndRestaurant(ctx context.Context, userID, restaurantID uuid.UUID) (*model.Favorite, error) {
	var orm FavoriteORM
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND restaurant_id = ?", userID, restaurantID).
		First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrFavoriteNotFound
		}
		return nil, err
	}

	return orm.ToDomain(), nil
}

// FindByUserID finds all favorites for a user
func (r *favoriteRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Favorite, error) {
	var orms []FavoriteORM
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orms).Error; err != nil {
		return nil, err
	}

	favorites := make([]*model.Favorite, len(orms))
	for i, orm := range orms {
		favorites[i] = orm.ToDomain()
	}

	return favorites, nil
}

// FindByRestaurantID finds all favorites for a restaurant
func (r *favoriteRepository) FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]*model.Favorite, error) {
	var orms []FavoriteORM
	if err := r.db.WithContext(ctx).
		Where("restaurant_id = ?", restaurantID).
		Find(&orms).Error; err != nil {
		return nil, err
	}

	favorites := make([]*model.Favorite, len(orms))
	for i, orm := range orms {
		favorites[i] = orm.ToDomain()
	}

	return favorites, nil
}

// Update updates an existing favorite
func (r *favoriteRepository) Update(ctx context.Context, favorite *model.Favorite) error {
	orm := FromDomainFavorite(favorite)

	result := r.db.WithContext(ctx).Model(&FavoriteORM{}).Where("id = ?", orm.ID).Updates(orm)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrFavoriteNotFound
	}

	return nil
}

// Delete soft-deletes a favorite by ID
func (r *favoriteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&FavoriteORM{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrFavoriteNotFound
	}

	return nil
}

// Exists checks if a favorite exists for a user and restaurant
func (r *favoriteRepository) Exists(ctx context.Context, userID, restaurantID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&FavoriteORM{}).
		Where("user_id = ? AND restaurant_id = ?", userID, restaurantID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CountByUserID returns the total count of favorites for a user
func (r *favoriteRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&FavoriteORM{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// FindByTag finds favorites by tag for a user
func (r *favoriteRepository) FindByTag(ctx context.Context, userID uuid.UUID, tag string) ([]*model.Favorite, error) {
	var orms []FavoriteORM
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND ? = ANY(tags)", userID, tag).
		Find(&orms).Error; err != nil {
		return nil, err
	}

	favorites := make([]*model.Favorite, len(orms))
	for i, orm := range orms {
		favorites[i] = orm.ToDomain()
	}

	return favorites, nil
}
