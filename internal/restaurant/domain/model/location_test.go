package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocation_Success(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{
			name:      "Valid Tokyo location",
			latitude:  35.6762,
			longitude: 139.6503,
		},
		{
			name:      "Valid boundary latitude max",
			latitude:  90.0,
			longitude: 0.0,
		},
		{
			name:      "Valid boundary latitude min",
			latitude:  -90.0,
			longitude: 0.0,
		},
		{
			name:      "Valid boundary longitude max",
			latitude:  0.0,
			longitude: 180.0,
		},
		{
			name:      "Valid boundary longitude min",
			latitude:  0.0,
			longitude: -180.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := NewLocation(tt.latitude, tt.longitude)

			assert.NoError(t, err)
			assert.NotNil(t, location)
			assert.Equal(t, tt.latitude, location.Latitude())
			assert.Equal(t, tt.longitude, location.Longitude())
		})
	}
}

func TestNewLocation_InvalidLatitude(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{
			name:      "Latitude too high",
			latitude:  90.1,
			longitude: 0.0,
		},
		{
			name:      "Latitude too low",
			latitude:  -90.1,
			longitude: 0.0,
		},
		{
			name:      "Latitude way out of range",
			latitude:  200.0,
			longitude: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := NewLocation(tt.latitude, tt.longitude)

			assert.Error(t, err)
			assert.Nil(t, location)
			assert.Contains(t, err.Error(), "invalid latitude")
		})
	}
}

func TestNewLocation_InvalidLongitude(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{
			name:      "Longitude too high",
			latitude:  0.0,
			longitude: 180.1,
		},
		{
			name:      "Longitude too low",
			latitude:  0.0,
			longitude: -180.1,
		},
		{
			name:      "Longitude way out of range",
			latitude:  0.0,
			longitude: 300.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := NewLocation(tt.latitude, tt.longitude)

			assert.Error(t, err)
			assert.Nil(t, location)
			assert.Contains(t, err.Error(), "invalid longitude")
		})
	}
}

func TestLocation_String(t *testing.T) {
	location, err := NewLocation(35.6762, 139.6503)

	assert.NoError(t, err)
	assert.Equal(t, "(35.676200, 139.650300)", location.String())
}

func TestLocation_Equals(t *testing.T) {
	location1, _ := NewLocation(35.6762, 139.6503)
	location2, _ := NewLocation(35.6762, 139.6503)
	location3, _ := NewLocation(35.6763, 139.6503)

	assert.True(t, location1.Equals(location2))
	assert.False(t, location1.Equals(location3))
	assert.False(t, location1.Equals(nil))
}
