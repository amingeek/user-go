package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"user-go/internal/config"
	"user-go/internal/handler"
	"user-go/internal/middleware"
	"user-go/internal/repository/postgres"
	"user-go/internal/repository/redis"
	"user-go/internal/service"
	"user-go/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	// Connect to PostgreSQL
	pgPool, err := pgxpool.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pgPool.Close()

	// Run migrations
	if err := runMigrations(cfg.PostgresURL); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// Connect to Redis
	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pgPool)
	otpRepo := redis.NewOTPRepository(redisClient, cfg.OTPExpiry)
	rateLimitRepo := redis.NewRateLimitRepository(redisClient, cfg.RateLimitWindow)

	// Initialize services
	otpService := service.NewOTPService(otpRepo)
	rateLimitService := service.NewRateLimitService(rateLimitRepo, cfg.RateLimitWindow)
	authService := service.NewAuthService(userRepo, otpService, cfg.JWTSecret)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, rateLimitService)
	userHandler := handler.NewUserHandler(userService)

	// Initialize middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(rateLimitService)
	jwtMiddleware := handler.JWTMiddleware(cfg.JWTSecret)

	// Create router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(rateLimitMiddleware.Handler) // Apply global rate limiting

	// Auth routes
	r.Post("/auth/send-otp", authHandler.SendOTP)
	r.Post("/auth/verify-otp", authHandler.VerifyOTP)

	// User routes (protected)
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{phone}", userHandler.GetUser)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info().Msgf("Server running on port %d", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

func runMigrations(postgresURL string) error {
	// In a real project, use a proper migration tool like golang-migrate
	// Here we assume the table is created by the docker-compose or manually
	return nil
}
