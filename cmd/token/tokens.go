package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sinde530/go-mancer/model"
)

// your secret key for signing the JWT, consider read it from env var for security
var jwtKey = []byte("your-secret-key")

// Email string `json:"email"`
type Claims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

func GenerateTokens(user *model.User) (*model.Tokens, error) {
	// Create the JWT claims, which includes the user email and expiry time
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	ss, err := token.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	// do the same for refreshToken
	// refreshClaims := &Claims{
	// 	User: user,
	// 	StandardClaims: jwt.StandardClaims{
	// 		// refresh token is typically longer lived
	// 		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	// 	},
	// }

	// refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// ssRefresh, err := refreshToken.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	return &model.Tokens{
		AccessToken: ss,
		// RefreshToken: ssRefresh,
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
