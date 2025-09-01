package repositories

import (
	"errors"

	"github.com/mann-som/zerodha/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func (r *UserRepository) Get(id string) (models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepository) Update(user models.User) (models.User, error) {
	result := r.db.Save(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) Delete(id string) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) List() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, result.Error
	}
	return user, nil
}
