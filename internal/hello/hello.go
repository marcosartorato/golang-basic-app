package hello

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcosartorato/myapp/internal/config"
	"github.com/marcosartorato/myapp/internal/metrics"
)

// HelloHandler is the handler for the "Hello, World!" path.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintln(w, "Hello, World!"); err != nil {
		// Log the error; since client likely went away, not much else to do
		log.Printf("failed to write response: %v", err)
	}
}

// Start run the hello-world server on dedicated goroutine.
func CreateServer(cfg *config.ServerConfig) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/hello", metrics.Instrument(http.HandlerFunc(HelloHandler)))
	mux.Handle("/api/message", metrics.Instrument(http.HandlerFunc(MessageHandler)))

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}
	return server
}
