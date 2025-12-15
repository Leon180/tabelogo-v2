package repository

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/auth/domain/model"
	"github.com/google/uuid"
)

// SessionRepository defines the interface for session persistence
type SessionRepository interface {
	// Create creates a new session
	Create(ctx context.Context, session *model.Session) error

	// GetByID retrieves a session by ID
	GetByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)

	// GetUserSessions retrieves all active sessions for a user
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error)

	// CountUserSessions counts active sessions for a user
	CountUserSessions(ctx context.Context, userID uuid.UUID) (int, error)

	// UpdateActivity updates session's last active timestamp
	UpdateActivity(ctx context.Context, sessionID uuid.UUID) error

	// Revoke revokes a specific session
	Revoke(ctx context.Context, sessionID uuid.UUID) error

	// RevokeAllForUser revokes all sessions for a user
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired deletes all expired sessions
	DeleteExpired(ctx context.Context) error
}
