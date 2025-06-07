package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/database"
	"github.com/KiskaLE/RustDeskServer/cmd/api/handler"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// main initializes and starts the HTTP server. It loads environment variables,
// sets up the server mux with route handlers, and configures server timeouts.
// The server listens on the specified port, defaulting to 8080 if not set in
// the environment. It logs server start and handles any startup errors.

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// migrate database
	database.MigrateDatabase(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := mux.NewRouter()

	ctx := handler.NewAPI(db)

	ctx.InitHandlers(mux)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	isHttps := os.Getenv("HTTPS")
	if isHttps == "true" {
		log.Println("Starting https server on :" + port)
		if err := s.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed")
		}
	} else {
		log.Println("Starting http server on :" + port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed")
		}
	}
}
