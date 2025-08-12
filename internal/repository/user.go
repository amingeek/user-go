package repository

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

type User struct {
	Phone            string
	RegistrationDate time.Time
}

type UserRepository interface {
	GetByPhone(phone string) (*User, error)
	Create(phone string) (*User, error)
	List(offset, limit int, search string) ([]User, error)
}

var ErrUserNotFound = errors.New("user not found")

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]User),
	}
}

func (r *InMemoryUserRepository) GetByPhone(phone string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[phone]
	if !exists {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *InMemoryUserRepository) Create(phone string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[phone]; exists {
		return nil, errors.New("user already exists")
	}

	user := User{
		Phone:            phone,
		RegistrationDate: time.Now(),
	}

	r.users[phone] = user
	return &user, nil
}

func (r *InMemoryUserRepository) List(offset, limit int, search string) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := []User{}
	search = strings.TrimSpace(search)

	for _, user := range r.users {
		if search == "" || strings.Contains(user.Phone, search) {
			result = append(result, user)
		}
	}

	// مرتب‌سازی خروجی بر اساس شماره تلفن به صورت صعودی
	sort.Slice(result, func(i, j int) bool {
		return result[i].Phone < result[j].Phone
	})

	// اعمال صفحه‌بندی
	start := offset
	if start > len(result) {
		start = len(result)
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}
