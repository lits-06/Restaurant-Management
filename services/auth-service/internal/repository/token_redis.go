package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"restaurant-management/services/auth-service/internal/domain"
)

// RedisRefreshTokenRepository is a Redis implementation of RefreshTokenRepository.
type RedisRefreshTokenRepository struct {
	client *redis.Client
}

// NewRedisRefreshTokenRepository creates a new Redis-backed refresh token repository.
func NewRedisRefreshTokenRepository(client *redis.Client) *RedisRefreshTokenRepository {
	return &RedisRefreshTokenRepository{client: client}
}

func (r *RedisRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	if err := token.Validate(); err != nil {
		return err
	}

	tokenTTL := time.Until(token.ExpiresAt)
	if tokenTTL <= 0 {
		return domain.ErrTokenExpired
	}

	tokenKey := refreshTokenKey(token.Token)

	if err := r.client.HSet(ctx, tokenKey,
		"user_id", token.UserID,
		"token", token.Token,
		"expires_at", token.ExpiresAt.Unix(),
		"created_at", token.CreatedAt.Unix(),
		"is_revoked", boolToRedisInt(token.IsRevoked),
	).Err(); err != nil {
		return err
	}

	return r.client.Expire(
		ctx,
		tokenKey,
		time.Until(token.ExpiresAt),
	).Err()
}

func (r *RedisRefreshTokenRepository) Update(ctx context.Context, token *domain.RefreshToken) error {
	if err := token.Validate(); err != nil {
		return err
	}

	tokenKey := refreshTokenKey(token.Token)
	exists, err := r.client.Exists(ctx, tokenKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return domain.ErrTokenNotFound
	}

	tokenTTL := time.Until(token.ExpiresAt)
	if tokenTTL <= 0 {
		return domain.ErrTokenExpired
	}

	if err := r.client.HSet(ctx, tokenKey,
		"user_id", token.UserID,
		"token", token.Token,
		"expires_at", token.ExpiresAt.Unix(),
		"created_at", token.CreatedAt.Unix(),
		"is_revoked", boolToRedisInt(token.IsRevoked),
	).Err(); err != nil {
		return err
	}

	return r.client.Expire(ctx, tokenKey, tokenTTL).Err()
}

func (r *RedisRefreshTokenRepository) Delete(ctx context.Context, token string) error {
	deleted, err := r.client.Del(ctx, refreshTokenKey(token)).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return domain.ErrTokenNotFound
	}

	return nil
}

func (r *RedisRefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	values, err := r.client.HGetAll(ctx, refreshTokenKey(token)).Result()
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, domain.ErrTokenNotFound
	}

	expiresAtUnix, err := strconv.ParseInt(values["expires_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_at for token %s: %w", token, err)
	}

	createdAtUnix, err := strconv.ParseInt(values["created_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at for token %s: %w", token, err)
	}

	isRevoked, err := redisIntToBool(values["is_revoked"])
	if err != nil {
		return nil, fmt.Errorf("invalid is_revoked for token %s: %w", token, err)
	}

	return &domain.RefreshToken{
		UserID:    values["user_id"],
		Token:     values["token"],
		ExpiresAt: time.Unix(expiresAtUnix, 0),
		CreatedAt: time.Unix(createdAtUnix, 0),
		IsRevoked: isRevoked,
	}, nil
}

func refreshTokenKey(token string) string {
	return "refresh_token:" + token
}

func boolToRedisInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func redisIntToBool(value string) (bool, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return false, err
	}

	return parsed == 1, nil
}
