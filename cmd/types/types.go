package types

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWT claims (payload tokenu)
type Claims struct {
	Email string `json:"email"`
	Token string `json:"token"`
	jwt.RegisteredClaims
}
