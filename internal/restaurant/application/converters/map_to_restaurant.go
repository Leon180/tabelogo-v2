package converters

import (
	"log"

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

	// Parse opening hours from proto to domain model
	openingHours := make(map[string]string)
	if place.OpeningHours != nil && len(place.OpeningHours.WeekdayText) > 0 {
		for i, text := range place.OpeningHours.WeekdayText {
			dayKey := getDayKey(i, text)
			openingHours[dayKey] = text
		}
	}

	// Build metadata with photos and other info
	metadata := make(map[string]interface{})
	if len(place.Photos) > 0 {
		photos := make([]map[string]interface{}, 0, len(place.Photos))
		for _, photo := range place.Photos {
			photos = append(photos, map[string]interface{}{
				"name":   photo.Name,
				"width":  photo.WidthPx,
				"height": photo.HeightPx,
			})
		}
		metadata["photos"] = photos
		metadata["photo_count"] = len(place.Photos)
	}

	// Add other useful metadata
	if place.OpeningHours != nil {
		metadata["open_now"] = place.OpeningHours.OpenNow
	}
	if place.EditorialSummary != "" {
		metadata["editorial_summary"] = place.EditorialSummary
	}
	if place.UserRatingsTotal > 0 {
		metadata["user_ratings_total"] = place.UserRatingsTotal
	}
	if len(place.Types) > 0 {
		metadata["types"] = place.Types
	}

	// Set price range from Google's price level
	priceRange := parsePriceLevel(place.PriceLevel)

	// Extract area from addressComponents (e.g., "Tokyo")
	area := extractAreaFromAddressComponents(place.AddressComponents)

	// DEBUG: Log conversion
	log.Printf("[MapPlaceToRestaurant] place_id=%s, name=%s, addressComponents_count=%d, extracted_area=%s",
		place.Id, place.Name, len(place.AddressComponents), area)

	// Create restaurant with all details in one call
	restaurant := model.NewRestaurantWithDetails(
		place.Name,
		area,
		model.SourceGoogle,
		place.Id,
		place.FormattedAddress,
		location,
		place.Rating,
		priceRange,
		"", // cuisineType - can be derived from types later
		place.PhoneNumber,
		place.Website,
		openingHours,
		metadata,
	)

	// Set area if extracted
	if area != "" {
		restaurant.UpdateArea(area)
	}

	return restaurant
}

// extractAreaFromAddressComponents extracts the administrative area level 1 from address components
// This is used to get the area (e.g., "Tokyo", "Osaka") for Tabelog search URLs
func extractAreaFromAddressComponents(components []*mapv1.AddressComponent) string {
	if len(components) == 0 {
		return ""
	}

	var fallbackArea string

	// DEBUG: Log input
	log.Printf("[Area Extraction] Received %d address components", len(components))

	for i, component := range components {
		log.Printf("[Area Extraction] Component %d: longText=%s, shortText=%s, languageCode=%s, types=%v",
			i, component.LongText, component.ShortText, component.LanguageCode, component.Types)

		for _, typ := range component.Types {
			if typ == "administrative_area_level_1" {
				log.Printf("[Area Extraction] Found administrative_area_level_1: shortText=%s, languageCode=%s",
					component.ShortText, component.LanguageCode)

				// Prioritize English language code
				if component.LanguageCode == "en" || component.LanguageCode == "" {
					if component.ShortText != "" {
						log.Printf("[Area Extraction] ✅ Using English area: %s", component.ShortText)
						return component.ShortText // "Tokyo"
					}
				}

				// Fallback for non-English
				if fallbackArea == "" {
					if component.ShortText != "" {
						fallbackArea = component.ShortText
						log.Printf("[Area Extraction] Set fallback area: %s", fallbackArea)
					} else if component.LongText != "" {
						fallbackArea = component.LongText
						log.Printf("[Area Extraction] Set fallback area (longText): %s", fallbackArea)
					}
				}
			}
		}
	}

	if fallbackArea != "" {
		log.Printf("[Area Extraction] ⚠️ Using fallback area: %s", fallbackArea)
	} else {
		log.Printf("[Area Extraction] ❌ No area found!")
	}

	return fallbackArea
}

// getDayKey extracts the day name from the opening hours text
// Example: "Monday: 8:30 AM – 5:30 PM" -> "Monday"
func getDayKey(index int, text string) string {
	// Try to extract day name from text
	if len(text) > 0 {
		// Find the colon position
		colonIdx := 0
		for i, ch := range text {
			if ch == ':' {
				colonIdx = i
				break
			}
		}
		if colonIdx > 0 {
			return text[:colonIdx]
		}
	}

	// Fallback to day index
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	if index >= 0 && index < len(days) {
		return days[index]
	}
	return "Unknown"
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
