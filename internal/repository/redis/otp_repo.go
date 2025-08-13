package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	otpPrefix        = "otp:"
	otpAttemptPrefix = "otp_attempt:"
)

type OTPRepository struct {
	client *redis.Client
	expiry time.Duration
}

func NewOTPRepository(client *redis.Client, expiry time.Duration) *OTPRepository {
	return &OTPRepository{
		client: client,
		expiry: expiry,
	}
}

func (r *OTPRepository) StoreOTP(ctx context.Context, phone, otp string) error {
	key := otpPrefix + phone
	err := r.client.Set(ctx, key, otp, r.expiry).Err()
	return err
}

func (r *OTPRepository) GetOTP(ctx context.Context, phone string) (string, time.Time, error) {
	key := otpPrefix + phone
	otp, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", time.Time{}, err
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return "", time.Time{}, err
	}

	expiry := time.Now().Add(ttl)
	return otp, expiry, nil
}

func (r *OTPRepository) GetOTPAttempts(ctx context.Context, phone string) (int, error) {
	key := otpAttemptPrefix + phone
	attempts, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return attempts, err
}

func (r *OTPRepository) IncrementOTPAttempt(ctx context.Context, phone string) error {
	key := otpAttemptPrefix + phone
	_, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	r.client.Expire(ctx, key, 30*time.Minute)
	return nil
}

func (r *OTPRepository) ResetOTPAttempts(ctx context.Context, phone string) error {
	key := otpAttemptPrefix + phone
	return r.client.Del(ctx, key).Err()
}
