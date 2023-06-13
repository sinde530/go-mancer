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
	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Failed to parse request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request"})
		return
	}

	fmt.Println("Email: ", user.Email)
	fmt.Println("Username: ", user.Username)
	fmt.Println("Password: ", user.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hashedPassword)
	fmt.Println("hashedPassword: ", user.Password)

	// Create uid
	user.UID = fmt.Sprintf("%d", time.Now().Unix())
	fmt.Println("Password: ", user.UID)

	// Create Date
	t := time.Now()
	user.CreateAT = t.Format("2006-01-02 15:04:05")
	fmt.Println("CreateAT: ", user.CreateAT)

	// Set default image
	defaultImage := "http://localhost:8080/assets/imgs/default-image.png"
	user.Image = defaultImage

	// err := db.SaveUser(&user)
	err = db.SaveUser(&user)
	if err != nil {
		if err.Error() == "Email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "I have an email. Please enter another email"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Failed to save user"})
		}
		return
	}

	retrievedUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	token, err := token.GenerateTokens(retrievedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tokens": token})
}

func HandleLogin(c *gin.Context) {
	var request model.LoginRequest

	// Parsing the JSON body (email and password)
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Failed to parse request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request"})
		return
	}

	// Authenticating the user
	user, err := db.AuthenticateUser(request.Email, request.Password)
	if err != nil {
		log.Printf("Failed to authenticate user: %v", err)
		if err.Error() == "User not found" {
			c.JSON(http.StatusConflict, gin.H{"error": "Incorrect email or password"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		}
		return
	}

	// Generating the tokens
	tokens, err := token.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func HandleLogout(c *gin.Context) {
	token, claims, err := token.VerifyToken(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": claims.User, "message": "Logged out successfully"})
}
