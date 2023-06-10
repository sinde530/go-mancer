package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	UID      string `json:"uid"`        // user ID
	Email    string `json:"email"`      // unique
	Username string `json:"username"`   // Min 7, max 30 characters.
	Password string `json:"password"`   // Min 6, max 150 characters.
	CreateAT string `json:"created_at"` // Date Created
}

func HandleRegister(c *gin.Context) {
	var request registerRequest

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

	c.JSON(http.StatusOK, gin.H{"message": "Register Successful"})
}
