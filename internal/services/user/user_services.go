package user

import (
	"context"
	utils "konsera-backend/internal/utils"
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
	

}