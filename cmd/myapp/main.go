package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	cfg "github.com/marcosartorato/myapp/internal/config"
	httpSrv "github.com/marcosartorato/myapp/internal/httpserver"
	metricsSrv "github.com/marcosartorato/myapp/internal/metricsserver"
)

const (
	envHTTPPort    = "HTTP_PORT"
	envHTTPHost    = "HTTP_HOST"
	envMetricsPort = "METRICS_PORT"
	envMetricsHost = "METRICS_HOST"
	envLogLevel    = "LOG_LEVEL"
)

// envOrDefault returns env value if set, otherwise fallback.
func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	// Get configuration
	httpHost := flag.String("http-host", envOrDefault(envHTTPHost, "localhost"), "host for the HTTP server (also via "+envHTTPHost+")")
	httpPort := flag.String("http-port", envOrDefault(envHTTPPort, "8080"), "port for the HTTP server (also via "+envHTTPPort+")")
	metricsHost := flag.String("metrics-host", envOrDefault(envMetricsHost, "localhost"), "host for the metrics server (also via "+envMetricsHost+")")
	metricsPort := flag.String("metrics-port", envOrDefault(envMetricsPort, "9090"), "port for the metrics server (also via "+envMetricsPort+")")
	logLevel := flag.String("log-level", envOrDefault(envLogLevel, "info"), "logging level (debug, info, warn, error)")
	flag.Parse()

	// Create logger based on the log level.
	level, err := zapcore.ParseLevel(*logLevel)
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

	// Start servers
	srvShutdown := httpSrv.RunServerWithShutdown(
		logger,
		cfg.WithHost(httpHost),
		cfg.WithPort(httpPort),
	)
	metricShutdown := metricsSrv.RunServerWithShutdown(
		logger,
		cfg.WithHost(metricsHost),
		cfg.WithPort(metricsPort),
	)

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
