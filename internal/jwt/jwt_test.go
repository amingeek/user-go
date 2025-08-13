package jwt

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	jwtManager := NewJWTManager("secret-key-for-test", 24*time.Hour)
	phoneNumber := "+989123456789"

	token, err := jwtManager.GenerateToken(phoneNumber)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Errorf("Error validating token: %v", err)
	}

	if claims.PhoneNumber != phoneNumber {
		t.Errorf("Expected phone number %s, but got %s", phoneNumber, claims.PhoneNumber)
	}
}
