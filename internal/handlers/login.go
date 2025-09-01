package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mann-som/zerodha/internal/services"
)

type LoginHandler struct {
	service *services.LoginService
}

func NewLoginHandler(service *services.LoginService) *LoginHandler {
	return &LoginHandler{service: service}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := h.service.Authenticate(loginReq.Email, loginReq.Password)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}
