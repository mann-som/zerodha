package main

import (

	// "net/http"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mann-som/zerodha/internal/handlers"
	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/repositories"
	"github.com/mann-som/zerodha/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading env file :%v", err)
		log.Println("No env file found. using default port")
	} else {
		log.Println("env. loaded successfully")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set in .env")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := gin.Default()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "Tradingin API is up\n")
	})

	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	fmt.Printf("Server starting on :%s Port \n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
