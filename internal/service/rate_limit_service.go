package service

import (
	"context"
	"time"

	"user-go/internal/repository"
)

type RateLimitService interface {
	Check(ctx context.Context, key string) (bool, error)
}

type rateLimitService struct {
	rateLimitRepo repository.RateLimitRepository
	maxRequests   int
	window        time.Duration
}

func NewRateLimitService(rateLimitRepo repository.RateLimitRepository, window time.Duration) RateLimitService {
	return &rateLimitService{
		rateLimitRepo: rateLimitRepo,
		maxRequests:   3, // Max 3 requests per window
		window:        window,
	}
}

func (s *rateLimitService) Check(ctx context.Context, key string) (bool, error) {
	count, err := s.rateLimitRepo.GetCount(ctx, key)
	if err != nil {
		return false, err
	}

	if count >= s.maxRequests {
		return true, nil
	}

	_, err = s.rateLimitRepo.Increment(ctx, key)
	if err != nil {
		return false, err
	}

	return false, nil
}
