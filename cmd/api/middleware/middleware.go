package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KiskaLE/RustDeskServer/utils"
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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
