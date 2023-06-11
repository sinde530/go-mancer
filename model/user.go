package model

type User struct {
	UID      string `json:"uid"`        // user ID
	Email    string `json:"email"`      // unique
	Username string `json:"username"`   // Min 7, max 30 characters.
	Password string `json:"password"`   // Min 6, max 150 characters.
	CreateAT string `json:"created_at"` // Date Created
	Image    string `json:"image"`      // Profile image url
}
