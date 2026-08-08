package user

import (
	"context"
	"fmt"
	dto "konsera-backend/internal/DTO/user"
	userRepo "konsera-backend/internal/repository/user"
)

type UserService struct {
	repo *userRepo.UserRepository
}

func NewUserService(repo *userRepo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *dto.UserRequest) (*dto.UserResponse, error) {
	// validate user input
	if user.Email == "" || user.Phone == "" || user.Password == "" {
		return nil, fmt.Errorf("[ERROR] Missing required fields")
	}

	return nil, nil
}