package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	cfg "github.com/marcosartorato/myapp/internal/config"
	httpSrv "github.com/marcosartorato/myapp/internal/httpserver"
	metricsSrv "github.com/marcosartorato/myapp/internal/metricsserver"
)

const (
	envHTTPPort                 = "HTTP_PORT"
	envHTTPHost                 = "HTTP_HOST"
	envHTTPReadHeaderTimeout    = "HTTP_READ_HEADER_TIMEOUT"
	envHTTPReadTimeout          = "HTTP_READ_TIMEOUT"
	envHTTPTimeoutHandler       = "HTTP_TIMEOUT_HANDLER"
	envHTTPIdleTimeout          = "HTTP_IDLE_TIMEOUT"
	envMetricsPort              = "METRICS_PORT"
	envMetricsHost              = "METRICS_HOST"
	envMetricsReadHeaderTimeout = "METRICS_READ_HEADER_TIMEOUT"
	envMetricsReadTimeout       = "METRICS_READ_TIMEOUT"
	envMetricsTimeoutHandler    = "METRICS_TIMEOUT_HANDLER"
	envMetricsIdleTimeout       = "METRICS_IDLE_TIMEOUT"
	envLogLevel                 = "LOG_LEVEL"
	defaultTimeoutMs            = 1000
)

// envOrDefaultStr returns env value if set, otherwise fallback to default.
func envOrDefaultStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// envOrDefaultInt64 like envOrDefaultString but also check the env value is a valid int64.
func envOrDefaultInt64(key string, def int64) int64 {
	if vStr := os.Getenv(key); vStr != "" {
		vInt, err := strconv.ParseInt(vStr, 10, 64)
		if err == nil {
			return vInt
		}
		// else fall through to return def
	}
	return def
}

func main() {
	// Get configuration
	httpHost := flag.String("http-host", envOrDefaultStr(envHTTPHost, "localhost"), "host for the HTTP server (also via "+envHTTPHost+")")
	httpPort := flag.String("http-port", envOrDefaultStr(envHTTPPort, "8080"), "port for the HTTP server (also via "+envHTTPPort+")")
	httpReadHeaderTimeout := flag.Int64("http-read-header-timeout", envOrDefaultInt64(envHTTPReadHeaderTimeout, defaultTimeoutMs), "max amount of time to read the request headers (also via "+envHTTPReadHeaderTimeout+")")
	httpReadTimeout := flag.Int64("http-read-timeout", envOrDefaultInt64(envHTTPReadTimeout, defaultTimeoutMs), "max amount of time to read the entire request (also via "+envHTTPReadTimeout+")")
	httpTimeoutHandler := flag.Int64("http-timeout-handler", envOrDefaultInt64(envHTTPTimeoutHandler, defaultTimeoutMs), "max amount of time for a handler to complete (also via "+envHTTPTimeoutHandler+")")
	httpIdleTimeout := flag.Int64("http-idle-timeout", envOrDefaultInt64(envHTTPIdleTimeout, defaultTimeoutMs), "max amount of time to wait for the next request when keep-alives are enabled (also via "+envHTTPIdleTimeout+")")

	metricsHost := flag.String("metrics-host", envOrDefaultStr(envMetricsHost, "localhost"), "host for the metrics server (also via "+envMetricsHost+")")
	metricsPort := flag.String("metrics-port", envOrDefaultStr(envMetricsPort, "9090"), "port for the metrics server (also via "+envMetricsPort+")")
	metricsReadHeaderTimeout := flag.Int64("metrics-read-header-timeout", envOrDefaultInt64(envMetricsReadHeaderTimeout, defaultTimeoutMs), "max amount of time to read the request headers (also via "+envHTTPReadHeaderTimeout+")")
	metricsReadTimeout := flag.Int64("metrics-read-timeout", envOrDefaultInt64(envMetricsReadTimeout, defaultTimeoutMs), "max amount of time to read the entire request (also via "+envHTTPReadTimeout+")")
	metricsTimeoutHandler := flag.Int64("metrics-timeout-handler", envOrDefaultInt64(envMetricsTimeoutHandler, defaultTimeoutMs), "max amount of time for a handler to complete (also via "+envHTTPTimeoutHandler+")")
	metricsIdleTimeout := flag.Int64("metrics-idle-timeout", envOrDefaultInt64(envMetricsIdleTimeout, defaultTimeoutMs), "max amount of time to wait for the next request when keep-alives are enabled (also via "+envMetricsIdleTimeout+")")

	logLevel := flag.String("log-level", envOrDefaultStr(envLogLevel, "info"), "logging level (debug, info, warn, error)")

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
		cfg.WithReadHeaderTimeout(*httpReadHeaderTimeout),
		cfg.WithReadTimeout(*httpReadTimeout),
		cfg.WithTimeoutHandler(*httpTimeoutHandler),
		cfg.WithIdleTimeout(*httpIdleTimeout),
	)
	metricShutdown := metricsSrv.RunServerWithShutdown(
		logger,
		cfg.WithHost(metricsHost),
		cfg.WithPort(metricsPort),
		cfg.WithReadHeaderTimeout(*metricsReadHeaderTimeout),
		cfg.WithReadTimeout(*metricsReadTimeout),
		cfg.WithTimeoutHandler(*metricsTimeoutHandler),
		cfg.WithIdleTimeout(*metricsIdleTimeout),
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
