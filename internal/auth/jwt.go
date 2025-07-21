package auth

import "github.com/golang-jwt/jwt/v5"

var JwtKey = []byte("your_very_secret_key_change_in_production")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
