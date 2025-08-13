package main

import (
	"context"
	"log"
	"os"
	"time"
	"user-go/internal/cache"
	"user-go/internal/handler"
	"user-go/internal/middleware"
	"user-go/internal/repository"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		// fallback برای توسعه محلی — در production حتماً مقداردهی کن
		secretKey = "mysecretjwtkey"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewPostgresUserRepository(pool)
	cache := cache.NewInMemoryCache()
	otpService := service.NewOtpService(cache, userRepo, secretKey)

	authHandler := handler.NewAuthHandler(otpService)
	userHandler := handler.NewUserHandler(userRepo)

	r := gin.Default()

	// Public routes
	r.POST("/auth/request-otp", authHandler.RequestOTP)
	r.POST("/auth/validate-otp", authHandler.ValidateOTP)

	// Protected routes (با JWT middleware)
	authGroup := r.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware([]byte(secretKey)))
	{
		authGroup.GET("/profile", userHandler.GetProfile)
		authGroup.GET("/users/:phone", userHandler.GetUser)
		authGroup.GET("/users", userHandler.ListUsers)
		authGroup.PUT("/users/:phone", userHandler.EditUser)
		authGroup.DELETE("/users/:phone", userHandler.DeleteUser)
	}

	log.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
