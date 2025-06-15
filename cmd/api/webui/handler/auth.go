// cmd/api/web/handler/auth.go
package webhandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/account"
	"github.com/KiskaLE/RustDeskServer/cmd/api/webui/view/pages"
)

type AuthHandler struct {
	accountService *account.AccountService
}

func NewAuthHandler(as *account.AccountService) *AuthHandler {
	return &AuthHandler{
		accountService: as,
	}
}

// Zobrazí přihlašovací stránku
func (h *AuthHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	// Zatím bez chybové zprávy
	pages.LoginPage("").Render(r.Context(), w)
}

// Zpracuje data z přihlašovacího formuláře
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		pages.LoginPage("Nepodařilo se zpracovat formulář.").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// --- Elegantní způsob, jak volat existující API logiku ---
	// Vytvoříme si falešný payload a request pro naši existující API metodu.
	payload := account.LoginPayload{Email: email, Password: password}
	payloadBytes, _ := json.Marshal(payload)

	// `httptest` nám pomůže "nasimulovat" API volání interně
	req, err := http.NewRequest("POST", "/api/v1/account/login", strings.NewReader(string(payloadBytes)))
	if err != nil {
		pages.LoginPage("Nepodařilo se zpracovat formulář.").Render(r.Context(), w)
		return
	}
	rr := httptest.NewRecorder()

	// Zavoláme existující handler
	h.accountService.LoginRoute(rr, req)

	// Zkontrolujeme výsledek
	if rr.Code != http.StatusOK {
		// Přihlášení selhalo, znovu zobrazíme login formulář s chybou
		pages.LoginPage("Nesprávný email nebo heslo.").Render(r.Context(), w)
		return
	}

	// Přihlášení bylo úspěšné!
	// V reálu bychom zde nastavili cookie s JWT tokenem.
	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Nastavíme httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    response["token"],
		Path:     "/",
		HttpOnly: true, // Klíčové pro bezpečnost!
		Secure:   true, // V produkci vždy true
		SameSite: http.SameSiteLaxMode,
		// MaxAge: 3600, // 1 hodina
	})

	// Nastavíme httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    response["refresh_token"],
		Path:     "/",
		HttpOnly: true, // Klíčové pro bezpečnost!
		Secure:   true, // V produkci vždy true
		SameSite: http.SameSiteLaxMode,
		// MaxAge: 3600, // 1 hodina
	})

	// A přesměrujeme uživatele na dashboard pomocí HTMX hlavičky
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// It reads the refresh_token from an HttpOnly cookie.
func (h *AuthHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	// 1. Get the refresh token from the cookie
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		// If cookie is not found, the session is invalid.
		http.Error(w, "Unauthorized: No refresh token", http.StatusUnauthorized)
		return
	}

	// 2. Use your existing API logic to refresh the token
	// We'll simulate the API call just like you do in HandleLogin
	payload := account.RefreshTokenPayload{RefreshToken: refreshCookie.Value}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/account/refresh-token", strings.NewReader(string(payloadBytes)))
	rr := httptest.NewRecorder()

	// Call the core API refresh logic
	h.accountService.RefreshTokenRoute(rr, req)

	// 3. Check the result and set new cookies
	if rr.Code != http.StatusOK {
		// The refresh token was invalid or expired. Force a full logout/login.
		// We clear the cookies to be safe.
		http.SetCookie(w, &http.Cookie{Name: "authtoken", Path: "/", MaxAge: -1})
		http.SetCookie(w, &http.Cookie{Name: "refresh_token", Path: "/", MaxAge: -1})
		http.Error(w, "Unauthorized: Refresh failed", http.StatusUnauthorized)
		return
	}

	// 4. Success! Extract new tokens and set them as cookies.
	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Set the new access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "authtoken",
		Value:    response["token"], // The new access token
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Use secure cookies if on HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Set the new refresh token cookie (Token Rotation)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    response["refresh_token"], // The new refresh token
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	// Send a 200 OK response to the JS interceptor.
	// The body can be empty.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token refreshed"))
}
