package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"fmt"
	"strings"
)

var jwtSecret = []byte("3hRkY5F30Q9aE1/oRVFwT8N1DW3B/9+8OS1lDFGhfzs=") // Replace with a secure key

// GenerateToken generates a JWT token for a given user
func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	// Print the received token string
	fmt.Println("Received Token String:", tokenStr)

	// Remove "Bearer " prefix if it exists
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = tokenStr[len("Bearer "):]
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// Print the parsed token structure
	fmt.Printf("Parsed Token: %+v\n", token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Print the claims
		fmt.Printf("Claims: %+v\n", claims)
		return claims, nil
	}

	// Print the error if token is invalid
	fmt.Println("Token parsing error:", err)
	return nil, err
}