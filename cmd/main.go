package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sinde530/go-mancer/cmd/handler"
	"github.com/sinde530/go-mancer/cmd/token"
)

func main() {
	log.Println("Starting server...")

	r := gin.Default()

	config := cors.DefaultConfig()                                                              // CORS Settings
	config.AllowAllOrigins = true                                                               // All Domain Allwo Path
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"} // 허용할 헤더를 지정
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
	r.POST("/account/register", handler.HandleRegister)
	r.POST("/account/login", handler.HandleLogin)
	r.POST("/account/logout", handler.HandleLogout)
	r.POST("/refresh", token.TokenChange)

	port := os.Getenv("PORT")
	log.Fatal(r.Run(":" + port))
}

func HandleTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hi",
	})
}
