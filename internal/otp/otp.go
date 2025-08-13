package otp

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type OTPData struct {
	Code       string
	Expiration time.Time
}

type OTPStore struct {
	mu         sync.Mutex
	store      sync.Map
	expiration time.Duration
}

func NewOTPStore() *OTPStore {
	return &OTPStore{
		expiration: 2 * time.Minute,
	}
}

func GenerateOTP() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (s *OTPStore) StoreOTP(phoneNumber, otpCode string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := OTPData{
		Code:       otpCode,
		Expiration: time.Now().Add(s.expiration),
	}
	s.store.Store(phoneNumber, data)

	fmt.Printf("Generated OTP for %s: %s\n", phoneNumber, otpCode)
}

func (s *OTPStore) ValidateOTP(phoneNumber, otpCode string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.store.Load(phoneNumber)
	if !ok {
		return false
	}

	data := value.(OTPData)
	if data.Code != otpCode {
		return false
	}

	if time.Now().After(data.Expiration) {
		s.store.Delete(phoneNumber)
		return false
	}

	s.store.Delete(phoneNumber)
	return true
}
