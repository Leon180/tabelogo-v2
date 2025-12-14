package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRestaurant(t *testing.T) {
	location, err := NewLocation(35.6762, 139.6503)
	require.NoError(t, err)

	restaurant := NewRestaurant(
		"Sushi Dai",
		"",
		SourceGoogle,
		"ChIJTest123",
		"Tokyo, Chuo-ku",
		location,
	)

	assert.NotNil(t, restaurant)
	assert.NotEqual(t, uuid.Nil, restaurant.ID())
	assert.Equal(t, "Sushi Dai", restaurant.Name())
	assert.Equal(t, SourceGoogle, restaurant.Source())
	assert.Equal(t, "ChIJTest123", restaurant.ExternalID())
	assert.Equal(t, "Tokyo, Chuo-ku", restaurant.Address())
	assert.Equal(t, location, restaurant.Location())
	assert.Equal(t, 0.0, restaurant.Rating())
	assert.Equal(t, "", restaurant.PriceRange())
	assert.Equal(t, "", restaurant.CuisineType())
	assert.Equal(t, int64(0), restaurant.ViewCount())
	assert.NotZero(t, restaurant.CreatedAt())
	assert.NotZero(t, restaurant.UpdatedAt())
	assert.Nil(t, restaurant.DeletedAt())
}

func TestRestaurant_UpdateRating(t *testing.T) {
	restaurant := createTestRestaurant(t)

	tests := []struct {
		name          string
		rating        float64
		expectedValue float64
	}{
		{
			name:          "Valid rating",
			rating:        4.5,
			expectedValue: 4.5,
		},
		{
			name:          "Rating below minimum",
			rating:        -1.0,
			expectedValue: 0.0,
		},
		{
			name:          "Rating above maximum",
			rating:        6.0,
			expectedValue: 5.0,
		},
		{
			name:          "Zero rating",
			rating:        0.0,
			expectedValue: 0.0,
		},
		{
			name:          "Maximum rating",
			rating:        5.0,
			expectedValue: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := restaurant.UpdatedAt()
			time.Sleep(1 * time.Millisecond)

			restaurant.UpdateRating(tt.rating)

			assert.Equal(t, tt.expectedValue, restaurant.Rating())
			assert.True(t, restaurant.UpdatedAt().After(oldUpdatedAt))
		})
	}
}

func TestRestaurant_IncrementViewCount(t *testing.T) {
	restaurant := createTestRestaurant(t)

	assert.Equal(t, int64(0), restaurant.ViewCount())

	restaurant.IncrementViewCount()
	assert.Equal(t, int64(1), restaurant.ViewCount())

	restaurant.IncrementViewCount()
	assert.Equal(t, int64(2), restaurant.ViewCount())
}

func TestRestaurant_UpdateDetails(t *testing.T) {
	restaurant := createTestRestaurant(t)

	restaurant.UpdateDetails(
		"Updated Name",
		"Updated Address",
		"$$",
		"Japanese",
		"03-1234-5678",
		"https://example.com",
	)

	assert.Equal(t, "Updated Name", restaurant.Name())
	assert.Equal(t, "Updated Address", restaurant.Address())
	assert.Equal(t, "$$", restaurant.PriceRange())
	assert.Equal(t, "Japanese", restaurant.CuisineType())
	assert.Equal(t, "03-1234-5678", restaurant.Phone())
	assert.Equal(t, "https://example.com", restaurant.Website())
}

func TestRestaurant_UpdateDetails_EmptyValues(t *testing.T) {
	restaurant := createTestRestaurant(t)
	restaurant.UpdateDetails("Original", "Address", "$$", "Sushi", "123", "web")

	// Update with empty values should not change
	restaurant.UpdateDetails("", "", "", "", "", "")

	assert.Equal(t, "Original", restaurant.Name())
	assert.Equal(t, "Address", restaurant.Address())
	assert.Equal(t, "$$", restaurant.PriceRange())
	assert.Equal(t, "Sushi", restaurant.CuisineType())
	assert.Equal(t, "123", restaurant.Phone())
	assert.Equal(t, "web", restaurant.Website())
}

func TestRestaurant_UpdateLocation(t *testing.T) {
	restaurant := createTestRestaurant(t)

	newLocation, err := NewLocation(34.0, 135.0)
	require.NoError(t, err)

	oldUpdatedAt := restaurant.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	restaurant.UpdateLocation(newLocation)

	assert.Equal(t, newLocation, restaurant.Location())
	assert.True(t, restaurant.UpdatedAt().After(oldUpdatedAt))
}

