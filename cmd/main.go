package main

import (
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		log.Println("Falling back to environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT not set, using default: 8080")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set in .env or environment variables")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL database")

	if err := db.AutoMigrate(&models.User{}, &models.Order{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	r := gin.Default()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Trading API is up!")
	})

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders/:id", orderHandler.GetOrder)
	r.PUT("/orders/:id", orderHandler.UpdateOrder)
	r.DELETE("/orders/:id", orderHandler.DeleteOrder)

	fmt.Printf("Server starting on :%s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
