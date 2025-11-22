package repository

import (
	"context"

	"github.com/Leon180/tabelogo-v2/internal/auth/domain/model"
	"github.com/google/uuid"
)

// TokenRepository defines the interface for token persistence
type TokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
