package model

import (
	"time"

	"github.com/google/uuid"
)

// RestaurantSource defines the source of restaurant data
type RestaurantSource string

const (
	SourceTabelog   RestaurantSource = "tabelog"
	SourceGoogle    RestaurantSource = "google"
	SourceOpenTable RestaurantSource = "opentable"
)

// Restaurant is the aggregate root for the restaurant domain
type Restaurant struct {
	id           uuid.UUID
	name         string
	nameJa       string // Japanese name for better Tabelog search
	source       RestaurantSource
	externalID   string
	address      string
	location     *Location
	rating       float64
	priceRange   string
	cuisineType  string
	phone        string
	website      string
	openingHours map[string]string
	metadata     map[string]interface{}
	viewCount    int64
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

// NewRestaurant creates a new restaurant
func NewRestaurant(
	name string,
	source RestaurantSource,
	externalID string,
	address string,
	location *Location,
) *Restaurant {
	now := time.Now()
	return &Restaurant{
		id:           uuid.New(),
		name:         name,
		nameJa:       "", // Will be set later via frontend or update API
		source:       source,
		externalID:   externalID,
		address:      address,
		location:     location,
		rating:       0.0,
		priceRange:   "",
		cuisineType:  "",
		phone:        "",
		website:      "",
		openingHours: make(map[string]string),
		metadata:     make(map[string]interface{}),
		viewCount:    0,
		createdAt:    now,
		updatedAt:    now,
		deletedAt:    nil,
	}
}

// NewRestaurantWithDetails creates a new restaurant with complete details
func NewRestaurantWithDetails(
	name string,
	source RestaurantSource,
	externalID string,
	address string,
	location *Location,
	rating float64,
	priceRange string,
	cuisineType string,
	phone string,
	website string,
	openingHours map[string]string,
	metadata map[string]interface{},
) *Restaurant {
	now := time.Now()

	// Initialize maps if nil
	if openingHours == nil {
		openingHours = make(map[string]string)
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &Restaurant{
		id:           uuid.New(),
		name:         name,
		nameJa:       "", // Will be set later
		source:       source,
		externalID:   externalID,
		address:      address,
		location:     location,
		rating:       rating,
		priceRange:   priceRange,
		cuisineType:  cuisineType,
		phone:        phone,
		website:      website,
		openingHours: openingHours,
		metadata:     metadata,
		viewCount:    0,
		createdAt:    now,
		updatedAt:    now,
		deletedAt:    nil,
	}
}

// ReconstructRestaurant is used by repository to reconstruct the Restaurant entity from persistence
// This should NOT be used by application layer to create new restaurants
func ReconstructRestaurant(
	id uuid.UUID,
	name string,
	nameJa string,
	source RestaurantSource,
	externalID string,
	address string,
	location *Location,
	rating float64,
	priceRange string,
	cuisineType string,
	phone string,
	website string,
	openingHours map[string]string,
	metadata map[string]interface{},
	viewCount int64,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *Restaurant {
	return &Restaurant{
		id:           id,
		name:         name,
		nameJa:       nameJa,
		source:       source,
		externalID:   externalID,
		address:      address,
		location:     location,
		rating:       rating,
		priceRange:   priceRange,
		cuisineType:  cuisineType,
		phone:        phone,
		website:      website,
		openingHours: openingHours,
		metadata:     metadata,
		viewCount:    viewCount,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		deletedAt:    deletedAt,
	}
}

// Getters
func (r *Restaurant) ID() uuid.UUID                    { return r.id }
func (r *Restaurant) Name() string                     { return r.name }
func (r *Restaurant) NameJa() string                   { return r.nameJa }
func (r *Restaurant) Source() RestaurantSource         { return r.source }
func (r *Restaurant) ExternalID() string               { return r.externalID }
func (r *Restaurant) Address() string                  { return r.address }
func (r *Restaurant) Location() *Location              { return r.location }
func (r *Restaurant) Rating() float64                  { return r.rating }
func (r *Restaurant) PriceRange() string               { return r.priceRange }
func (r *Restaurant) CuisineType() string              { return r.cuisineType }
func (r *Restaurant) Phone() string                    { return r.phone }
func (r *Restaurant) Website() string                  { return r.website }
func (r *Restaurant) OpeningHours() map[string]string  { return r.openingHours }
func (r *Restaurant) Metadata() map[string]interface{} { return r.metadata }
func (r *Restaurant) ViewCount() int64                 { return r.viewCount }
func (r *Restaurant) CreatedAt() time.Time             { return r.createdAt }
func (r *Restaurant) UpdatedAt() time.Time             { return r.updatedAt }
func (r *Restaurant) DeletedAt() *time.Time            { return r.deletedAt }

// Domain Methods

// UpdateRating updates the restaurant's rating
func (r *Restaurant) UpdateRating(rating float64) {
	if rating < 0.0 {
		rating = 0.0
	}
	if rating > 5.0 {
		rating = 5.0
	}
	r.rating = rating
	r.updatedAt = time.Now()
}

// IncrementViewCount increments the view count by 1
func (r *Restaurant) IncrementViewCount() {
	r.viewCount++
	r.updatedAt = time.Now()
}

// UpdateDetails updates the restaurant's details
func (r *Restaurant) UpdateDetails(
	name, address, priceRange, cuisineType, phone, website string,
) {
	if name != "" {
		r.name = name
	}
	if address != "" {
		r.address = address
	}
	if priceRange != "" {
		r.priceRange = priceRange
	}
	if cuisineType != "" {
		r.cuisineType = cuisineType
	}
	if phone != "" {
		r.phone = phone
	}
	if website != "" {
		r.website = website
	}
	r.updatedAt = time.Now()
}

// UpdateLocation updates the restaurant's location
func (r *Restaurant) UpdateLocation(location *Location) {
	if location != nil {
		r.location = location
		r.updatedAt = time.Now()
	}
}

// SetOpeningHours sets the opening hours for a specific day
func (r *Restaurant) SetOpeningHours(day, hours string) {
	if r.openingHours == nil {
		r.openingHours = make(map[string]string)
	}
	r.openingHours[day] = hours
	r.updatedAt = time.Now()
}

// UpdateOpeningHours replaces all opening hours with a new map
func (r *Restaurant) UpdateOpeningHours(openingHours map[string]string) {
	if openingHours != nil {
		r.openingHours = openingHours
		r.updatedAt = time.Now()
	}
}

// SetMetadata sets a metadata key-value pair
func (r *Restaurant) SetMetadata(key string, value interface{}) {
	if r.metadata == nil {
		r.metadata = make(map[string]interface{})
	}
	r.metadata[key] = value
	r.updatedAt = time.Now()
}

// UpdateMetadata replaces all metadata with a new map
func (r *Restaurant) UpdateMetadata(metadata map[string]interface{}) {
	if metadata != nil {
		r.metadata = metadata
		r.updatedAt = time.Now()
	}
}

// SoftDelete marks the restaurant as deleted
func (r *Restaurant) SoftDelete() {
	now := time.Now()
	r.deletedAt = &now
	r.updatedAt = now
}

// IsDeleted checks if the restaurant is soft-deleted
func (r *Restaurant) IsDeleted() bool {
	return r.deletedAt != nil
}
