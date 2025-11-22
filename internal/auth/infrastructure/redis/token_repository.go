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
)

type tokenRepository struct {
	client *redis.Client
}

// NewTokenRepository creates a new redis token repository
func NewTokenRepository(client *redis.Client) repository.TokenRepository {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) key(tokenHash string) string {
	return fmt.Sprintf("refresh_token:%s", tokenHash)
}

func (r *tokenRepository) userKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_tokens:%s", userID.String())
}

func (r *tokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	dto := toDTO(token)
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()
	pipe.Set(ctx, r.key(token.TokenHash()), data, time.Until(token.ExpiresAt()))
	pipe.SAdd(ctx, r.userKey(token.UserID()), token.TokenHash())
	pipe.Expire(ctx, r.userKey(token.UserID()), 30*24*time.Hour)
	_, err = pipe.Exec(ctx)
	return err
}

// TokenDTO is used for JSON serialization since Domain Entity has private fields
type TokenDTO struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

func toDTO(t *model.RefreshToken) TokenDTO {
	return TokenDTO{
		ID:        t.ID(),
		UserID:    t.UserID(),
		TokenHash: t.TokenHash(),
		ExpiresAt: t.ExpiresAt(),
		CreatedAt: t.CreatedAt(),
		RevokedAt: t.RevokedAt(),
	}
}

func fromDTO(dto TokenDTO) *model.RefreshToken {
	return model.ReconstructRefreshToken(
		dto.ID,
		dto.UserID,
		dto.TokenHash,
		dto.ExpiresAt,
		dto.CreatedAt,
		dto.RevokedAt,
	)
}

func (r *tokenRepository) GetByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	data, err := r.client.Get(ctx, r.key(tokenHash)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Or specific error
		}
		return nil, err
	}

	var dto TokenDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}

	return fromDTO(dto), nil
}

func (r *tokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	// This is tricky with just ID if we key by Hash.
	// The interface asks to revoke by ID, but our storage is by Hash.
	// We might need to change the interface or storage strategy.
	// For now, let's assume we can't easily revoke by ID without an index.
	// But wait, the domain method Revoke() is on the entity.
	// Usually we retrieve the token, call Revoke(), then Save().
	// But here the interface is Revoke(ctx, id).
	// Let's assume we pass the token object to Update, or change interface to Revoke(token).
	// Or we store by ID as well?
	// Given the KV nature, storing by Hash is best for validation.
	// Let's stick to the interface but maybe we need to fetch by ID?
	// Actually, for a refresh token, we usually have the token string when we want to refresh/revoke.
	// If we want to revoke by ID (e.g. admin action), we need a mapping ID -> Hash.

	// For simplicity in this phase, I'll implement RevokeAllForUser correctly,
	// and for Revoke(id), I'll leave a TODO or change interface to Revoke(token *RefreshToken).
	// Actually, let's update the interface in the next step if needed.
	// For now, I'll implement Revoke assuming we might need to look it up or just not support it efficiently yet.
	return fmt.Errorf("revoke by ID not fully supported in redis without index, use RevokeAllForUser")
}

func (r *tokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	// Get all token hashes for user
	hashes, err := r.client.SMembers(ctx, r.userKey(userID)).Result()
	if err != nil {
		return err
	}

	if len(hashes) == 0 {
		return nil
	}

	// Delete all tokens
	pipe := r.client.Pipeline()
	for _, hash := range hashes {
		pipe.Del(ctx, r.key(hash))
	}
	pipe.Del(ctx, r.userKey(userID))
	_, err = pipe.Exec(ctx)
	return err
}
