package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sinde530/go-mancer/cmd/token"
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

	// Set default image
	defaultImage := "assets/imgs/default-image.png"
	request.Image = defaultImage

	err := db.SaveUser(&request)
	if err != nil {
		if err.Error() == "Email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "이메일이 있습니다. 다른 이메일 또는 비밀번호 찾기를 해주세요."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		}
		return
	}

	user, err := db.GetUserByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	tokens, err := token.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// c.JSON(http.StatusCreated, user)
	// c.JSON(http.StatusCreated, gin.H{"user": user, "tokens": tokens})
	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})
}
