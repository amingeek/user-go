package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimitRepository struct {
	client *redis.Client
	window time.Duration
}

func NewRateLimitRepository(client *redis.Client, window time.Duration) *RateLimitRepository {
	return &RateLimitRepository{
		client: client,
		window: window,
	}
}

func (r *RateLimitRepository) Increment(ctx context.Context, key string) (int, error) {
	now := time.Now().UnixNano()
	windowStart := now - int64(r.window.Nanoseconds())

	// Add current timestamp to sorted set
	_, err := r.client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: now,
	}).Result()
	if err != nil {
		return 0, err
	}

	// Remove old entries (outside the window)
	_, err = r.client.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10)).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration
	_, err = r.client.Expire(ctx, key, r.window).Result()
	if err != nil {
		return 0, err
	}

	// Get current count
	count, err := r.client.ZCount(ctx, key, strconv.FormatInt(windowStart, 10), "+inf").Result()
	return int(count), err
}

func (r *RateLimitRepository) GetCount(ctx context.Context, key string) (int, error) {
	windowStart := time.Now().Add(-r.window).UnixNano()

	// Remove old entries
	_, err := r.client.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10)).Result()
	if err != nil {
		return 0, err
	}

	// Get current count
	count, err := r.client.ZCard(ctx, key).Result()
	return int(count), err
}
