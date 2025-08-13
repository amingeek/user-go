package service

import (
	"context"

	"user-go/internal/model"
	"user-go/internal/repository"
)

type UserService interface {
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUsers(ctx context.Context, page, pageSize int, search string) ([]model.User, int, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	return s.userRepo.GetUserByPhone(ctx, phone)
}

func (s *userService) GetUsers(ctx context.Context, page, pageSize int, search string) ([]model.User, int, error) {
	// Validate and set default values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	return s.userRepo.GetUsers(ctx, page, pageSize, search)
}
