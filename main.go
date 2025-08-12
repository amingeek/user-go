// main.go (یا هر جایی که gin router ساخته شده)
package main

import (
	"user-go/internal/cache"
	"user-go/internal/handler"
	"user-go/internal/middleware"
	"user-go/internal/repository"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cache := cache.NewInMemoryCache()
	userRepo := repository.NewInMemoryUserRepository()
	otpService := service.NewOtpService(cache, userRepo, "mysecretjwtkey")

	authHandler := handler.NewAuthHandler(otpService)
	userHandler := handler.NewUserHandler(userRepo)

	r.POST("/auth/request-otp", authHandler.RequestOTP)
	r.POST("/auth/validate-otp", authHandler.ValidateOTP)

	// گروه روت‌های محافظت شده با JWT
	authGroup := r.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware([]byte("mysecretjwtkey")))
	{
		authGroup.GET("/profile", userHandler.GetProfile)
	}

	r.Run(":8080")
}
