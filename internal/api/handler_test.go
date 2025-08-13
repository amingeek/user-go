package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amingeek/user-go-complete/internal/jwt"
	"github.com/amingeek/user-go-complete/internal/otp"
	"github.com/amingeek/user-go-complete/internal/user"
	"github.com/gorilla/mux"
)

// MockUserService is a mock implementation of the user.ServiceInterface.
type MockUserService struct{}

func (m *MockUserService) GetUserByPhoneNumber(phoneNumber string) (*user.User, error) {
	if phoneNumber == "+989123456789" {
		return &user.User{ID: 1, PhoneNumber: phoneNumber}, nil
	}
	return nil, nil
}
func (m *MockUserService) RegisterUser(u *user.User) error { return nil }
func (m *MockUserService) GetUserByID(id int) (*user.User, error) {
	if id == 1 {
		return &user.User{ID: 1, PhoneNumber: "+989123456789"}, nil
	}
	return nil, errors.New("not found")
}
func (m *MockUserService) ListUsers(limit, offset int, search string) ([]*user.User, error) {
	return []*user.User{{ID: 1, PhoneNumber: "+989123456789"}}, nil
}

func TestRequestOTP(t *testing.T) {
	handler := NewHandler(&MockUserService{}, otp.NewOTPStore(), otp.NewRateLimiter(), nil)
	body := []byte(`{"phone_number": "+989123456789"}`)
	req, _ := http.NewRequest("POST", "/otp/request", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.RequestOTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLoginValidateOTP(t *testing.T) {
	otpStore := otp.NewOTPStore()
	phoneNumber := "+989123456789"
	otpCode := otp.GenerateOTP()
	otpStore.StoreOTP(phoneNumber, otpCode)

	handler := NewHandler(&MockUserService{}, otpStore, otp.NewRateLimiter(), jwt.NewJWTManager("test-key", time.Hour))
	reqBody := map[string]string{"phone_number": phoneNumber, "otp": otpCode}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/otp/login", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	handler.LoginValidateOTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var res TokenResponse
	json.Unmarshal(rr.Body.Bytes(), &res)
	if res.Token == "" {
		t.Error("expected a JWT token, but got an empty string")
	}
}

func TestGetUserByID(t *testing.T) {
	handler := NewHandler(&MockUserService{}, nil, nil, nil)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	handler.GetUserByID(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/users/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rr = httptest.NewRecorder()
	handler.GetUserByID(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