func TestRestaurant_UpdateLocation_Nil(t *testing.T) {
	restaurant := createTestRestaurant(t)
	originalLocation := restaurant.Location()
	oldUpdatedAt := restaurant.UpdatedAt()

	time.Sleep(1 * time.Millisecond)
	restaurant.UpdateLocation(nil)

	assert.Equal(t, originalLocation, restaurant.Location())
	assert.Equal(t, oldUpdatedAt, restaurant.UpdatedAt())
}

func TestRestaurant_SetOpeningHours(t *testing.T) {
	restaurant := createTestRestaurant(t)

	restaurant.SetOpeningHours("Monday", "09:00-17:00")
	restaurant.SetOpeningHours("Tuesday", "09:00-18:00")

	hours := restaurant.OpeningHours()
	assert.Len(t, hours, 2)
	assert.Equal(t, "09:00-17:00", hours["Monday"])
	assert.Equal(t, "09:00-18:00", hours["Tuesday"])
}

func TestRestaurant_SetMetadata(t *testing.T) {
	restaurant := createTestRestaurant(t)

	restaurant.SetMetadata("price_level", "MODERATE")
	restaurant.SetMetadata("photo_count", 10)
	restaurant.SetMetadata("verified", true)

	metadata := restaurant.Metadata()
	assert.Len(t, metadata, 3)
	assert.Equal(t, "MODERATE", metadata["price_level"])
	assert.Equal(t, 10, metadata["photo_count"])
	assert.Equal(t, true, metadata["verified"])
}

func TestRestaurant_SoftDelete(t *testing.T) {
	restaurant := createTestRestaurant(t)

	assert.False(t, restaurant.IsDeleted())
	assert.Nil(t, restaurant.DeletedAt())

	restaurant.SoftDelete()

	assert.True(t, restaurant.IsDeleted())
	assert.NotNil(t, restaurant.DeletedAt())
	assert.True(t, restaurant.DeletedAt().Before(time.Now().Add(1*time.Second)))
}

func TestReconstructRestaurant(t *testing.T) {
	id := uuid.New()
	location, _ := NewLocation(35.6762, 139.6503)
	openingHours := map[string]string{"Monday": "09:00-17:00"}
	metadata := map[string]interface{}{"key": "value"}
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()
	deletedAt := time.Now()

	restaurant := ReconstructRestaurant(
		id,
		"Test Restaurant",
		"Test Restaurant",
		"Test Restaurant",
		SourceTabelog,
		"external123",
		"Test Address",
		location,
		4.5,
		"$$",
		"Sushi",
		"03-1234-5678",
		"https://example.com",
		openingHours,
		metadata,
		100,
		createdAt,
		updatedAt,
		&deletedAt,
	)

	assert.Equal(t, id, restaurant.ID())
	assert.Equal(t, "Test Restaurant", restaurant.Name())
	assert.Equal(t, SourceTabelog, restaurant.Source())
	assert.Equal(t, "external123", restaurant.ExternalID())
	assert.Equal(t, "Test Address", restaurant.Address())
	assert.Equal(t, location, restaurant.Location())
	assert.Equal(t, 4.5, restaurant.Rating())
	assert.Equal(t, "$$", restaurant.PriceRange())
	assert.Equal(t, "Sushi", restaurant.CuisineType())
	assert.Equal(t, "03-1234-5678", restaurant.Phone())
	assert.Equal(t, "https://example.com", restaurant.Website())
	assert.Equal(t, openingHours, restaurant.OpeningHours())
	assert.Equal(t, metadata, restaurant.Metadata())
	assert.Equal(t, int64(100), restaurant.ViewCount())
	assert.Equal(t, createdAt, restaurant.CreatedAt())
	assert.Equal(t, updatedAt, restaurant.UpdatedAt())
	assert.Equal(t, &deletedAt, restaurant.DeletedAt())
	assert.True(t, restaurant.IsDeleted())
}

// Helper function
func createTestRestaurant(t *testing.T) *Restaurant {
	location, err := NewLocation(35.6762, 139.6503)
	require.NoError(t, err)

	return NewRestaurant(
		"Test Restaurant",
		"Test Area",
		SourceGoogle,
		"test-external-id",
		"Test Address",
		location,
	)
}
