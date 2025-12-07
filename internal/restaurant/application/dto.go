package application

import "github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"

// CreateRestaurantRequest represents a request to create a restaurant
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

// UpdateRestaurantRequest represents a request to update a restaurant
type UpdateRestaurantRequest struct {
	Name        string  `json:"name"`
	NameJa      string  `json:"name_ja"`
	Address     string  `json:"address"`
	Rating      float64 `json:"rating"`
	PriceRange  string  `json:"price_range"`
	CuisineType string  `json:"cuisine_type"`
	Phone       string  `json:"phone"`
	Website     string  `json:"website"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
