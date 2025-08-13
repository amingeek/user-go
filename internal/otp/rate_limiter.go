package otp

import (
	"sync"
	"time"
)

type RequestRecord struct {
	Timestamps []time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	limits  sync.Map
	window  time.Duration
	maxReqs int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		window:  10 * time.Minute,
		maxReqs: 3,
	}
}

func (r *RateLimiter) Allow(phoneNumber string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	value, _ := r.limits.LoadOrStore(phoneNumber, &RequestRecord{})
	record := value.(*RequestRecord)

	now := time.Now()
	var newTimestamps []time.Time
	for _, t := range record.Timestamps {
		if now.Sub(t) <= r.window {
			newTimestamps = append(newTimestamps, t)
		}
	}
	record.Timestamps = newTimestamps

	if len(record.Timestamps) >= r.maxReqs {
		return false
	}

	record.Timestamps = append(record.Timestamps, now)
	r.limits.Store(phoneNumber, record)

	return true
}
