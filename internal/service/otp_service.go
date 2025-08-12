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

func (s *OtpService) ValidateOTP(phone, otp string) (string, error) {
	otpKey := "otp:" + phone
	stored, err := s.cache.Get(otpKey)
	if err != nil {
		return "", err
	}
	if stored != otp {
		return "", errors.New("invalid OTP")
	}

	_ = s.cache.Delete(otpKey)

	user, err := s.users.GetByPhone(phone)
	if err == repository.ErrUserNotFound {
		user, err = s.users.Create(phone)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": user.Phone,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
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

func (s *OtpService) RequestOTP(phone string) (string, error) {
	reqKey := "otp_req:" + phone
	count, err := s.cache.IncrWithExpire(reqKey, 600)
	if err != nil {
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
		return "", err
	}

	fmt.Printf("Generated OTP for %s: %s (expires in 2 min)\n", phone, otp)

	return otp, nil
}
