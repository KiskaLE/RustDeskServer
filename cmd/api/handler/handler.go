package handler

import (
	"net/http"

	"github.com/KiskaLE/RustDeskServer/cmd/api/middleware"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/computer"
	"github.com/KiskaLE/RustDeskServer/cmd/api/routes/test"
	"github.com/gorilla/mux"
)

func publicHandler(path string, handler http.HandlerFunc, mux *mux.Router) {
	mux.Handle(path, middleware.Logging(http.HandlerFunc(handler))).Methods("GET", "POST")
}

func privateHandler(path string, handler http.HandlerFunc, mux *mux.Router) {
	mux.Handle(path, middleware.ApiAuth(middleware.Logging(http.HandlerFunc(handler)))).Methods("GET", "POST", "PUT", "DELETE")
}

func InitHandlers(mux *mux.Router) {
	publicHandler("/api/v1/test", test.HelloRoute, mux)

	privateHandler("/api/v1/computer/refresh", computer.RefreshComputerRoute, mux)
}
