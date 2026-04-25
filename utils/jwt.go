package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT signing uses HMAC (hashing), which works on raw bytes, not strings.
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Unix() → convert to seconds format
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (int, error) {
	// jwt.Parse does Call the function below to get secret key
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid Claims")
	}

	// JSON numbers are decoded as float64 by default so even if you stored int → you get float64 here
	userID := int(claims["user_id"].(float64))

	return userID, nil
}
