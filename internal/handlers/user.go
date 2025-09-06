package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if user.Email == "" || user.Name == "" || user.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "email, name, and password are required"})
		return
	}
	if user.Role == "" {
		user.Role = "user"
	}
	if user.Balance == 0 {
		user.Balance = 0
	}

	createdUser, err := h.service.CreateUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser.Password = ""
	c.IndentedJSON(http.StatusCreated, createdUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	createdUser, err := h.service.CreateUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser.Password = ""
	c.IndentedJSON(http.StatusCreated, createdUser)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUser(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	c.IndentedJSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	c.IndentedJSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	user.ID = c.Param("id")
	updatedUser, err := h.service.UpdateUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser.Password = ""
	c.IndentedJSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(id); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted"})
}
