package main

import (

	// "net/http"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No env file found. using default port")
	}

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	
}
