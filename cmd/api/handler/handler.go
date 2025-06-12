package handler

import (
	"net/http"
	"os"

	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/computer"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/test"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/user"
	"github.com/valkey-io/valkey-glide/go/api"
	"gorm.io/gorm"
)

type API struct {
	db     *gorm.DB
	valkey api.GlideClientCommands
}

func NewAPI(db *gorm.DB, valkey api.GlideClientCommands) *API {
	return &API{
		db:     db,
		valkey: valkey,
	}
}

func (api *API) publicHandler(path string, handler http.HandlerFunc, mux *http.ServeMux) {
	mux.Handle(path, middleware.Logging(http.HandlerFunc(handler)))
}

func (api *API) privateHandler(path string, handler http.HandlerFunc, mux *http.ServeMux) {
	mux.Handle(path, middleware.ApiAuth(middleware.Logging(http.HandlerFunc(handler))))
}

func (api *API) privateCredentialHandler(path string, handler http.HandlerFunc, mux *http.ServeMux) {
	mux.Handle(path, middleware.CredentialAuth(middleware.Logging(http.HandlerFunc(handler))))
}

func (api *API) InitHandlers(mux *http.ServeMux) {
	computerService := computer.NewComputerService(api.db)
	userService := user.NewUserService(api.db, api.valkey)

	api.publicHandler("GET /api/v1/test", test.HelloRoute, mux)
	api.publicHandler("GET /api/v1/computer/{computerName}/get-rustdesk-id", computerService.GetComputerRustDeskIDRoute, mux)
	api.publicHandler("POST /api/v1/user/login", userService.LoginRoute, mux)
	if os.Getenv("ALLOW_REGISTER") == "true" {
		api.publicHandler("POST /api/v1/user/register", userService.RegisterRoute, mux)
	}

	api.privateHandler("POST /api/v1/computer/refresh", computerService.RefreshComputerRoute, mux)

}
