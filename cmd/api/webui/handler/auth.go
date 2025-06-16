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
