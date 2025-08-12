package main

import (
	"user-go/internal/cache"
	"user-go/internal/handler"
	"user-go/internal/repository"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cache := cache.NewInMemoryCache()
	userRepo := repository.NewInMemoryUserRepository()
	otpService := service.NewOtpService(cache, userRepo, "your_jwt_secret_here")

	authHandler := handler.NewAuthHandler(otpService)
	userHandler := handler.NewUserHandler(userRepo)

	// Routes
	api := r.Group("/api")

	api.POST("/auth/request-otp", authHandler.RequestOTP)
	api.POST("/auth/validate-otp", authHandler.ValidateOTP)

	api.GET("/users/:phone", userHandler.GetUser)
	api.GET("/users", userHandler.ListUsers)

	r.Run(":8080")
}
