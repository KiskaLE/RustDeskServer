package webmw

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/account"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valkey-io/valkey-glide/go/api"
)

type SessionMiddleware struct {
	valkey         api.GlideClientCommands
	accountService *account.AccountService
}

func New(valkey api.GlideClientCommands, accountService *account.AccountService) SessionMiddleware {
	return SessionMiddleware{valkey: valkey, accountService: accountService}
}

// VerifySession kontroluje, zda je v cookie platný JWT a není na blacklistu.
func (sm SessionMiddleware) VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("auth_token")
		if err != nil || c.Value == "" {
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
			refreshCookie, err := r.Cookie("refresh_token")
			if err != nil || refreshCookie.Value == "" {
				utils.WriteError(w, http.StatusUnauthorized, errors.New("refreshToken is missing: "+err.Error()))
				return
			}
			refreshedTokenStr, refreshToken, err := sm.handleRefreshToken(refreshCookie.Value)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to handle refresh token: "+err.Error()))
				return
			}

			// Set the new access token cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    refreshedTokenStr, // The new access token
				Path:     "/",
				HttpOnly: true,
				Secure:   r.TLS != nil, // Use secure cookies if on HTTPS
				SameSite: http.SameSiteLaxMode,
			})

			// Set the new refresh token cookie (Token Rotation)
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    refreshToken, // The new refresh token
				Path:     "/",
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteLaxMode,
			})

			http.Redirect(w, r, r.URL.Path, http.StatusFound)
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

func (sm SessionMiddleware) handleRefreshToken(refreshToken string) (token string, refresh string, err error) {
	// 2. Use your existing API logic to refresh the token
	// We'll simulate the API call just like you do in HandleLogin
	payload := account.RefreshTokenPayload{RefreshToken: refreshToken}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/api/v1/account/refresh-token", strings.NewReader(string(payloadBytes)))
	if err != nil {
		return "", "", err
	}
	rr := httptest.NewRecorder()

	// Call the core API refresh logic
	sm.accountService.RefreshTokenRoute(rr, req)

	// 3. Check the result and set new cookies
	if rr.Code != http.StatusOK {
		return "", "", errors.New("failed to refresh token")
	}

	// Success! Extract new tokens and set them as cookies.
	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	return response["token"], response["refresh_token"], nil
}
