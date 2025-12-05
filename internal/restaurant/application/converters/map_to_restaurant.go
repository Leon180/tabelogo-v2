package converters

import (
	mapv1 "github.com/Leon180/tabelogo-v2/api/gen/map/v1"
	"github.com/Leon180/tabelogo-v2/internal/restaurant/domain/model"
)

// MapPlaceToRestaurant converts a Map Service Place proto to Restaurant domain model
func MapPlaceToRestaurant(place *mapv1.Place) *model.Restaurant {
	if place == nil {
		return nil
	}

	// Create location if coordinates are available
	var location *model.Location
	if place.Location != nil {
		var err error
		location, err = model.NewLocation(
			place.Location.Latitude,
			place.Location.Longitude,
		)
		if err != nil {
			// Log error but continue - location is optional
			location = nil
		}
	}

	// Create new restaurant with basic info
	restaurant := model.NewRestaurant(
		place.Name,
		model.SourceGoogle, // Data from Google Maps
		place.Id,           // Google Place ID as external ID
		place.FormattedAddress,
		location,
	)

	// Update additional details
	restaurant.UpdateRating(place.Rating)

	// Set price range from Google's price level
	priceRange := parsePriceLevel(place.PriceLevel)
	restaurant.UpdateDetails(
		place.Name,             // name
		place.FormattedAddress, // address
		priceRange,             // priceRange
		"",                     // cuisineType - can be derived from types later
		place.PhoneNumber,      // phone
		place.Website,          // website
	)

	return restaurant
}

// parsePriceLevel converts Google's price level string to our price range format
func parsePriceLevel(priceLevel string) string {
	switch priceLevel {
	case "PRICE_LEVEL_FREE":
		return "$"
	case "PRICE_LEVEL_INEXPENSIVE":
		return "$"
	case "PRICE_LEVEL_MODERATE":
		return "$$"
	case "PRICE_LEVEL_EXPENSIVE":
		return "$$$"
	case "PRICE_LEVEL_VERY_EXPENSIVE":
		return "$$$$"
	default:
		return ""
	}
}

// MapPlacesToRestaurants converts multiple Map Service Places to Restaurant domain models
func MapPlacesToRestaurants(places []*mapv1.Place) []*model.Restaurant {
	if len(places) == 0 {
		return nil
	}

	restaurants := make([]*model.Restaurant, 0, len(places))
	for _, place := range places {
		if restaurant := MapPlaceToRestaurant(place); restaurant != nil {
			restaurants = append(restaurants, restaurant)
		}
	}

	return restaurants
}
