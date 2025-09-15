package hello

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metrics.RequestsTotal.WithLabelValues("/hello").Inc()
		HelloHandler(w, r)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("/hello").Observe(duration)
	})
	mux.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metrics.RequestsTotal.WithLabelValues("/api/message").Inc()
		MessageHandler(w, r)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("/api/message").Observe(duration)
	})

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}
	return server
}
