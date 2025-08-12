package internal_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"user-go/internal/cache"
	"user-go/internal/handler"
	"user-go/internal/middleware"
	"user-go/internal/repository"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	cache := cache.NewInMemoryCache()
	userRepo := repository.NewInMemoryUserRepository()
	otpService := service.NewOtpService(cache, userRepo, "testsecretkey")

	authHandler := handler.NewAuthHandler(otpService)
	userHandler := handler.NewUserHandler(userRepo)

	r.POST("/auth/request-otp", authHandler.RequestOTP)
	r.POST("/auth/validate-otp", authHandler.ValidateOTP)

	authGroup := r.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware([]byte("testsecretkey")))
	{
		authGroup.GET("/profile", userHandler.GetProfile)
	}

	return r
}

func TestEndToEnd(t *testing.T) {
	r := setupRouter() // تابعی که router اصلی رو با همه چیز راه اندازی میکنه

	phone := "+1234567890"

	// درخواست OTP
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
	require.NotEmpty(t, otp)

	// اعتبارسنجی OTP و دریافت توکن JWT
	reqBody = strings.NewReader(fmt.Sprintf(`{"phone":"%s", "otp":"%s"}`, phone, otp))
	req = httptest.NewRequest("POST", "/auth/validate-otp", reqBody)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var valResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &valResp)
	require.NoError(t, err)
	token := valResp["token"]
	require.NotEmpty(t, token)

	// درخواست پروفایل با هدر Authorization
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
