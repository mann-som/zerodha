package services

import (
	"errors"
	"time"

	// "github.com/mann-som/zerodha/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mann-som/zerodha/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewLoginService(userRepo *repositories.UserRepository, jwtSecret string) *LoginService {
	return &LoginService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *LoginService) Authenticate(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}
