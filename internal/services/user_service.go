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
	if user.ID != "" {
		return models.User{}, errors.New("ID should be empty; it will be auto-generated")
	}
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

func (s *UserService) GetUser(id string) (models.User, error) {
	if id == "" {
		return models.User{}, errors.New("id is required")
	}
	return s.repo.Get(id)
}

func (s *UserService) UpdateUser(user models.User) (models.User, error) {
	if user.ID == "" {
		return models.User{}, errors.New("id is required")
	}
	if user.Email == "" {
		return models.User{}, errors.New("email is required")
	}
	if user.Name == "" {
		return models.User{}, errors.New("name is required")
	}
	if user.Balance < 0 {
		return models.User{}, errors.New("balance cannot be negative")
	}
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}
