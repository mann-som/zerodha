package main

import (

	// "net/http"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mann-som/zerodha/internal/handlers"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No env file found. using default port")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "Tradingin API is up\n")
	})

	router.GET("/user", handlers.UserHandler)

	fmt.Printf("Server starting on :%s Port \n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
