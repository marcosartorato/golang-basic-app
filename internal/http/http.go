package hello

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/marcosartorato/myapp/internal/config"
	"github.com/marcosartorato/myapp/internal/metrics"
)

// Start run the HTTP on dedicated goroutine.
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

// Start run the HTTP server on dedicated goroutine and return the shutdown function.
func RunServerWithShutdown(cfg *config.ServerConfig) func(context.Context) error {
	srv := CreateServer(cfg)

	go func() {
		addr := srv.Addr
		fmt.Println("App server listening on " + addr)
		if err := http.ListenAndServe(addr, srv.Handler); err != nil {
			log.Fatalf("app server failed: %v", err)
		}
	}()

	return srv.Shutdown
}
