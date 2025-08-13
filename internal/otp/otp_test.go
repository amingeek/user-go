package otp

import (
	"testing"
	//"time"
)

func TestGenerateOTP(t *testing.T) {
	otpCode := GenerateOTP()
	if len(otpCode) != 6 {
		t.Errorf("Expected OTP length of 6, but got %d", len(otpCode))
	}
}

func TestStoreAndValidateOTP(t *testing.T) {
	otpStore := NewOTPStore()
	phoneNumber := "+989123456789"
	otpCode := "123456"

	otpStore.StoreOTP(phoneNumber, otpCode)
	isValid := otpStore.ValidateOTP(phoneNumber, otpCode)
	if !isValid {
		t.Error("Expected OTP to be valid, but it was not")
	}

	isValid = otpStore.ValidateOTP(phoneNumber, "987654")
	if isValid {
		t.Error("Expected OTP to be invalid, but it was valid")
	}
}
