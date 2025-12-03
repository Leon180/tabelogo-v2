package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewFavorite(t *testing.T) {
	userID := uuid.New()
	restaurantID := uuid.New()

	favorite := NewFavorite(userID, restaurantID)

	assert.NotNil(t, favorite)
	assert.NotEqual(t, uuid.Nil, favorite.ID())
	assert.Equal(t, userID, favorite.UserID())
	assert.Equal(t, restaurantID, favorite.RestaurantID())
	assert.Equal(t, "", favorite.Notes())
	assert.Empty(t, favorite.Tags())
	assert.Equal(t, 0, favorite.VisitCount())
	assert.Nil(t, favorite.LastVisitedAt())
	assert.NotZero(t, favorite.CreatedAt())
	assert.NotZero(t, favorite.UpdatedAt())
	assert.Nil(t, favorite.DeletedAt())
	assert.False(t, favorite.IsDeleted())
}

func TestFavorite_AddVisit(t *testing.T) {
	favorite := createTestFavorite(t)

	assert.Equal(t, 0, favorite.VisitCount())
	assert.Nil(t, favorite.LastVisitedAt())

	oldUpdatedAt := favorite.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	favorite.AddVisit()

	assert.Equal(t, 1, favorite.VisitCount())
	assert.NotNil(t, favorite.LastVisitedAt())
	assert.True(t, favorite.LastVisitedAt().Before(time.Now().Add(1*time.Second)))
	assert.True(t, favorite.UpdatedAt().After(oldUpdatedAt))

	time.Sleep(1 * time.Millisecond)
	favorite.AddVisit()
	assert.Equal(t, 2, favorite.VisitCount())
}

func TestFavorite_UpdateNotes(t *testing.T) {
	favorite := createTestFavorite(t)

	oldUpdatedAt := favorite.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	favorite.UpdateNotes("Great sushi place!")

	assert.Equal(t, "Great sushi place!", favorite.Notes())
	assert.True(t, favorite.UpdatedAt().After(oldUpdatedAt))

	favorite.UpdateNotes("Updated notes")
	assert.Equal(t, "Updated notes", favorite.Notes())
}

func TestFavorite_AddTag(t *testing.T) {
	favorite := createTestFavorite(t)

	favorite.AddTag("sushi")
	assert.Len(t, favorite.Tags(), 1)
	assert.Contains(t, favorite.Tags(), "sushi")

	favorite.AddTag("affordable")
	assert.Len(t, favorite.Tags(), 2)
	assert.Contains(t, favorite.Tags(), "affordable")

	// Adding duplicate tag should not add it again
	favorite.AddTag("sushi")
	assert.Len(t, favorite.Tags(), 2)
}

func TestFavorite_AddTag_EmptyString(t *testing.T) {
	favorite := createTestFavorite(t)

	favorite.AddTag("")
	assert.Empty(t, favorite.Tags())
}

func TestFavorite_RemoveTag(t *testing.T) {
	favorite := createTestFavorite(t)
	favorite.AddTag("sushi")
	favorite.AddTag("affordable")
	favorite.AddTag("cozy")

	assert.Len(t, favorite.Tags(), 3)

	oldUpdatedAt := favorite.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	favorite.RemoveTag("affordable")

	assert.Len(t, favorite.Tags(), 2)
	assert.Contains(t, favorite.Tags(), "sushi")
	assert.Contains(t, favorite.Tags(), "cozy")
	assert.NotContains(t, favorite.Tags(), "affordable")
	assert.True(t, favorite.UpdatedAt().After(oldUpdatedAt))
}

func TestFavorite_RemoveTag_NonExistent(t *testing.T) {
	favorite := createTestFavorite(t)
	favorite.AddTag("sushi")

	oldUpdatedAt := favorite.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	favorite.RemoveTag("nonexistent")

	assert.Len(t, favorite.Tags(), 1)
	assert.Equal(t, oldUpdatedAt, favorite.UpdatedAt())
}

func TestFavorite_RemoveTag_EmptyString(t *testing.T) {
	favorite := createTestFavorite(t)
	favorite.AddTag("sushi")

	favorite.RemoveTag("")

	assert.Len(t, favorite.Tags(), 1)
}

func TestFavorite_SetTags(t *testing.T) {
	favorite := createTestFavorite(t)

	newTags := []string{"sushi", "expensive", "michelin"}
	favorite.SetTags(newTags)

	assert.Equal(t, newTags, favorite.Tags())

	// Set new tags should replace old ones
	anotherTags := []string{"ramen", "casual"}
	favorite.SetTags(anotherTags)

	assert.Equal(t, anotherTags, favorite.Tags())
	assert.Len(t, favorite.Tags(), 2)
}

func TestFavorite_HasTag(t *testing.T) {
	favorite := createTestFavorite(t)
	favorite.AddTag("sushi")
	favorite.AddTag("affordable")

	assert.True(t, favorite.HasTag("sushi"))
	assert.True(t, favorite.HasTag("affordable"))
	assert.False(t, favorite.HasTag("expensive"))
	assert.False(t, favorite.HasTag(""))
}

func TestFavorite_SoftDelete(t *testing.T) {
	favorite := createTestFavorite(t)

	assert.False(t, favorite.IsDeleted())
	assert.Nil(t, favorite.DeletedAt())

	oldUpdatedAt := favorite.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	favorite.SoftDelete()

	assert.True(t, favorite.IsDeleted())
	assert.NotNil(t, favorite.DeletedAt())
	assert.True(t, favorite.DeletedAt().Before(time.Now().Add(1*time.Second)))
	assert.True(t, favorite.UpdatedAt().After(oldUpdatedAt))
}

func TestReconstructFavorite(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	restaurantID := uuid.New()
	tags := []string{"sushi", "affordable"}
	lastVisited := time.Now().Add(-1 * time.Hour)
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()
	deletedAt := time.Now()

	favorite := ReconstructFavorite(
		id,
		userID,
		restaurantID,
		"Great place!",
		tags,
		5,
		&lastVisited,
		createdAt,
		updatedAt,
		&deletedAt,
	)

	assert.Equal(t, id, favorite.ID())
	assert.Equal(t, userID, favorite.UserID())
	assert.Equal(t, restaurantID, favorite.RestaurantID())
	assert.Equal(t, "Great place!", favorite.Notes())
	assert.Equal(t, tags, favorite.Tags())
	assert.Equal(t, 5, favorite.VisitCount())
	assert.Equal(t, &lastVisited, favorite.LastVisitedAt())
	assert.Equal(t, createdAt, favorite.CreatedAt())
	assert.Equal(t, updatedAt, favorite.UpdatedAt())
	assert.Equal(t, &deletedAt, favorite.DeletedAt())
	assert.True(t, favorite.IsDeleted())
}

// Helper function
func createTestFavorite(t *testing.T) *Favorite {
	userID := uuid.New()
	restaurantID := uuid.New()
	return NewFavorite(userID, restaurantID)
}
