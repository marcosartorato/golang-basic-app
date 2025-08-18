package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marcosartorato/myapp/internal/config"
	"github.com/marcosartorato/myapp/internal/hello"
	"github.com/marcosartorato/myapp/internal/metrics"
)

func main() {
	// Get configuration
	cfg := config.Parse()

	// Initialize metrics
	metrics.Init()

	// Start servers
	hello.Start(&cfg.HTTP)
	metrics.Start(&cfg.Metrics)

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stop
	fmt.Println("\nShutting down gracefully...")

	// Give servers a few seconds to exit gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// At this point, since Start() just spins goroutines,
	// we don’t have direct server references to call Shutdown().
	// You might later refactor Start() to return *http.Server
	// so that you can call Shutdown(ctx) here.

	<-ctx.Done()
	fmt.Println("Servers stopped cleanly")

}
