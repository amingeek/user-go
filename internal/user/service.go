package user

import (
	"database/sql"
	"errors"
)

// ServiceInterface defines the methods required for user business logic.
type ServiceInterface interface {
	GetUserByPhoneNumber(phoneNumber string) (*User, error)
	RegisterUser(user *User) error
	GetUserByID(id int) (*User, error)
	ListUsers(limit, offset int, search string) ([]*User, error)
}

// Service provides business logic for user management.
type Service struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetUserByPhoneNumber retrieves a user by their phone number.
func (s *Service) GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	return s.repo.GetUserByPhoneNumber(phoneNumber)
}

// RegisterUser registers a new user.
func (s *Service) RegisterUser(user *User) error {
	return s.repo.RegisterUser(user)
}

// GetUserByID retrieves a user by their ID.
func (s *Service) GetUserByID(id int) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// ListUsers retrieves a paginated list of users with optional search filter.
func (s *Service) ListUsers(limit, offset int, search string) ([]*User, error) {
	return s.repo.ListUsers(limit, offset, search)
}
