package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sinde530/go-mancer/db"
	"github.com/sinde530/go-mancer/model"
)

func HandleRegister(c *gin.Context) {
	var request model.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Failed to parse request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request"})
		return
	}

	fmt.Println("Email: ", request.Email)
	fmt.Println("Username: ", request.Username)
	fmt.Println("Password: ", request.Password)

	// Create uid
	request.UID = fmt.Sprintf("%d", time.Now().Unix())
	fmt.Println("Password: ", request.UID)

	// Create Date
	t := time.Now()
	request.CreateAT = t.Format("2006-01-02 15:04:05")
	fmt.Println("Password: ", request.CreateAT)

	err := db.SaveUser(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register Successful"})
}
