package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token entity
type RefreshToken struct {
	id        uuid.UUID
	userID    uuid.UUID
	tokenHash string
	expiresAt time.Time
	createdAt time.Time
	revokedAt *time.Time
}

// NewRefreshToken creates a new refresh token
func NewRefreshToken(userID uuid.UUID, tokenHash string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		id:        uuid.New(),
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		createdAt: time.Now(),
	}
}

// ReconstructRefreshToken reconstructs a refresh token from persistence
func ReconstructRefreshToken(
	id uuid.UUID,
	userID uuid.UUID,
	tokenHash string,
	expiresAt, createdAt time.Time,
	revokedAt *time.Time,
) *RefreshToken {
	return &RefreshToken{
		id:        id,
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		createdAt: createdAt,
		revokedAt: revokedAt,
	}
}

// Getters
func (t *RefreshToken) ID() uuid.UUID         { return t.id }
func (t *RefreshToken) UserID() uuid.UUID     { return t.userID }
func (t *RefreshToken) TokenHash() string     { return t.tokenHash }
func (t *RefreshToken) ExpiresAt() time.Time  { return t.expiresAt }
func (t *RefreshToken) CreatedAt() time.Time  { return t.createdAt }
func (t *RefreshToken) RevokedAt() *time.Time { return t.revokedAt }

// Domain Methods

// IsExpired checks if the token is expired
func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.expiresAt)
}

// IsRevoked checks if the token is revoked
func (t *RefreshToken) IsRevoked() bool {
	return t.revokedAt != nil
}

// Revoke marks the token as revoked
func (t *RefreshToken) Revoke() {
	now := time.Now()
	t.revokedAt = &now
}
