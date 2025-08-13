package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"user-go/internal/service"
	"user-go/internal/utils"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authService      service.AuthService
	rateLimitService service.RateLimitService
}

func NewAuthHandler(authService service.AuthService, rateLimitService service.RateLimitService) *AuthHandler {
	return &AuthHandler{
		authService:      authService,
		rateLimitService: rateLimitService,
	}
}

func (h *AuthHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate phone number
	if !utils.ValidatePhoneNumber(req.Phone) {
		http.Error(w, "Invalid phone number format", http.StatusBadRequest)
		return
	}

	// Check rate limit
	if limited, err := h.rateLimitService.Check(r.Context(), req.Phone); err != nil || limited {
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	otp, err := h.authService.SendOTP(r.Context(), req.Phone)
	if err != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	// In production, we wouldn't return the OTP to the client
	response := map[string]interface{}{
		"message": "OTP sent successfully",
		"otp":     otp, // Only for development
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate phone number
	if !utils.ValidatePhoneNumber(req.Phone) {
		http.Error(w, "Invalid phone number format", http.StatusBadRequest)
		return
	}

	token, err := h.authService.VerifyOTP(r.Context(), req.Phone, req.OTP)
	if err != nil {
		switch err {
		case service.ErrOTPNotFound, service.ErrInvalidOTP:
			http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		case service.ErrOTPExpired:
			http.Error(w, "OTP expired", http.StatusUnauthorized)
		case service.ErrTooManyAttempts:
			http.Error(w, "Too many attempts", http.StatusTooManyRequests)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
