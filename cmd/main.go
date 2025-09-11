package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	if filepath.Base(projectRoot) == "cmd" {
		projectRoot = filepath.Dir(projectRoot)
	}
	envPath := filepath.Join(projectRoot, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		log.Println("Falling back to environment variables")
	} else {
		log.Printf(".env file loaded successfully from %s", envPath)
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

	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.Stock{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	userRepo := repositories.NewUserRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	stockRepo := repositories.NewStockRepository(db)

	userService := services.NewUserService(userRepo)
	loginService := services.NewLoginService(userRepo, jwtSecret)
	orderService := services.NewOrderService(orderRepo, userRepo)
	stockService := services.NewStockService(stockRepo)

	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	stockHandler := handlers.NewStockHandler(stockService)
	loginHandler := handlers.NewLoginHandler(loginService)

	r.POST("/register", userHandler.Register)
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Trading API is up!\n")
	})

	r.POST("/login", loginHandler.Login)

	// Public stock routes (all users can view stocks)
	apiPublic := r.Group("/api")
	{
		apiPublic.GET("/stocks", stockHandler.ListStocks)
		apiPublic.GET("/stocks/:id", stockHandler.GetStock)
	}

	// Protected routes
	api := r.Group("/api", middleware.AuthMiddleware(jwtSecret))
	{
		api.POST("/users", middleware.AdminMiddleware(), userHandler.CreateUser)
		api.GET("/users", middleware.AdminMiddleware(), userHandler.ListUsers)
		api.GET("/users/:id", middleware.AdminMiddleware(), userHandler.GetUser)
		api.PUT("/users/:id", middleware.AdminMiddleware(), userHandler.UpdateUser)
		api.DELETE("/users/:id", middleware.AdminMiddleware(), userHandler.DeleteUser)

		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders", orderHandler.ListOrders)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.PUT("/orders/:id", middleware.AdminMiddleware(), orderHandler.UpdateOrder)
		api.DELETE("/orders/:id", middleware.AdminMiddleware(), orderHandler.DeleteOrder)

		// Admin-only stock routes
		api.POST("/stocks", middleware.AdminMiddleware(), stockHandler.CreateStock)
		api.PUT("/stocks/:id", middleware.AdminMiddleware(), stockHandler.UpdateStock)
		api.DELETE("/stocks/:id", middleware.AdminMiddleware(), stockHandler.DeleteStock)
	}

	fmt.Printf("Server starting on :%s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
