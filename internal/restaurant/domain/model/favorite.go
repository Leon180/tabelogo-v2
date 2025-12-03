package model

import (
	"time"

	"github.com/google/uuid"
)

// Favorite is the aggregate root for user's favorite restaurants
type Favorite struct {
	id            uuid.UUID
	userID        uuid.UUID
	restaurantID  uuid.UUID
	notes         string
	tags          []string
	visitCount    int
	lastVisitedAt *time.Time
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
}

// NewFavorite creates a new favorite
func NewFavorite(userID, restaurantID uuid.UUID) *Favorite {
	now := time.Now()
	return &Favorite{
		id:            uuid.New(),
		userID:        userID,
		restaurantID:  restaurantID,
		notes:         "",
		tags:          []string{},
		visitCount:    0,
		lastVisitedAt: nil,
		createdAt:     now,
		updatedAt:     now,
		deletedAt:     nil,
	}
}

// ReconstructFavorite is used by repository to reconstruct the Favorite entity from persistence
// This should NOT be used by application layer to create new favorites
func ReconstructFavorite(
	id uuid.UUID,
	userID uuid.UUID,
	restaurantID uuid.UUID,
	notes string,
	tags []string,
	visitCount int,
	lastVisitedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *Favorite {
	return &Favorite{
		id:            id,
		userID:        userID,
		restaurantID:  restaurantID,
		notes:         notes,
		tags:          tags,
		visitCount:    visitCount,
		lastVisitedAt: lastVisitedAt,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		deletedAt:     deletedAt,
	}
}

// Getters
func (f *Favorite) ID() uuid.UUID            { return f.id }
func (f *Favorite) UserID() uuid.UUID        { return f.userID }
func (f *Favorite) RestaurantID() uuid.UUID  { return f.restaurantID }
func (f *Favorite) Notes() string            { return f.notes }
func (f *Favorite) Tags() []string           { return f.tags }
func (f *Favorite) VisitCount() int          { return f.visitCount }
func (f *Favorite) LastVisitedAt() *time.Time { return f.lastVisitedAt }
func (f *Favorite) CreatedAt() time.Time     { return f.createdAt }
func (f *Favorite) UpdatedAt() time.Time     { return f.updatedAt }
func (f *Favorite) DeletedAt() *time.Time    { return f.deletedAt }

// Domain Methods

// AddVisit increments the visit count and updates last visited time
func (f *Favorite) AddVisit() {
	f.visitCount++
	now := time.Now()
	f.lastVisitedAt = &now
	f.updatedAt = now
}

// UpdateNotes updates the user's private notes
func (f *Favorite) UpdateNotes(notes string) {
	f.notes = notes
	f.updatedAt = time.Now()
}

// AddTag adds a tag if it doesn't already exist
func (f *Favorite) AddTag(tag string) {
	if tag == "" {
		return
	}

	// Check if tag already exists
	for _, existingTag := range f.tags {
		if existingTag == tag {
			return
		}
	}

	f.tags = append(f.tags, tag)
	f.updatedAt = time.Now()
}

// RemoveTag removes a tag if it exists
func (f *Favorite) RemoveTag(tag string) {
	if tag == "" {
		return
	}

	newTags := make([]string, 0, len(f.tags))
	for _, existingTag := range f.tags {
		if existingTag != tag {
			newTags = append(newTags, existingTag)
		}
	}

	if len(newTags) != len(f.tags) {
		f.tags = newTags
		f.updatedAt = time.Now()
	}
}

// SetTags replaces all tags with the provided list
func (f *Favorite) SetTags(tags []string) {
	f.tags = tags
	f.updatedAt = time.Now()
}

// HasTag checks if the favorite has a specific tag
func (f *Favorite) HasTag(tag string) bool {
	for _, existingTag := range f.tags {
		if existingTag == tag {
			return true
		}
	}
	return false
}

// SoftDelete marks the favorite as deleted
func (f *Favorite) SoftDelete() {
	now := time.Now()
	f.deletedAt = &now
	f.updatedAt = now
}

// IsDeleted checks if the favorite is soft-deleted
func (f *Favorite) IsDeleted() bool {
	return f.deletedAt != nil
}
