package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mann-som/zerodha/internal/handlers"
	"github.com/mann-som/zerodha/internal/middleware"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set in .env or environment variables")
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.Static("/ui", "./frontend")

	userRepo := repositories.NewUserRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	userService := services.NewUserService(userRepo)
	loginService := services.NewLoginService(userRepo, jwtSecret)
	orderService := services.NewOrderService(orderRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	loginHandler := handlers.NewLoginHandler(loginService)

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Trading API is up!")
	})

	r.POST("/login", loginHandler.Login)

	api := r.Group("/api", middleware.AuthMiddleware(jwtSecret))
	{

		api.POST("/users", middleware.AdminMiddleware(), userHandler.CreateUser)
		api.GET("/users", middleware.AdminMiddleware(), userHandler.ListUsers)
		api.GET("/users/:id", middleware.AdminMiddleware(), userHandler.GetUser)
		api.PUT("/users/:id", middleware.AdminMiddleware(), userHandler.UpdateUser)
		r.DELETE("/users/:id", middleware.AdminMiddleware(), userHandler.DeleteUser)

		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders", orderHandler.ListOrders)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.PUT("/orders/:id", middleware.AdminMiddleware(), orderHandler.UpdateOrder)
		api.DELETE("/orders/:id", middleware.AdminMiddleware(), orderHandler.DeleteOrder)
	}

	fmt.Printf("Server starting on :%s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
