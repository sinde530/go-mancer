package model

type RegisterRequest struct {
	Email    string `json:"email"`    // unique
	Username string `json:"username"` // Min 7, max 30 characters.
	Password string `json:"password"` // Min 6, max 150 characters.
}
