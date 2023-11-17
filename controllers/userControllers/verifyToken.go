package userControllers

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// Verify JWT token
func verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("cat_is_a_cat_but_cat_plus_cap_are_turtle"), nil // Use the same secret used for generating the token
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["user_id"])
		fmt.Println("verify success")
		return claims, nil
	}

	return nil, err
}
