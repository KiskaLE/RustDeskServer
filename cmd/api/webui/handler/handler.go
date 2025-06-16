package webhandler

import (
	"net/http"

	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/account"
	webmw "github.com/KiskaLE/RustDeskServer/cmd/api/webui/middleware"
	"github.com/valkey-io/valkey-glide/go/api"
	"gorm.io/gorm"
)

type UI struct {
	db     *gorm.DB
	valkey api.GlideClientCommands
}

func NewUI(db *gorm.DB, valkey api.GlideClientCommands) *UI {
	return &UI{db: db, valkey: valkey}
}

func (ui *UI) InitHandlers(mux *http.ServeMux) {
	accountService := account.NewAccountService(ui.db, ui.valkey)
	authWebHandler := NewAuthHandler(accountService)

	mux.HandleFunc("GET /login", authWebHandler.HandleLoginPage)
	mux.HandleFunc("POST /login", authWebHandler.HandleLogin)

	mux.HandleFunc("POST /web/refresh-token", authWebHandler.HandleRefreshToken)

	sessionMW := webmw.New(ui.valkey)
	mux.Handle("GET /dashboard",
		sessionMW.VerifySession(http.HandlerFunc(HandleDashboardPage)))
	mux.Handle("GET /",
		sessionMW.VerifySession(http.HandlerFunc(HandleRoot)))
}
