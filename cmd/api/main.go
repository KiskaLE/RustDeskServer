package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/handler"
	"github.com/KiskaLE/RustDeskServer/cmd/db"
	"github.com/KiskaLE/RustDeskServer/utils"
	"github.com/gorilla/mux"
)

// main initializes and starts the HTTP server. It loads environment variables,
// sets up the server mux with route handlers, and configures server timeouts.
// The server listens on the specified port, defaulting to 8080 if not set in
// the environment. It logs server start and handles any startup errors.

func main() {

	utils.LoadEnv()

	// migrate database
	db.MigrateDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := mux.NewRouter()

	handler.InitHandlers(mux)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	log.Println("Starting server on :" + port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
	s.ListenAndServe()
}
