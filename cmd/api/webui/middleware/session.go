package webmw

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valkey-io/valkey-glide/go/api"
)

type SessionMiddleware struct {
	valkey api.GlideClientCommands
}

func New(valkey api.GlideClientCommands) SessionMiddleware {
	return SessionMiddleware{valkey: valkey}
}

// VerifySession kontroluje, zda je v cookie platný JWT a není na blacklistu.
func (sm SessionMiddleware) VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHtmx := r.Header.Get("HX-Request") == "true"
		c, err := r.Cookie("auth_token")
		if err != nil || c.Value == "" {
			if isHtmx {
				w.WriteHeader(http.StatusUnauthorized) // 401
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenStr := c.Value
		jwtKey := []byte(os.Getenv("SALT"))

		claims := &middleware.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if errors.Is(err, jwt.ErrTokenExpired) {
			utils.WriteError(w, http.StatusUnauthorized, errors.New("token_expired"))
			return
		}
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// kontrola blacklistu (logout)
		bl, err := sm.valkey.Get("jwt_blacklist:" + claims.Jti)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !bl.IsNil() {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// vložíme claims do contextu, pokud by je handler potřeboval
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sm SessionMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
