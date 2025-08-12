package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
	"user-go/internal/cache"
	"user-go/internal/handler"
	"user-go/internal/middleware"
	"user-go/internal/repository"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := os.Getenv("DATABASE_URL") // مثل postgres://user:pass@localhost:5432/dbname
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewPostgresUserRepository(pool)

	r := gin.Default()

	cache := cache.NewInMemoryCache()
	otpService := service.NewOtpService(cache, userRepo, "mysecretjwtkey")

	authHandler := handler.NewAuthHandler(otpService)
	userHandler := handler.NewUserHandler(userRepo)

	r.POST("/auth/request-otp", authHandler.RequestOTP)
	r.POST("/auth/validate-otp", authHandler.ValidateOTP)

	authGroup := r.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware([]byte("mysecretjwtkey")))
	{
		authGroup.GET("/profile", userHandler.GetProfile)
		authGroup.GET("/users/:phone", userHandler.GetUser)
		authGroup.GET("/users", userHandler.ListUsers)
		authGroup.PUT("/users/:phone", userHandler.EditUser)
		authGroup.DELETE("/users/:phone", userHandler.DeleteUser)
	}

	r.Run(":8080")
}
