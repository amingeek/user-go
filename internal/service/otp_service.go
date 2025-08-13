package service

import (
	"context"
	"math/rand"
	"time"

	"user-go/internal/repository"
)

// تعریف مستقیم اینترفیس در صورت مشکل
type OTPRepository interface {
	StoreOTP(ctx context.Context, phone, otp string) error
}

type otpService struct {
	otpRepo OTPRepository
}

func NewOTPService(otpRepo OTPRepository) *otpService {
	return &otpService{
		otpRepo: otpRepo,
	}
}

func (s *otpService) GenerateOTP(ctx context.Context, phone string) (string, error) {
	otp := generateRandomOTP()
	err := s.otpRepo.StoreOTP(ctx, phone, otp)
	return otp, err
}

func generateRandomOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := make([]byte, 6)
	for i := range digits {
		digits[i] = byte(r.Intn(10)) + '0'
	}
	return string(digits)
}
