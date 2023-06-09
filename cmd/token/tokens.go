package token

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sinde530/go-mancer/model"
)

var jwtKey = []byte("your-secret-key")

type Claims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

func GenerateTokens(user *model.User) (*model.Tokens, error) {
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	refreshClaims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			// refresh token is typically longer lived
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	ssRefresh, err := refreshToken.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  ss,
		RefreshToken: ssRefresh,
	}, nil
}

func VerifyToken(c *gin.Context) (*jwt.Token, *Claims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, nil, fmt.Errorf("authorization header is missing")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		return nil, nil, fmt.Errorf("invalid token format")
	}

	tokenStr := bearerToken[1]
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}

func TokenChange(c *gin.Context) {
	// We need to parse the Refresh Token from Authorization Header
	_, claims, err := VerifyToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	newTokens, err := GenerateTokens(claims.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": newTokens})
}
