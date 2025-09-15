package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	helloSrv := hello.CreateServer(&cfg.HTTP)
	metricSrv := metrics.CreateServer(&cfg.Metrics)

	// Run main app server in a goroutine
	go func() {
		addr := helloSrv.Addr
		fmt.Println("App server listening on " + addr)
		if err := http.ListenAndServe(addr, helloSrv.Handler); err != nil {
			log.Fatalf("app server failed: %v", err)
		}
	}()
	// Run metrics server in a goroutine
	go func() {
		addr := metricSrv.Addr
		fmt.Println("Metrics server listening on  " + addr)
		if err := http.ListenAndServe(addr, metricSrv.Handler); err != nil {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()

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
	if err := helloSrv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	if err := metricSrv.Shutdown(ctx); err != nil {
		fmt.Printf("metrics server shutdown error: %v\n", err)
	}

	<-ctx.Done()
	fmt.Println("Servers stopped cleanly")

}
