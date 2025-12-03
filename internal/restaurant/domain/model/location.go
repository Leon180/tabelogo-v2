package model

import "fmt"

// Location is a value object representing geographic coordinates
type Location struct {
	latitude  float64
	longitude float64
}

// NewLocation creates a new Location with validation
func NewLocation(latitude, longitude float64) (*Location, error) {
	if latitude < -90 || latitude > 90 {
		return nil, fmt.Errorf("invalid latitude: must be between -90 and 90, got %f", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return nil, fmt.Errorf("invalid longitude: must be between -180 and 180, got %f", longitude)
	}

	return &Location{
		latitude:  latitude,
		longitude: longitude,
	}, nil
}

// Getters
func (l *Location) Latitude() float64  { return l.latitude }
func (l *Location) Longitude() float64 { return l.longitude }

// String returns a string representation of the location
func (l *Location) String() string {
	return fmt.Sprintf("(%f, %f)", l.latitude, l.longitude)
}

// Equals checks if two locations are equal
func (l *Location) Equals(other *Location) bool {
	if other == nil {
		return false
	}
	return l.latitude == other.latitude && l.longitude == other.longitude
}
