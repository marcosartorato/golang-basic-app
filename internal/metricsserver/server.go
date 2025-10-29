package metricsserver

import (
	"context"
	"net"
	"net/http"

	cfg "github.com/marcosartorato/myapp/internal/config"
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
func CreateServer(opt cfg.Options) *http.Server {
	mux := http.NewServeMux()

	mux.Handle(
		"/metrics",
		http.TimeoutHandler(
			Handler(),
			opt.TimeoutHandler,
			"Service Timeout",
		),
	)

	addr := net.JoinHostPort(*opt.Host, *opt.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       opt.ReadTimeout,
		ReadHeaderTimeout: opt.ReadHeaderTimeout,
		IdleTimeout:       opt.IdleTimeout,
	}
	return server
}

// Start run the HTTP server on dedicated goroutine and return the shutdown function.
func RunServerWithShutdown(logger *zap.Logger, opts ...cfg.Option) func(context.Context) error {
	var options cfg.Options
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			panic(err)
		}
	}

	srv := CreateServer(options)

	// Run metrics server in a goroutine
	go func() {
		addr := srv.Addr
		logger.Info("Metrics server listening on  " + addr)
		if err := http.ListenAndServe(addr, srv.Handler); err != nil {
			logger.Error("metrics server failed: %v", zap.Error(err))
		}
	}()

	return srv.Shutdown
}
