package cache

// Cache abstracts a simple Redis-like behaviour needed for OTP and rate-limit.
// We'll mock this in unit tests and implement a real Redis-backed version later.
type Cache interface {
    IncrWithExpire(key string, expireSeconds int) (int, error)
    SetWithTTL(key string, value string, ttlSeconds int) error
}
