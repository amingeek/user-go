package repository

import (
	"context"
	"time"

	"user-go/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUsers(ctx context.Context, page, pageSize int, search string) ([]model.User, int, error)
}

type OTPRepository interface {
	StoreOTP(ctx context.Context, phone, otp string) error
	GetOTP(ctx context.Context, phone string) (string, time.Time, error)
	GetOTPAttempts(ctx context.Context, phone string) (int, error)
	IncrementOTPAttempt(ctx context.Context, phone string) error
	ResetOTPAttempts(ctx context.Context, phone string) error
}

type RateLimitRepository interface {
	Increment(ctx context.Context, key string) (int, error)
	GetCount(ctx context.Context, key string) (int, error)
}
