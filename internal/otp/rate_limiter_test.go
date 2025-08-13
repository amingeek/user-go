package otp

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter()
	phoneNumber := "+989123456789"
	limiter.window = 1 * time.Second

	for i := 0; i < 3; i++ {
		if !limiter.Allow(phoneNumber) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if limiter.Allow(phoneNumber) {
		t.Error("4th request should not be allowed")
	}

	time.Sleep(2 * time.Second)

	if !limiter.Allow(phoneNumber) {
		t.Error("Request should be allowed after rate limit window expires")
	}
}
