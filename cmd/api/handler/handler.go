package handler

import (
	"net/http"

	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/computer"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/test"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type API struct {
	db *gorm.DB
}

func NewAPI(db *gorm.DB) *API {
	return &API{
		db: db,
	}
}

func (api *API) publicHandler(path string, handler http.HandlerFunc, mux *mux.Router) {
	mux.Handle(path, middleware.Logging(http.HandlerFunc(handler))).Methods("GET", "POST")
}

func (api *API) privateHandler(path string, handler http.HandlerFunc, mux *mux.Router) {
	mux.Handle(path, middleware.ApiAuth(middleware.Logging(http.HandlerFunc(handler)))).Methods("GET", "POST", "PUT", "DELETE")
}

func (api *API) InitHandlers(mux *mux.Router) {
	computerService := computer.NewComputerService(api.db)

	api.publicHandler("/api/v1/test", test.HelloRoute, mux)
	api.publicHandler("/api/v1/computer/{computerName}/get-rustdesk-id", computerService.GetComputerRustDeskIDRoute, mux)

	api.privateHandler("/api/v1/computer/refresh", computerService.RefreshComputerRoute, mux)
}
