package services

import (
	"errors"

	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user models.User) (models.User, error) {
	if user.Email == "" {
		return models.User{}, errors.New("email is required")
	}
	if user.Name == "" {
		return models.User{}, errors.New("name is required")
	}
	if user.Balance < 0 {
		return models.User{}, errors.New("balance can't be negative")
	}
	return s.repo.Create(user)
}
