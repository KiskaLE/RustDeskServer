package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// ParseJSON reads the request body as JSON and decodes it into the provided value.
// It returns an error if the request body is missing or if the JSON decoding fails.
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJSON writes a JSON response with the provided HTTP status code.
// It accepts an http.ResponseWriter to write the response, an integer status code, and a value of any type to be encoded as JSON.
// The function sets the "Content-Type" header to "application/json" and returns an error if JSON encoding fails.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// WriteError writes an error message as a JSON response with the provided HTTP status code.
// It accepts an http.ResponseWriter to write the response, an integer status code, and an error object.
// The function returns an error if the JSON encoding fails.
func WriteError(w http.ResponseWriter, status int, err error) error {
	return WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func LoadEnv() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
