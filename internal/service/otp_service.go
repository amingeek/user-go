package service

import (
    "errors"
    "time"

    "user-go/internal/cache"
)

// ErrRateLimited is returned when phone exceeded allowed OTP requests.
var ErrRateLimited = errors.New("rate limited")

// OtpService handles OTP generation and rate limiting.
type OtpService struct {
    cache cache.Cache
}

func NewOtpService(c cache.Cache) *OtpService {
    return &OtpService{cache: c}
}

// RequestOTP generates an OTP for a phone number and stores it in cache.
// NOTE: Minimal stub here so tests (written first) can compile. We'll implement
// behaviour in the next TDD step.
func (s *OtpService) RequestOTP(phone string) (string, error) {
    // TODO: implement according to tests
    _ = time.Now()
    return "", nil
}
