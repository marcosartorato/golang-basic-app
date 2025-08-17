package hello

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/marcosartorato/myapp/internal/metrics"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

// Start hello-world server
func Start() {
	// --- Application server (port 8080) ---
	appMux := http.NewServeMux()
	appMux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metrics.RequestsTotal.WithLabelValues("/hello").Inc()
		HelloHandler(w, r)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("/hello").Observe(duration)
	})

	// Run main app server
	go func() {
		fmt.Println("App server listening on :8080")
		if err := http.ListenAndServe(":8080", appMux); err != nil {
			log.Fatalf("app server failed: %v", err)
		}
	}()
}
