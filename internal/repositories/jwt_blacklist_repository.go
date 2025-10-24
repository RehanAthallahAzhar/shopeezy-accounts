package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/pkg/redisclient"

	"github.com/go-redis/redis/v8"
)

type JWTBlacklistRepository interface {
	AddToBlacklist(ctx context.Context, jti string, expiration time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

type jwtBlacklistRepository struct {
	redisClient *redisclient.RedisClient
}

func NewJWTBlacklistRepository(redisClient *redisclient.RedisClient) JWTBlacklistRepository {
	return &jwtBlacklistRepository{redisClient: redisClient}
}

// JWT ID (JTI) to the Redis blacklist with the validity period must match the validity period of the original JWT.
func (r *jwtBlacklistRepository) AddToBlacklist(ctx context.Context, jti string, expiration time.Duration) error {
	// Key in Redis will be "jwt:blacklist:<jti>"
	// Value will be "blacklisted"
	// Expire will automatically delete the key after a certain time.
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	return r.redisClient.Client.Set(ctx, key, "blacklisted", expiration).Err()
}

// IsBlacklisted checks if the JWT ID (JTI) is on the Redis blacklist
func (r *jwtBlacklistRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	_, err := r.redisClient.Client.Get(ctx, key).Result()
	if err == nil {
		return true, nil // Found in blacklist
	}
	if err == redis.Nil { // redis.Nil -> key not found
		return false, nil // Not found in blacklist
	}
	return false, fmt.Errorf("gagal memeriksa blacklist Redis: %w", err)
}
