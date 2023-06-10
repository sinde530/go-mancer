package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server...")

	// Load env from .env file
	if gin.Mode() != gin.ReleaseMode {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
}
