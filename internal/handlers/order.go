package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/services"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	createdOrder, err := h.service.CreateOrder(order)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, createdOrder)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.service.GetOrder(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	order.ID = c.Param("id")
	updatedOrder, err := h.service.UpdateOrder(order)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, updatedOrder)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteOrder(id); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "order deleted"})
}
