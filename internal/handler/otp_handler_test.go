package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-go/internal/cache"
	"user-go/internal/handler"
	"user-go/internal/repository"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, *handler.AuthHandler, *service.OtpService) {
	cache := cache.NewInMemoryCache()
	users := repository.NewInMemoryUserRepository()
	svc := service.NewOtpService(cache, users, "testsecret")
	authHandler := handler.NewAuthHandler(svc)

	r := gin.Default()
	r.POST("/request-otp", authHandler.RequestOTP)
	r.POST("/validate-otp", authHandler.ValidateOTP)
	return r, authHandler, svc
}

func TestRequestOTP_Success(t *testing.T) {
	r, _, _ := setupRouter()

	payload := map[string]string{"phone": "+1234567890"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/request-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "otp")
	assert.Len(t, resp["otp"], 6)
}

func TestValidateOTP_Success(t *testing.T) {
	r, _, svc := setupRouter()

	phone := "+1234567890"
	otp, err := svc.RequestOTP(phone)
	assert.NoError(t, err)

	payload := map[string]string{"phone": phone, "otp": otp}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/validate-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["token"])
}

func TestValidateOTP_Fail(t *testing.T) {
	r, _, _ := setupRouter()

	payload := map[string]string{"phone": "+1234567890", "otp": "000000"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/validate-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
