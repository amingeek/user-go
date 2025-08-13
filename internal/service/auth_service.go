package service

import (
	"context"
	"errors"
	"time"

	"user-go/internal/model"
	"user-go/internal/repository"
	"user-go/internal/utils"
)

var (
	ErrOTPNotFound     = errors.New("OTP not found")
	ErrInvalidOTP      = errors.New("invalid OTP")
	ErrOTPExpired      = errors.New("OTP expired")
	ErrTooManyAttempts = errors.New("too many attempts")
)

type AuthService interface {
	SendOTP(ctx context.Context, phone string) (string, error)
	VerifyOTP(ctx context.Context, phone, otp string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	otpRepo   repository.OTPRepository
	jwtSecret string
}

func NewAuthService(
	userRepo repository.UserRepository,
	otpRepo repository.OTPRepository,
	jwtSecret string,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) SendOTP(ctx context.Context, phone string) (string, error) {
	otp, err := s.otpRepo.GenerateOTP(ctx, phone)
	if err != nil {
		return "", err
	}
	return otp, nil
}

func (s *authService) VerifyOTP(ctx context.Context, phone, otp string) (string, error) {
	// Check OTP attempts
	attempts, err := s.otpRepo.GetOTPAttempts(ctx, phone)
	if err != nil {
		return "", err
	}
	if attempts >= 5 {
		return "", ErrTooManyAttempts
	}

	// Validate OTP
	storedOTP, expiry, err := s.otpRepo.GetOTP(ctx, phone)
	if err != nil {
		return "", ErrOTPNotFound
	}

	if time.Now().After(expiry) {
		return "", ErrOTPExpired
	}

	if storedOTP != otp {
		// Increment attempt counter
		s.otpRepo.IncrementOTPAttempt(ctx, phone)
		return "", ErrInvalidOTP
	}

	// Reset attempts on success
	s.otpRepo.ResetOTPAttempts(ctx, phone)

	// Check if user exists
	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		// Create new user if not exists
		newUser := model.User{
			Phone:     phone,
			CreatedAt: time.Now(),
		}
		user, err = s.userRepo.CreateUser(ctx, newUser)
		if err != nil {
			return "", err
		}
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
