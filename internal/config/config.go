package config

import (
	"context"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	PostgresURL     string
	RedisURL        string
	JWTSecret       string
	OTPExpiry       time.Duration
	RateLimitWindow time.Duration
	Context         context.Context
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	otpExpiry, _ := strconv.Atoi(getEnv("OTP_EXPIRY", "120"))
	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "600"))

	return &Config{
		Port:            port,
		PostgresURL:     getEnv("POSTGRES_URL", "postgres://user:pass@localhost:5432/userdb?sslmode=disable"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:       getEnv("JWT_SECRET", "supersecretkey"),
		OTPExpiry:       time.Duration(otpExpiry) * time.Second,
		RateLimitWindow: time.Duration(rateLimitWindow) * time.Second,
		Context:         context.Background(),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
