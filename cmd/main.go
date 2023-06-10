package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server...")

	r := gin.Default()

	config := cors.DefaultConfig()      // CORS Settings
	config.AllowOrigins = []string{"*"} // All Domain Allwo Path
	r.Use(cors.New(config))

	// Load env from .env file
	if gin.Mode() != gin.ReleaseMode {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}

		log.Printf("Successd env load")
	}

	r.GET("/", HandleTest)

	port := os.Getenv("PORT")
	log.Fatal(r.Run(":" + port))
}

func HandleTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hi",
	})
}
