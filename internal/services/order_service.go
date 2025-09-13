package services

import (
	"errors"
	"fmt"

	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/repositories"
)

type OrderService struct {
	repo     *repositories.OrderRepository
	userRepo *repositories.UserRepository
}

func NewOrderService(repo *repositories.OrderRepository, userRepo *repositories.UserRepository) *OrderService {
	return &OrderService{repo: repo, userRepo: userRepo}
}

func (s *OrderService) CreateOrder(order models.Order, userID string) (models.Order, error) {
	if userID == "" {
		return models.Order{}, errors.New("user_id is required")
	}
	user, err := s.userRepo.Get(userID)
	if err != nil {
		return models.Order{}, errors.New("user not found")
	}
	order.UserID = userID
	if order.Symbol == "" {
		return models.Order{}, errors.New("symbol is required")
	}
	if order.Side != "buy" && order.Side != "sell" {
		return models.Order{}, errors.New("side must be 'buy' or 'sell'")
	}
	if order.Quantity <= 0 {
		return models.Order{}, errors.New("quantity must be positive")
	}
	if order.Price <= 0 {
		return models.Order{}, errors.New("price must be positive")
	}
	if order.Status != "" && order.Status != "pending" && order.Status != "executed" && order.Status != "cancelled" {
		return models.Order{}, errors.New("status must be 'pending', 'executed', or 'cancelled'")
	}
	if order.Status == "" {
		order.Status = "pending"
	}

	if order.Side == "buy" {
		totalCost := float64(order.Quantity) * order.Price
		if totalCost > user.Balance {
			return models.Order{}, errors.New("insufficient balance: required " + fmt.Sprintf("%.2f", totalCost) + ", available " + fmt.Sprintf("%.2f", user.Balance))
		}
	}
	return s.repo.Create(order)
}

func (s *OrderService) GetOrder(id string) (models.Order, error) {
	if id == "" {
		return models.Order{}, errors.New("id is required")
	}
	return s.repo.Get(id)
}

func (s *OrderService) ListOrders() ([]models.Order, error) {
	return s.repo.List()
}

func (s *OrderService) UpdateOrder(order models.Order, userID string) (models.Order, error) {
	if order.ID == "" {
		return models.Order{}, errors.New("id is required")
	}
	_, err := s.userRepo.Get(userID)
	if err != nil {
		return models.Order{}, errors.New("user not found")
	}
	order.UserID = userID
	if order.Symbol == "" {
		return models.Order{}, errors.New("symbol is required")
	}
	if order.Side != "buy" && order.Side != "sell" {
		return models.Order{}, errors.New("side must be 'buy' or 'sell'")
	}
	if order.Quantity <= 0 {
		return models.Order{}, errors.New("quantity must be positive")
	}
	if order.Price <= 0 {
		return models.Order{}, errors.New("price must be positive")
	}
	if order.Status != "pending" && order.Status != "executed" && order.Status != "cancelled" {
		return models.Order{}, errors.New("status must be 'pending', 'executed', or 'cancelled'")
	}
	return s.repo.Update(order)
}

func (s *OrderService) DeleteOrder(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}
