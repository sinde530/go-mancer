package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
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
	r.Static("/assets", "./assets")
	r.Static("/uploads", "./uploads")

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
	r.GET("/iptrack/:ip", IPTrack)
	r.POST("/create/group", handler.HandleCreateGroup)
	r.POST("/upload-url", handler.HandleUploadURL)

	port := os.Getenv("PORT")
	log.Fatal(r.Run(":" + port))
}

func HandleTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hi",
	})
}

func IPTrack(c *gin.Context) {
	ip := c.Param("ip")

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		c.JSON(400, gin.H{
			"error": "Invalid IP address",
		})
		return
	}

	// IP 위치 정보 조회
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s", url.PathEscape(ip))
	response, err := http.Get(apiURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// JSON 파싱
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	if result["status"] == "fail" {
		c.JSON(500, gin.H{
			"error": "Failed to retrieve location information",
		})
		return
	}

	c.JSON(200, gin.H{
		"country_code": result["countryCode"],
		"city":         result["city"],
		"latitude":     result["lat"],
		"longitude":    result["lon"],
		"isp":          result["isp"],
		"mobile":       result["mobile"],
		"as":           result["as"],
		"proxy":        result["proxy"],
		"hosting":      result["hosting"],
	})
}
