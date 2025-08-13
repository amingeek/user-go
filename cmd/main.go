package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amingeek/user-go-complete/internal/api"
	"github.com/amingeek/user-go-complete/internal/jwt"
	"github.com/amingeek/user-go-complete/internal/otp"
	"github.com/amingeek/user-go-complete/internal/user"
	"github.com/amingeek/user-go-complete/pkg/database"
)

func main() {
	db, err := database.NewPostgresDB("postgres://user:password@db:5432/userdb?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = database.CreateUserTable(db)
	if err != nil {
		log.Fatalf("Failed to create user table: %v", err)
	}

	userRepo := user.NewPostgresRepository(db)
	userService := user.NewService(userRepo)
	otpStore := otp.NewOTPStore()
	otpLimiter := otp.NewRateLimiter()
	jwtManager := jwt.NewJWTManager("your-super-secret-key", 24*time.Hour)

	handler := api.NewHandler(userService, otpStore, otpLimiter, jwtManager)
	router := api.NewRouter(handler, jwtManager)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// end
