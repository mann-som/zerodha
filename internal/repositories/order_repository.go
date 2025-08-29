package repositories

import (
	"errors"

	"github.com/mann-som/zerodha/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order models.Order) (models.Order, error) {
	if order.ID != "" {
		return models.Order{}, errors.New("ID should be empty")
	}
	result := r.db.Create(&order)
	if result.Error != nil {
		return models.Order{}, result.Error
	}
	return order, nil
}

func (r *OrderRepository) Get(id string) (models.Order, error) {
	var order models.Order
	result := r.db.First(&order, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Order{}, errors.New("order not found")
		}
		return models.Order{}, result.Error
	}
	return order, nil
}

func (r *OrderRepository) Update(order models.Order) (models.Order, error) {
	result := r.db.Save(&order)
	if result.Error != nil {
		return models.Order{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *OrderRepository) Delete(id string) error {
	result := r.db.Delete(&models.Order{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *OrderRepository) List() ([]models.Order, error) {
	var orders []models.Order
	result := r.db.Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
