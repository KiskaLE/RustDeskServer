package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware checks if the user is authenticated
func ApiAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a valid token (this is just an example, implement your own logic)
		token := os.Getenv("AUTH_TOKEN")
		user_token := r.Header.Get("api_key")

		if token != user_token {
			utils.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CredentialAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO implement
		tokenString := r.Header.Get("Authorization")
		// remove the "Bearer " prefix from the token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			utils.WriteError(w, http.StatusUnauthorized, errors.New("token is missing"))
			return
		}

		jwtKey := []byte(os.Getenv("SALT"))

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
