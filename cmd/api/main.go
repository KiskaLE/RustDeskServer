package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KiskaLE/RustDeskServer/cmd/api/database"
	"github.com/KiskaLE/RustDeskServer/cmd/api/handler"
	webhandler "github.com/KiskaLE/RustDeskServer/cmd/api/webui/handler"
	"github.com/joho/godotenv"
	"github.com/valkey-io/valkey-glide/go/api"
	valkeyAPI "github.com/valkey-io/valkey-glide/go/api"
)

// main initializes and starts the HTTP server. It loads environment variables,
// sets up the server mux with route handlers, and configures server timeouts.
// The server listens on the specified port, defaulting to 8080 if not set in
// the environment. It logs server start and handles any startup errors.

var runningInDocker string = "false"

func main() {

	if runningInDocker == "false" {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	valkeyHost := os.Getenv("VALKEY_HOST")
	valkeyPort, err := strconv.Atoi(os.Getenv("VALKEY_PORT"))
	if err != nil {
		valkeyPort = 6379
	}

	// init valkey
	valkeyConfig := valkeyAPI.NewGlideClientConfiguration().WithAddress(&valkeyAPI.NodeAddress{Host: valkeyHost, Port: valkeyPort})

	valkey, err := api.NewGlideClient(valkeyConfig)
	if err != nil {
		log.Fatalf("Error creating valkey client: %v", err)
	}
	defer valkey.Close()

	// Test the connection
	result, err := valkey.Ping()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(result)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// migrate database
	database.MigrateDatabase(db)

	port := "3000"

	mux := http.NewServeMux()

	apiHandler := handler.NewAPI(db, valkey)

	apiHandler.InitHandlers(mux)

	webHandler := webhandler.NewUI(db, valkey)
	webHandler.InitHandlers(mux)

	isHttps := os.Getenv("HTTPS")
	if isHttps == "true" {
		certFilePath := os.Getenv("CERT_FILE_PATH")
		keyFilePath := os.Getenv("KEY_FILE_PATH")

		serverTLSCert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
		if err != nil {
			log.Fatal(err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverTLSCert},
		}

		s := &http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			IdleTimeout:  5 * time.Second,
			TLSConfig:    tlsConfig,
		}

		log.Println("Starting https server on :" + port)
		if err := s.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed")
		}
	} else {
		s := &http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			IdleTimeout:  5 * time.Second,
		}

		log.Println("Starting http server on :" + port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed")
		}
	}
}
