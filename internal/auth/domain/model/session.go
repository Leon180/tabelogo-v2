package model

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session
type Session struct {
	id         uuid.UUID
	userID     uuid.UUID
	deviceInfo string
	ipAddress  string
	createdAt  time.Time
	lastActive time.Time
	expiresAt  time.Time
	isActive   bool
}

// NewSession creates a standard session (24 hours)
func NewSession(userID uuid.UUID, deviceInfo, ipAddress string) *Session {
	now := time.Now()
	return &Session{
		id:         uuid.New(),
		userID:     userID,
		deviceInfo: deviceInfo,
		ipAddress:  ipAddress,
		createdAt:  now,
		lastActive: now,
		expiresAt:  now.Add(24 * time.Hour),
		isActive:   true,
	}
}

// NewRememberMeSession creates an extended session (30 days)
func NewRememberMeSession(userID uuid.UUID, deviceInfo, ipAddress string) *Session {
	now := time.Now()
	return &Session{
		id:         uuid.New(),
		userID:     userID,
		deviceInfo: deviceInfo,
		ipAddress:  ipAddress,
		createdAt:  now,
		lastActive: now,
		expiresAt:  now.Add(30 * 24 * time.Hour), // 30 days
		isActive:   true,
	}
}

// ReconstructSession reconstructs a session from persistence
func ReconstructSession(
	id, userID uuid.UUID,
	deviceInfo, ipAddress string,
	createdAt, lastActive, expiresAt time.Time,
	isActive bool,
) *Session {
	return &Session{
		id:         id,
		userID:     userID,
		deviceInfo: deviceInfo,
		ipAddress:  ipAddress,
		createdAt:  createdAt,
		lastActive: lastActive,
		expiresAt:  expiresAt,
		isActive:   isActive,
	}
}

// Getters
func (s *Session) ID() uuid.UUID         { return s.id }
func (s *Session) UserID() uuid.UUID     { return s.userID }
func (s *Session) DeviceInfo() string    { return s.deviceInfo }
func (s *Session) IPAddress() string     { return s.ipAddress }
func (s *Session) CreatedAt() time.Time  { return s.createdAt }
func (s *Session) LastActive() time.Time { return s.lastActive }
func (s *Session) ExpiresAt() time.Time  { return s.expiresAt }
func (s *Session) IsActive() bool        { return s.isActive }

// Domain Methods

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// UpdateActivity updates the last active timestamp
func (s *Session) UpdateActivity() {
	s.lastActive = time.Now()
}

// Revoke marks the session as inactive
func (s *Session) Revoke() {
	s.isActive = false
}

// IsValid checks if the session is valid (active and not expired)
func (s *Session) IsValid() bool {
	return s.isActive && !s.IsExpired()
}
