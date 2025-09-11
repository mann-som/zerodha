package services

import (
	"errors"

	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/repositories"
)

type StockService struct {
	repo *repositories.StockRepository
}

func NewStockService(repo *repositories.StockRepository) *StockService {
	return &StockService{repo: repo}
}

func (s *StockService) CreateStock(stock models.Stock) (models.Stock, error) {
	if stock.Symbol == "" {
		return models.Stock{}, errors.New("symbol is required")
	}
	if stock.InitialPrice <= 0 {
		return models.Stock{}, errors.New("initial_price must be positive")
	}
	if stock.Description == "" {
		return models.Stock{}, errors.New("description is required")
	}
	stock.CurrentPrice = stock.InitialPrice // Set current to initial
	return s.repo.Create(stock)
}

func (s *StockService) GetStock(id string) (models.Stock, error) {
	if id == "" {
		return models.Stock{}, errors.New("id is required")
	}
	return s.repo.Get(id)
}

func (s *StockService) ListStocks() ([]models.Stock, error) {
	return s.repo.List()
}

func (s *StockService) UpdateStock(stock models.Stock) (models.Stock, error) {
	if stock.ID == "" {
		return models.Stock{}, errors.New("id is required")
	}
	return s.repo.Update(stock)
}

func (s *StockService) DeleteStock(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}
