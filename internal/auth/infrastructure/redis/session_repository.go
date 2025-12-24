package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/auth/domain/model"
	"github.com/Leon180/tabelogo-v2/internal/auth/domain/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisSessionRepository struct {
	client *redis.Client
	logger *zap.Logger
}

func NewSessionRepository(client *redis.Client, logger *zap.Logger) repository.SessionRepository {
	return &RedisSessionRepository{
		client: client,
		logger: logger,
	}
}

// Redis key helpers
func (r *RedisSessionRepository) sessionKey(sessionID uuid.UUID) string {
	return fmt.Sprintf("session:%s", sessionID.String())
}

func (r *RedisSessionRepository) userSessionsKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_sessions:%s", userID.String())
}

// sessionDTO is the data transfer object for Redis storage
type sessionDTO struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	DeviceInfo string    `json:"device_info"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsActive   bool      `json:"is_active"`
}

func (r *RedisSessionRepository) Create(ctx context.Context, session *model.Session) error {
	// Serialize session
	dto := sessionDTO{
		ID:         session.ID().String(),
		UserID:     session.UserID().String(),
		DeviceInfo: session.DeviceInfo(),
		IPAddress:  session.IPAddress(),
		CreatedAt:  session.CreatedAt(),
		LastActive: session.LastActive(),
		ExpiresAt:  session.ExpiresAt(),
		IsActive:   session.IsActive(),
	}

	data, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Calculate TTL
	ttl := time.Until(session.ExpiresAt())
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	// Store session data
	sessionKey := r.sessionKey(session.ID())
	if err := r.client.Set(ctx, sessionKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}

	// Add to user's session set
	userSessionsKey := r.userSessionsKey(session.UserID())
	if err := r.client.SAdd(ctx, userSessionsKey, session.ID().String()).Err(); err != nil {
		return fmt.Errorf("failed to add session to user set: %w", err)
	}

	// Set expiry on user sessions set (same as longest session)
	r.client.Expire(ctx, userSessionsKey, ttl)

	r.logger.Info("Session created",
		zap.String("session_id", session.ID().String()),
		zap.String("user_id", session.UserID().String()),
		zap.String("device", session.DeviceInfo()),
		zap.Duration("ttl", ttl),
	)

	return nil
}

func (r *RedisSessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	sessionKey := r.sessionKey(sessionID)

	data, err := r.client.Get(ctx, sessionKey).Result()
	if err == redis.Nil {
		return nil, nil // Session not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var dto sessionDTO
	if err := json.Unmarshal([]byte(data), &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Reconstruct domain model
	id, _ := uuid.Parse(dto.ID)
	userID, _ := uuid.Parse(dto.UserID)

	return model.ReconstructSession(
		id, userID,
		dto.DeviceInfo, dto.IPAddress,
		dto.CreatedAt, dto.LastActive, dto.ExpiresAt,
		dto.IsActive,
	), nil
}

func (r *RedisSessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	userSessionsKey := r.userSessionsKey(userID)

	// Get all session IDs for user
	sessionIDs, err := r.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	sessions := make([]*model.Session, 0, len(sessionIDs))
	for _, idStr := range sessionIDs {
		sessionID, err := uuid.Parse(idStr)
		if err != nil {
			r.logger.Warn("Invalid session ID in user set",
				zap.String("session_id", idStr),
				zap.Error(err),
			)
			// Remove invalid session ID from set
			r.client.SRem(ctx, userSessionsKey, idStr)
			continue
		}

		session, err := r.GetByID(ctx, sessionID)
		if err != nil {
			r.logger.Error("Failed to get session",
				zap.String("session_id", sessionID.String()),
				zap.Error(err),
			)
			continue
		}

		if session == nil {
			// Session expired or deleted, remove from set
			r.client.SRem(ctx, userSessionsKey, idStr)
			continue
		}

		// Only return active, non-expired sessions
		if session.IsValid() {
			sessions = append(sessions, session)
		} else {
			// Clean up invalid session
			r.client.SRem(ctx, userSessionsKey, idStr)
		}
	}

	return sessions, nil
}

func (r *RedisSessionRepository) CountUserSessions(ctx context.Context, userID uuid.UUID) (int, error) {
	sessions, err := r.GetUserSessions(ctx, userID)
	if err != nil {
		return 0, err
	}
	return len(sessions), nil
}

func (r *RedisSessionRepository) UpdateActivity(ctx context.Context, sessionID uuid.UUID) error {
	session, err := r.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	session.UpdateActivity()

	// Re-save session with updated activity
	dto := sessionDTO{
		ID:         session.ID().String(),
		UserID:     session.UserID().String(),
		DeviceInfo: session.DeviceInfo(),
		IPAddress:  session.IPAddress(),
		CreatedAt:  session.CreatedAt(),
		LastActive: session.LastActive(),
		ExpiresAt:  session.ExpiresAt(),
		IsActive:   session.IsActive(),
	}

	data, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	sessionKey := r.sessionKey(sessionID)
	ttl := time.Until(session.ExpiresAt())

	return r.client.Set(ctx, sessionKey, data, ttl).Err()
}

func (r *RedisSessionRepository) Revoke(ctx context.Context, sessionID uuid.UUID) error {
	session, err := r.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	session.Revoke()

	// Update session in Redis (mark as inactive)
	dto := sessionDTO{
		ID:         session.ID().String(),
		UserID:     session.UserID().String(),
		DeviceInfo: session.DeviceInfo(),
		IPAddress:  session.IPAddress(),
		CreatedAt:  session.CreatedAt(),
		LastActive: session.LastActive(),
		ExpiresAt:  session.ExpiresAt(),
		IsActive:   false,
	}

	data, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	sessionKey := r.sessionKey(sessionID)
	ttl := time.Until(session.ExpiresAt())

	// Keep for audit trail but mark as inactive
	if err := r.client.Set(ctx, sessionKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	// Remove from user's active sessions
	userSessionsKey := r.userSessionsKey(session.UserID())
	r.client.SRem(ctx, userSessionsKey, sessionID.String())

	r.logger.Info("Session revoked",
		zap.String("session_id", sessionID.String()),
		zap.String("user_id", session.UserID().String()),
	)

	return nil
}

func (r *RedisSessionRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	sessions, err := r.GetUserSessions(ctx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := r.Revoke(ctx, session.ID()); err != nil {
			r.logger.Error("Failed to revoke session",
				zap.String("session_id", session.ID().String()),
				zap.Error(err),
			)
		}
	}

	r.logger.Info("All sessions revoked for user",
		zap.String("user_id", userID.String()),
		zap.Int("count", len(sessions)),
	)

	return nil
}

func (r *RedisSessionRepository) DeleteExpired(ctx context.Context) error {
	// Redis TTL handles this automatically
	// This method is for manual cleanup if needed
	r.logger.Info("Expired sessions are automatically cleaned up by Redis TTL")
	return nil
}
