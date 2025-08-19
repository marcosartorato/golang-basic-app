package config

import (
	"flag"
	"fmt"
)

const (
	envHTTPPort    = "HTTP_PORT"
	envHTTPHost    = "HTTP_HOST"
	envMetricsPort = "METRICS_PORT"
	envMetricsHost = "METRICS_HOST"
)

// ServerConfig holds host and port information for a server.
type ServerConfig struct {
	Host string
	Port string
}

// Config holds application configuration.
type Config struct {
	HTTP    ServerConfig
	Metrics ServerConfig
}

// Parse builds a Config struct from flags or environment variables.
func Parse() *Config {

	httpHost := flag.String("http-host", envOrDefault(envHTTPHost, "localhost"), "host for the HTTP server (also via "+envHTTPHost+")")
	httpPort := flag.String("http-port", envOrDefault(envHTTPPort, "8080"), "port for the HTTP server (also via "+envHTTPPort+")")
	metricsHost := flag.String("metrics-host", envOrDefault(envMetricsHost, "localhost"), "host for the metrics server (also via "+envMetricsHost+")")
	metricsPort := flag.String("metrics-port", envOrDefault(envMetricsPort, "9090"), "port for the metrics server (also via "+envMetricsPort+")")

	flag.Parse()

	return &Config{
		HTTP: ServerConfig{
			Host: *httpHost,
			Port: *httpPort,
		},
		Metrics: ServerConfig{
			Host: *metricsHost,
			Port: *metricsPort,
		},
	}
}

// Addr returns the server configuration information using the format "${HOST}:${PORT}" .
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
