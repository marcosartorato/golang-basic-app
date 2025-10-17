package metrics

import (
	"context"
	"net/http"

	"github.com/marcosartorato/myapp/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
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
		cfg.Logger.Info("Metrics server listening on  " + addr)
		if err := http.ListenAndServe(addr, srv.Handler); err != nil {
			cfg.Logger.Error("metrics server failed: %v", zap.Error(err))
		}
	}()

	return srv.Shutdown
}
