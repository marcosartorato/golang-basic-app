package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marcosartorato/myapp/internal/config"
	httpSrv "github.com/marcosartorato/myapp/internal/http"
	"github.com/marcosartorato/myapp/internal/metrics"
)

func main() {
	// Get configuration
	cfg := config.Parse()

	// Start servers
	srvShutdown := httpSrv.RunServerWithShutdown(&cfg.HTTP)
	metricShutdown := metrics.RunServerWithShutdown(&cfg.Metrics)

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stop
	fmt.Println("\nShutting down gracefully...")

	// Give servers a few seconds to exit gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Shutdown on the server
	if err := srvShutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	if err := metricShutdown(ctx); err != nil {
		fmt.Printf("metrics server shutdown error: %v\n", err)
	}

	<-ctx.Done()
	fmt.Println("Servers stopped cleanly")

}
