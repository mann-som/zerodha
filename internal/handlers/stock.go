package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/services"
)

type StockHandler struct {
	service *services.StockService
}

func NewStockHandler(service *services.StockService) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) CreateStock(c *gin.Context) {
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	createdStock, err := h.service.CreateStock(stock)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, createdStock)
}

func (h *StockHandler) GetStock(c *gin.Context) {
	id := c.Param("id")
	stock, err := h.service.GetStock(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, stock)
}

func (h *StockHandler) ListStocks(c *gin.Context) {
	stocks, err := h.service.ListStocks()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, stocks)
}

func (h *StockHandler) UpdateStock(c *gin.Context) {
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	stock.ID = c.Param("id")
	updatedStock, err := h.service.UpdateStock(stock)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, updatedStock)
}

func (h *StockHandler) DeleteStock(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteStock(id); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "stock deleted"})
}
