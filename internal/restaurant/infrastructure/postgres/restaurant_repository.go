package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	domainerrors "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RestaurantORM is the database model for Restaurant
type RestaurantORM struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name         string         `gorm:"type:varchar(255);not null"`
	NameJa       string         `gorm:"type:varchar(255)"`
	Area         string         `gorm:"type:varchar(100)"`
	Source       string         `gorm:"type:varchar(50);not null"`
	ExternalID   string         `gorm:"type:varchar(255);not null"`
	Address      string         `gorm:"type:text"`
	Latitude     float64        `gorm:"type:decimal(10,8)"`
	Longitude    float64        `gorm:"type:decimal(11,8)"`
	Rating       float64        `gorm:"type:decimal(3,2)"`
	PriceRange   string         `gorm:"type:varchar(10)"`
	CuisineType  string         `gorm:"type:varchar(50)"`
	Phone        string         `gorm:"type:varchar(20)"`
	Website      string         `gorm:"type:varchar(500)"`
	OpeningHours string         `gorm:"type:jsonb"` // JSON string
	Metadata     string         `gorm:"type:jsonb"` // JSON string
	ViewCount    int64          `gorm:"type:bigint;default:0"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name
func (RestaurantORM) TableName() string {
	return "restaurants"
}

// ToDomain converts ORM model to Domain entity
func (r *RestaurantORM) ToDomain() (*model.Restaurant, error) {
	// Parse Location
	location, err := model.NewLocation(r.Latitude, r.Longitude)
	if err != nil {
		return nil, err
	}

	// Parse OpeningHours
	var openingHours map[string]string
	if r.OpeningHours != "" {
		if err := json.Unmarshal([]byte(r.OpeningHours), &openingHours); err != nil {
			openingHours = make(map[string]string)
		}
	} else {
		openingHours = make(map[string]string)
	}

	// Parse Metadata
	var metadata map[string]interface{}
	if r.Metadata != "" {
		if err := json.Unmarshal([]byte(r.Metadata), &metadata); err != nil {
			metadata = make(map[string]interface{})
		}
	} else {
		metadata = make(map[string]interface{})
	}

	// Handle DeletedAt
	var deletedAt *time.Time
	if r.DeletedAt.Valid {
		deletedAt = &r.DeletedAt.Time
	}

	return model.ReconstructRestaurant(
		r.ID,
		r.Name,
		r.NameJa,
		r.Area,
		model.RestaurantSource(r.Source),
		r.ExternalID,
		r.Address,
		location,
		r.Rating,
		r.PriceRange,
		r.CuisineType,
		r.Phone,
		r.Website,
		openingHours,
		metadata,
		r.ViewCount,
		r.CreatedAt,
		r.UpdatedAt,
		deletedAt,
	), nil
}

// FromDomain converts Domain entity to ORM model
func FromDomain(r *model.Restaurant) (*RestaurantORM, error) {
	// Marshal OpeningHours
	openingHoursJSON, err := json.Marshal(r.OpeningHours())
	if err != nil {
		return nil, err
	}

	// Marshal Metadata
	metadataJSON, err := json.Marshal(r.Metadata())
	if err != nil {
		return nil, err
	}

	// Handle DeletedAt
	var deletedAt gorm.DeletedAt
	if r.DeletedAt() != nil {
		deletedAt = gorm.DeletedAt{Time: *r.DeletedAt(), Valid: true}
	}

	orm := &RestaurantORM{
		ID:           r.ID(),
		Name:         r.Name(),
		NameJa:       r.NameJa(),
		Area:         r.Area(),
		Source:       string(r.Source()),
		ExternalID:   r.ExternalID(),
		Address:      r.Address(),
		Rating:       r.Rating(),
		PriceRange:   r.PriceRange(),
		CuisineType:  r.CuisineType(),
		Phone:        r.Phone(),
		Website:      r.Website(),
		OpeningHours: string(openingHoursJSON),
		Metadata:     string(metadataJSON),
		ViewCount:    r.ViewCount(),
		CreatedAt:    r.CreatedAt(),
		UpdatedAt:    r.UpdatedAt(),
		DeletedAt:    deletedAt,
	}

	// Handle Location
	if r.Location() != nil {
		orm.Latitude = r.Location().Latitude()
		orm.Longitude = r.Location().Longitude()
	}

	return orm, nil
}

type restaurantRepository struct {
	db *gorm.DB
}

// NewRestaurantRepository creates a new postgres restaurant repository
func NewRestaurantRepository(db *gorm.DB) repository.RestaurantRepository {
	return &restaurantRepository{db: db}
}

// Create creates a new restaurant
func (r *restaurantRepository) Create(ctx context.Context, restaurant *model.Restaurant) error {
	orm, err := FromDomain(restaurant)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(orm).Error; err != nil {
		return err
	}

	return nil
}

// FindByID finds a restaurant by ID
func (r *restaurantRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Restaurant, error) {
	var orm RestaurantORM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrRestaurantNotFound
		}
		return nil, err
	}

	return orm.ToDomain()
}

// FindByExternalID finds a restaurant by source and external ID
func (r *restaurantRepository) FindByExternalID(ctx context.Context, source model.RestaurantSource, externalID string) (*model.Restaurant, error) {
	var orm RestaurantORM
	if err := r.db.WithContext(ctx).
		Where("source = ? AND external_id = ?", string(source), externalID).
		First(&orm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrRestaurantNotFound
		}
		return nil, err
	}

	return orm.ToDomain()
}

// Update updates an existing restaurant
func (r *restaurantRepository) Update(ctx context.Context, restaurant *model.Restaurant) error {
	orm, err := FromDomain(restaurant)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Model(&RestaurantORM{}).Where("id = ?", orm.ID).Updates(orm)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrRestaurantNotFound
	}

	return nil
}

// Delete soft-deletes a restaurant by ID
func (r *restaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&RestaurantORM{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrRestaurantNotFound
	}

	return nil
}

// Search searches restaurants by query string
func (r *restaurantRepository) Search(ctx context.Context, query string, limit, offset int) ([]*model.Restaurant, error) {
	var orms []RestaurantORM
	if err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR address ILIKE ? OR cuisine_type ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Find(&orms).Error; err != nil {
		return nil, err
	}

	restaurants := make([]*model.Restaurant, 0, len(orms))
	for _, orm := range orms {
		restaurant, err := orm.ToDomain()
		if err != nil {
			continue // Skip invalid records
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

// FindByLocation finds restaurants within a radius from a location
// Note: This is a simplified implementation. In production, use PostGIS for better performance
func (r *restaurantRepository) FindByLocation(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*model.Restaurant, error) {
	var orms []RestaurantORM

	// Simple bounding box query (should use PostGIS in production)
	// Approximate: 1 degree latitude ≈ 111 km
	latDelta := radiusKm / 111.0
	lngDelta := radiusKm / (111.0 * 0.9) // Approximate for mid-latitudes

	if err := r.db.WithContext(ctx).
		Where("latitude BETWEEN ? AND ?", lat-latDelta, lat+latDelta).
		Where("longitude BETWEEN ? AND ?", lng-lngDelta, lng+lngDelta).
		Limit(limit).
		Find(&orms).Error; err != nil {
		return nil, err
	}

	restaurants := make([]*model.Restaurant, 0, len(orms))
	for _, orm := range orms {
		restaurant, err := orm.ToDomain()
		if err != nil {
			continue
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

// List lists all restaurants with pagination
func (r *restaurantRepository) List(ctx context.Context, limit, offset int) ([]*model.Restaurant, error) {
	var orms []RestaurantORM
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&orms).Error; err != nil {
		return nil, err
	}

	restaurants := make([]*model.Restaurant, 0, len(orms))
	for _, orm := range orms {
		restaurant, err := orm.ToDomain()
		if err != nil {
			continue
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

// Count returns the total count of restaurants
func (r *restaurantRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RestaurantORM{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindByCuisineType finds restaurants by cuisine type
func (r *restaurantRepository) FindByCuisineType(ctx context.Context, cuisineType string, limit, offset int) ([]*model.Restaurant, error) {
	var orms []RestaurantORM
	if err := r.db.WithContext(ctx).
		Where("cuisine_type = ?", cuisineType).
		Limit(limit).
		Offset(offset).
		Find(&orms).Error; err != nil {
		return nil, err
	}

	restaurants := make([]*model.Restaurant, 0, len(orms))
	for _, orm := range orms {
		restaurant, err := orm.ToDomain()
		if err != nil {
			continue
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

// FindBySource finds restaurants by source
func (r *restaurantRepository) FindBySource(ctx context.Context, source model.RestaurantSource, limit, offset int) ([]*model.Restaurant, error) {
	var orms []RestaurantORM
	if err := r.db.WithContext(ctx).
		Where("source = ?", string(source)).
		Limit(limit).
		Offset(offset).
		Find(&orms).Error; err != nil {
		return nil, err
	}

	restaurants := make([]*model.Restaurant, 0, len(orms))
	for _, orm := range orms {
		restaurant, err := orm.ToDomain()
		if err != nil {
			continue
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}
