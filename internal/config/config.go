package config

import (
	"flag"
	"fmt"

	"go.uber.org/zap"
)

const (
	envHTTPPort    = "HTTP_PORT"
	envHTTPHost    = "HTTP_HOST"
	envMetricsPort = "METRICS_PORT"
	envMetricsHost = "METRICS_HOST"
	envLogLevel    = "LOG_LEVEL"
)

// ServerConfig holds host and port information for a server.
type ServerConfig struct {
	Host   string
	Port   string
	Logger *zap.Logger
}

// Config holds application configuration.
type Config struct {
	HTTP     ServerConfig
	Metrics  ServerConfig
	LogLevel *string
}

// Parse builds a Config struct from flags or environment variables.
func Parse() *Config {

	httpHost := flag.String("http-host", envOrDefault(envHTTPHost, "localhost"), "host for the HTTP server (also via "+envHTTPHost+")")
	httpPort := flag.String("http-port", envOrDefault(envHTTPPort, "8080"), "port for the HTTP server (also via "+envHTTPPort+")")
	metricsHost := flag.String("metrics-host", envOrDefault(envMetricsHost, "localhost"), "host for the metrics server (also via "+envMetricsHost+")")
	metricsPort := flag.String("metrics-port", envOrDefault(envMetricsPort, "9090"), "port for the metrics server (also via "+envMetricsPort+")")

	// log level has to be validated after the parsing
	logLevel := flag.String("log-level", envOrDefault(envLogLevel, "info"), "logging level (debug, info, warn, error)") // envOrDefault(envLogLevel, "info")

	// Parse flags
	flag.Parse()

	// Validate log level
	allowed := map[string]struct{}{
		"debug": {},
		"info":  {},
		"warn":  {},
		"error": {},
	}
	if _, ok := allowed[*logLevel]; !ok {
		panic(fmt.Sprintf("invalid log level: %s", *logLevel))
	}

	return &Config{
		HTTP: ServerConfig{
			Host: *httpHost,
			Port: *httpPort,
		},
		Metrics: ServerConfig{
			Host: *metricsHost,
			Port: *metricsPort,
		},
		LogLevel: logLevel,
	}
}

// Addr returns the server configuration information using the format "${HOST}:${PORT}" .
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
