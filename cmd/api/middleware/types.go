package middleware

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Jti   string
	ID    string `json:"Email"`
	Token string `json:"Token"`
	jwt.RegisteredClaims
}
