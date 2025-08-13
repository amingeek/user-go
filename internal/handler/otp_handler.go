package handler

import (
	"net/http"
	"user-go/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	otpService *service.OtpService
}

func NewAuthHandler(otpService *service.OtpService) *AuthHandler {
	return &AuthHandler{otpService: otpService}
}

// Request OTP
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone is required"})
		return
	}

	otp, err := h.otpService.RequestOTP(req.Phone)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent (printed on server)", "otp": otp})
}

// Validate OTP and login/register
func (h *AuthHandler) ValidateOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		OTP   string `json:"otp" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and 6-digit otp required"})
		return
	}

	token, err := h.otpService.ValidateOTP(req.Phone, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
