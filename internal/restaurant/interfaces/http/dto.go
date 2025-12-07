package http

import (
	"time"

	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
)

// Request DTOs

type CreateRestaurantRequest struct {
	Name        string                 `json:"name" binding:"required"`
	NameJa      string                 `json:"name_ja"`
	Source      model.RestaurantSource `json:"source" binding:"required"`
	ExternalID  string                 `json:"external_id" binding:"required"`
	Address     string                 `json:"address"`
	Latitude    float64                `json:"latitude" binding:"required"`
	Longitude   float64                `json:"longitude" binding:"required"`
	Rating      float64                `json:"rating"`
	PriceRange  string                 `json:"price_range"`
	CuisineType string                 `json:"cuisine_type"`
	Phone       string                 `json:"phone"`
	Website     string                 `json:"website"`
}

// UpdateRestaurantRequest represents the HTTP request to update a restaurant
type UpdateRestaurantRequest struct {
	NameJa string `json:"name_ja"`
}

type AddFavoriteRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	RestaurantID string `json:"restaurant_id" binding:"required"`
}

// Response DTOs

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type RestaurantDTO struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	NameJa       string                 `json:"name_ja,omitempty"`
	Source       string                 `json:"source"`
	ExternalID   string                 `json:"external_id"`
	Address      string                 `json:"address"`
	Latitude     float64                `json:"latitude"`
	Longitude    float64                `json:"longitude"`
	Rating       float64                `json:"rating"`
	PriceRange   string                 `json:"price_range"`
	CuisineType  string                 `json:"cuisine_type"`
	Phone        string                 `json:"phone"`
	Website      string                 `json:"website"`
	OpeningHours map[string]string      `json:"opening_hours,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ViewCount    int64                  `json:"view_count"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type RestaurantResponse struct {
	Restaurant RestaurantDTO `json:"restaurant"`
}

type RestaurantListResponse struct {
	Restaurants []RestaurantDTO `json:"restaurants"`
	Total       int             `json:"total"`
}

type FavoriteDTO struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	RestaurantID  string     `json:"restaurant_id"`
	Notes         string     `json:"notes"`
	Tags          []string   `json:"tags"`
	VisitCount    int        `json:"visit_count"`
	LastVisitedAt *time.Time `json:"last_visited_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type FavoriteResponse struct {
	Favorite FavoriteDTO `json:"favorite"`
}

type FavoriteListResponse struct {
	Favorites []FavoriteDTO `json:"favorites"`
	Total     int           `json:"total"`
}

// Mapper functions

func toRestaurantDTO(r *model.Restaurant) RestaurantDTO {
	var lat, lng float64
	if r.Location() != nil {
		lat = r.Location().Latitude()
		lng = r.Location().Longitude()
	}

	return RestaurantDTO{
		ID:           r.ID().String(),
		Name:         r.Name(),
		NameJa:       r.NameJa(),
		Source:       string(r.Source()),
		ExternalID:   r.ExternalID(),
		Address:      r.Address(),
		Latitude:     lat,
		Longitude:    lng,
		Rating:       r.Rating(),
		PriceRange:   r.PriceRange(),
		CuisineType:  r.CuisineType(),
		Phone:        r.Phone(),
		Website:      r.Website(),
		OpeningHours: r.OpeningHours(),
		Metadata:     r.Metadata(),
		ViewCount:    r.ViewCount(),
		CreatedAt:    r.CreatedAt(),
		UpdatedAt:    r.UpdatedAt(),
	}
}

func toRestaurantDTOList(restaurants []*model.Restaurant) []RestaurantDTO {
	dtos := make([]RestaurantDTO, len(restaurants))
	for i, r := range restaurants {
		dtos[i] = toRestaurantDTO(r)
	}
	return dtos
}

func toFavoriteDTO(f *model.Favorite) FavoriteDTO {
	return FavoriteDTO{
		ID:            f.ID().String(),
		UserID:        f.UserID().String(),
		RestaurantID:  f.RestaurantID().String(),
		Notes:         f.Notes(),
		Tags:          f.Tags(),
		VisitCount:    f.VisitCount(),
		LastVisitedAt: f.LastVisitedAt(),
		CreatedAt:     f.CreatedAt(),
	}
}

func toFavoriteDTOList(favorites []*model.Favorite) []FavoriteDTO {
	dtos := make([]FavoriteDTO, len(favorites))
	for i, f := range favorites {
		dtos[i] = toFavoriteDTO(f)
	}
	return dtos
}
