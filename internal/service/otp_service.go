package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"time"
	"user-go/internal/cache"
	"user-go/internal/repository"
)

var (
	ErrRateLimited  = errors.New("too many OTP requests, please wait")
	ErrInvalidToken = errors.New("invalid token")
)

type OtpService struct {
	cache     cache.Cache
	users     repository.UserRepository
	jwtSecret []byte
}

func NewOtpService(c cache.Cache, u repository.UserRepository, secret string) *OtpService {
	return &OtpService{
		cache:     c,
		users:     u,
		jwtSecret: []byte(secret),
	}
}

// ValidateOTP reads OTP from cache, compares, creates user if needed and returns signed JWT.
// Added diagnostic logs to help debug integration tests.
func (s *OtpService) ValidateOTP(phone, otp string) (string, error) {
	otpKey := "otp:" + phone
	stored, err := s.cache.Get(otpKey)
	if err != nil {
		// واضح و قابل تشخیص در لاگ‌ها
		fmt.Printf("[OtpService] cache.Get error for key=%s: %v\n", otpKey, err)
		return "", errors.New("OTP not found or expired")
	}

	// لاگ مقدار ذخیره‌شده و مقدار دریافتی برای دیباگ
	fmt.Printf("[OtpService] stored OTP=%q for phone=%s, provided OTP=%q\n", stored, phone, otp)

	if stored != otp {
		fmt.Printf("[OtpService] invalid otp: stored=%q provided=%q\n", stored, otp)
		return "", errors.New("invalid OTP")
	}

	// حذف OTP بعد از استفاده (لاگ در صورت خطا)
	if err := s.cache.Delete(otpKey); err != nil {
		fmt.Printf("[OtpService] warning: failed to delete otp key %s: %v\n", otpKey, err)
	}

	// ثبت‌نام یا فراخوانی یوزر
	user, err := s.users.GetByPhone(phone)
	if err == repository.ErrUserNotFound {
		user, err = s.users.Create(phone)
		if err != nil {
			fmt.Printf("[OtpService] users.Create error for phone=%s: %v\n", phone, err)
			return "", err
		}
	} else if err != nil {
		fmt.Printf("[OtpService] users.GetByPhone error for phone=%s: %v\n", phone, err)
		return "", err
	}

	// ساخت JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": user.Phone,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		fmt.Printf("[OtpService] token.SignedString error: %v\n", err)
		return "", err
	}
	fmt.Printf("[OtpService] token generated for phone=%s\n", phone)
	return signed, nil
}

func generateOTP() (string, error) {
	max := big.NewInt(int64(1000000))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// RequestOTP generates OTP, rate-limits and stores it in cache
func (s *OtpService) RequestOTP(phone string) (string, error) {
	reqKey := "otp_req:" + phone
	count, err := s.cache.IncrWithExpire(reqKey, 600)
	if err != nil {
		fmt.Printf("[OtpService] IncrWithExpire error for key=%s: %v\n", reqKey, err)
		return "", err
	}

	if count > 3 {
		return "", ErrRateLimited
	}

	otp, err := generateOTP()
	if err != nil {
		return "", err
	}

	otpKey := "otp:" + phone
	if err := s.cache.SetWithTTL(otpKey, otp, 120); err != nil {
		fmt.Printf("[OtpService] SetWithTTL error for key=%s: %v\n", otpKey, err)
		return "", err
	}

	// Show OTP in stdout for tests/debug (no SMS)
	fmt.Printf("Generated OTP for %s: %s (expires in 2 min)\n", phone, otp)

	return otp, nil
}
