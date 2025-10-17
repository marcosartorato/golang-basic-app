package hello

import (
	"fmt"
	"log"
	"net/http"
)

// HelloHandler is the handler for the "Hello, World!" path.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintln(w, "Hello, World!"); err != nil {
		// Log the error; since client likely went away, not much else to do
		log.Printf("failed to write response: %v", err)
	}
}
