package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/amingeek/user-go-complete/internal/jwt"
	"github.com/amingeek/user-go-complete/internal/otp"
	"github.com/amingeek/user-go-complete/internal/user"
	"github.com/gorilla/mux"
)

type OTPRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type OTPValidationRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID           int       `json:"id"`
	PhoneNumber  string    `json:"phone_number"`
	RegisteredAt time.Time `json:"registered_at"`
}

type Handler struct {
	userService user.ServiceInterface
	otpStore    *otp.OTPStore
	otpLimiter  *otp.RateLimiter
	jwtManager  *jwt.JWTManager
}

func NewHandler(userService user.ServiceInterface, otpStore *otp.OTPStore, otpLimiter *otp.RateLimiter, jwtManager *jwt.JWTManager) *Handler {
	return &Handler{
		userService: userService,
		otpStore:    otpStore,
		otpLimiter:  otpLimiter,
		jwtManager:  jwtManager,
	}
}

func (h *Handler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req OTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !h.otpLimiter.Allow(req.PhoneNumber) {
		http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
		return
	}

	otpCode := otp.GenerateOTP()
	h.otpStore.StoreOTP(req.PhoneNumber, otpCode)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent to console"})
}

func (h *Handler) LoginValidateOTP(w http.ResponseWriter, r *http.Request) {
	var req OTPValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !h.otpStore.ValidateOTP(req.PhoneNumber, req.OTP) {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	existingUser, err := h.userService.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}

	if existingUser == nil {
		newUser := &user.User{
			PhoneNumber:  req.PhoneNumber,
			RegisteredAt: time.Now(),
		}
		if err := h.userService.RegisterUser(newUser); err != nil {
			http.Error(w, "Failed to register new user", http.StatusInternalServerError)
			return
		}
		existingUser = newUser
	}

	token, err := h.jwtManager.GenerateToken(existingUser.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TokenResponse{Token: token})
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	foundUser, err := h.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{
		ID:           foundUser.ID,
		PhoneNumber:  foundUser.PhoneNumber,
		RegisteredAt: foundUser.RegisteredAt,
	})
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	search := r.URL.Query().Get("search")

	users, err := h.userService.ListUsers(limit, offset, search)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
