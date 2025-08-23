package handlers

import (
	// "encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mann-som/zerodha/internal/models"
)

func UserHandler(c *gin.Context) {
	user := models.User{ID: "USER123", Balance: 1000.50}
	c.IndentedJSON(http.StatusOK, user)
}
