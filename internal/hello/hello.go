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
func Start(cfg *config.ServerConfig) {
	appMux := http.NewServeMux()

	appMux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metrics.RequestsTotal.WithLabelValues("/hello").Inc()
		HelloHandler(w, r)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("/hello").Observe(duration)
	})

	// Run main app server in a goroutine
	go func() {
		addr := cfg.Addr()
		fmt.Println("App server listening on " + addr)
		if err := http.ListenAndServe(addr, appMux); err != nil {
			log.Fatalf("app server failed: %v", err)
		}
	}()
}
