package userControllers

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// Generate a JWT token
func generateToken(userID uint) string {
	// Set the expiration time to 5 hours from now
	expirationTime := time.Now().Add(5 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(), // Set the expiration time as a UNIX timestamp
	})

	tokenString, _ := token.SignedString([]byte("cat_is_a_cat_but_cat_plus_cap_are_turtle")) // Use a secure secret for production
	return tokenString
}
