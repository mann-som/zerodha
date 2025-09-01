package services

import (
	"errors"

	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/repositories"
	"golang.org/x/crypto/bcrypt"
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
	if user.Password == "" {
		return models.User{}, errors.New("password is required")
	}
	if user.Balance < 0 {
		return models.User{}, errors.New("balance cannot be negative")
	}
	if user.Role != "" && user.Role != "user" && user.Role != "admin" {
		return models.User{}, errors.New("role must be 'user' or 'admin'")
	}
	if user.Role == "" {
		user.Role = "user"
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *UserService) GetUser(id string) (models.User, error) {
	if id == "" {
		return models.User{}, errors.New("id is required")
	}
	return s.repo.Get(id)
}

func (s *UserService) ListUsers() ([]models.User, error) {
	return s.repo.List()
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
	if user.Password != "" {
		// Hash new password if provided
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}
	if user.Balance < 0 {
		return models.User{}, errors.New("balance cannot be negative")
	}
	if user.Role != "user" && user.Role != "admin" {
		return models.User{}, errors.New("role must be 'user' or 'admin'")
	}
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}
