package webhandler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valkey-io/valkey-glide/go/api"
)

type MiddlewareService struct {
	valkey api.GlideClientCommands
}

func NewMiddlewareService(valkey api.GlideClientCommands) *MiddlewareService {
	return &MiddlewareService{valkey: valkey}
}

func (ms *MiddlewareService) CredentialAuth(next http.Handler) http.Handler {
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

		claims := &middleware.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check jwt blacklist
		blacklisted, err := ms.valkey.Get("jwt_blacklist:" + claims.Jti)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !blacklisted.IsNil() {
			utils.WriteError(w, http.StatusUnauthorized, errors.New("token is blacklisted"))
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
