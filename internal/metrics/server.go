package metrics

import (
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
