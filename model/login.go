package model

type LoginRequest struct {
	Email    string `json:"email"`    // unique
	Password string `json:"password"` // Min 6, max 150 characters.
}
