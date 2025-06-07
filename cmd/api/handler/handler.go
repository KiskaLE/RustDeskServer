package handler

import (
	"net/http"

	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/computer"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/test"
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

func (api *API) publicHandler(path string, handler http.HandlerFunc, mux *http.ServeMux) {
	mux.Handle(path, middleware.Logging(http.HandlerFunc(handler)))
}

func (api *API) privateHandler(path string, handler http.HandlerFunc, mux *http.ServeMux) {
	mux.Handle(path, middleware.ApiAuth(middleware.Logging(http.HandlerFunc(handler))))
}

func (api *API) InitHandlers(mux *http.ServeMux) {
	computerService := computer.NewComputerService(api.db)

	api.publicHandler("GET /api/v1/test", test.HelloRoute, mux)
	api.publicHandler("GET /api/v1/computer/{computerName}/get-rustdesk-id", computerService.GetComputerRustDeskIDRoute, mux)

	api.privateHandler("POST /api/v1/computer/refresh", computerService.RefreshComputerRoute, mux)
}
