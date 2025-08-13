package internal_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"user-go/internal/cache"
	"user-go/internal/handler"
	"user-go/internal/middleware"
	"user-go/internal/repository"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// secretKey باید با main.go هماهنگ باشد یا از env خوانده شود
var secretKey = func() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "mysecretjwtkey"
}()

func setupRouterWithPostgres(t *testing.T) *gin.Engine {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set, skipping Postgres integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err, "failed to connect to Postgres")

	userRepo := repository.NewPostgresUserRepository(pool)

	r := gin.Default()

	cache := cache.NewInMemoryCache()
	otpService := service.NewOtpService(cache, userRepo, secretKey)

	authHandler := handler.NewAuthHandler(otpService)
	userHandler := handler.NewUserHandler(userRepo)

	r.POST("/auth/request-otp", authHandler.RequestOTP)
	r.POST("/auth/validate-otp", authHandler.ValidateOTP)

	authGroup := r.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware([]byte(secretKey)))
	{
		authGroup.GET("/profile", userHandler.GetProfile)
	}

	return r
}

func TestEndToEndPostgres(t *testing.T) {
	r := setupRouterWithPostgres(t)

	phone := "+19998887777"

	// 1. Request OTP
	reqBody := strings.NewReader(fmt.Sprintf(`{"phone":"%s"}`, phone))
	req := httptest.NewRequest("POST", "/auth/request-otp", reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var otpResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &otpResp)
	require.NoError(t, err)
	otp := otpResp["otp"]
	// لاگ OTP برای تطبیق با لاگ سرور
	t.Logf("Test received OTP: %s", otp)
	require.NotEmpty(t, otp)

	// 2. Validate OTP → Get token
	reqBody = strings.NewReader(fmt.Sprintf(`{"phone":"%s", "otp":"%s"}`, phone, otp))
	req = httptest.NewRequest("POST", "/auth/validate-otp", reqBody)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// اگر وضعیت موفق نبود، بدنه را چاپ کن تا دلیل دقیق را ببینی
	if w.Code != http.StatusOK {
		t.Fatalf("validate-otp failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var valResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &valResp)
	require.NoError(t, err)
	token := valResp["token"]
	require.NotEmpty(t, token)

	// 3. Request profile
	req = httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var userResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &userResp)
	require.NoError(t, err)

	phoneResp, ok := userResp["Phone"].(string)
	require.True(t, ok)
	assert.Equal(t, phone, phoneResp)
}
