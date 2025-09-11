package repositories

import (
	"github.com/mann-som/zerodha/internal/models"
	"gorm.io/gorm"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) Create(stock models.Stock) (models.Stock, error) {
	result := r.db.Create(&stock)
	if result.Error != nil {
		return models.Stock{}, result.Error
	}

	return stock, nil
}

func (r *StockRepository) Get(id string) (models.Stock, error) {
	var stock models.Stock
	result := r.db.First(&stock, "id = ?", id)
	if result.Error != nil {
		return models.Stock{}, result.Error
	}
	return stock, nil
}

func (r *StockRepository) List() ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.Find(&stocks)
	if result.Error != nil {
		return nil, result.Error
	}
	return stocks, nil
}

func (r *StockRepository) Update(stock models.Stock) (models.Stock, error) {
	result := r.db.Save(&stock)
	if result.Error != nil {
		return models.Stock{}, result.Error
	}

	return stock, nil
}

func (r *StockRepository) Delete(id string) error {
	result := r.db.Delete(&models.Stock{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
