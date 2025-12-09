// Package database provides Redis connection and client management
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis client
type RedisClient struct {
	Client *redis.Client
}

// ConnectRedis establishes a connection to Redis
func ConnectRedis(url string) (*RedisClient, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Configure connection pool settings
	opts.PoolSize = 10
	opts.MinIdleConns = 5
	opts.ConnMaxIdleTime = 5 * time.Minute
	opts.ConnMaxLifetime = 30 * time.Minute

	client := redis.NewClient(opts)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &RedisClient{
		Client: client,
	}, nil
}

// Close closes the Redis connection
func (rc *RedisClient) Close() error {
	if rc.Client != nil {
		return rc.Client.Close()
	}
	return nil
}

// Get returns the value for a key
func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// Set sets a key-value pair with optional expiration
func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := rc.Client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Delete deletes one or more keys
func (rc *RedisClient) Delete(ctx context.Context, keys ...string) error {
	if err := rc.Client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}
	return nil
}

// Exists checks if a key exists
func (rc *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rc.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return count > 0, nil
}

// SetNX sets a key only if it does not exist (atomic operation)
func (rc *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := rc.Client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set NX key %s: %w", key, err)
	}
	return result, nil
}

// Expire sets an expiration time on a key
func (rc *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := rc.Client.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration on key %s: %w", key, err)
	}
	return nil
}

// Increment atomically increments a key's value by 1
// Returns the new value after increment
func (rc *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	val, err := rc.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return val, nil
}

// IncrementBy atomically increments a key's value by the specified amount
// Returns the new value after increment
func (rc *RedisClient) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := rc.Client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s by %d: %w", key, value, err)
	}
	return val, nil
}
