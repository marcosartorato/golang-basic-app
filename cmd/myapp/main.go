package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/marcosartorato/myapp/internal/config"
	httpSrv "github.com/marcosartorato/myapp/internal/http"
	"github.com/marcosartorato/myapp/internal/metrics"
)

func main() {
	// Get configuration
	cfg := config.Parse()

	// Create logger based on the log level.
	level, err := zapcore.ParseLevel(*cfg.LogLevel)
	if err != nil {
		log.Fatalf("invalid log level: %v", err)
	}
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	logger, err := zapConfig.Build()
	logger.Info("Logger initialized", zap.String("level", level.String()))
	if err != nil {
		log.Fatalf("invalid log level: %v", err)
	}
	defer func() { // flushes buffer, if any
		err := logger.Sync()
		if err != nil {
			log.Fatalf("error flushing logger buffer: %v", err)
		}
	}()

	// Save the logger in the server configs.
	cfg.HTTP.Logger = logger
	cfg.Metrics.Logger = logger

	// Start servers
	srvShutdown := httpSrv.RunServerWithShutdown(&cfg.HTTP)
	metricShutdown := metrics.RunServerWithShutdown(&cfg.Metrics)

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stop
	logger.Info("Shutting down gracefully...")

	// Give servers a few seconds to exit gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Shutdown on the server
	if err := srvShutdown(ctx); err != nil {
		logger.Sugar().Fatalf("server forced to shutdown: %v", err)
	}
	if err := metricShutdown(ctx); err != nil {
		logger.Sugar().Fatalf("metrics server shutdown error: %v", err)
	}

	<-ctx.Done()
	logger.Info("Servers stopped cleanly")
}
