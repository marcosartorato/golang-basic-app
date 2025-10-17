package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/marcosartorato/myapp/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns a /metrics handler exposing this registry.
func Handler() http.Handler {
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
}

// Start metric server
func CreateServer(cfg *config.ServerConfig) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", Handler())

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}
	return server
}

// Start run the HTTP server on dedicated goroutine and return the shutdown function.
func RunServerWithShutdown(cfg *config.ServerConfig) func(context.Context) error {
	srv := CreateServer(cfg)

	// Run metrics server in a goroutine
	go func() {
		addr := srv.Addr
		fmt.Println("Metrics server listening on  " + addr)
		if err := http.ListenAndServe(addr, srv.Handler); err != nil {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()

	return srv.Shutdown
}
