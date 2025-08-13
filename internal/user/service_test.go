package user

import (
	"errors"
	"testing"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	users []*User
}

func (m *MockRepository) GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	for _, user := range m.users {
		if user.PhoneNumber == phoneNumber {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) RegisterUser(user *User) error {
	m.users = append(m.users, user)
	return nil
}

func (m *MockRepository) GetUserByID(id int) (*User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockRepository) ListUsers(limit, offset int, search string) ([]*User, error) {
	return m.users, nil
}

func TestGetUserByPhoneNumberService(t *testing.T) {
	mockRepo := &MockRepository{
		users: []*User{{ID: 1, PhoneNumber: "+989123456789"}},
	}
	service := NewService(mockRepo)

	user, err := service.GetUserByPhoneNumber("+989123456789")
	if err != nil {
		t.Errorf("expected no error, but got '%s'", err)
	}
	if user == nil {
		t.Error("expected user, but got nil")
	}

	user, err = service.GetUserByPhoneNumber("+989121111111")
	if err != nil {
		t.Errorf("expected no error, but got '%s'", err)
	}
	if user != nil {
		t.Error("expected nil, but got user")
	}
}
