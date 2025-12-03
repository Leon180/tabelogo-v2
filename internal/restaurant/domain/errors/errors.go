package errors

import "errors"

var (
	// Restaurant errors
	ErrRestaurantNotFound      = errors.New("restaurant not found")
	ErrRestaurantAlreadyExists = errors.New("restaurant already exists")
	ErrInvalidLocation         = errors.New("invalid location coordinates")
	ErrInvalidRating           = errors.New("invalid rating value")

	// Favorite errors
	ErrFavoriteNotFound      = errors.New("favorite not found")
	ErrFavoriteAlreadyExists = errors.New("favorite already exists")
	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrInvalidRestaurantID   = errors.New("invalid restaurant ID")
)
