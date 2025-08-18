package metrics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marcosartorato/myapp/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define your custom metrics here
var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

// Init registers all metrics
func Init() {
	prometheus.MustRegister(RequestsTotal, RequestDuration)
}

// Start metric server
func Start(cfg *config.ServerConfig) {
	metricsMux := http.NewServeMux()

	metricsMux.Handle("/metrics", promhttp.Handler())

	// Run metrics server in a goroutine
	go func() {
		addr := cfg.Addr()
		fmt.Println("Metrics server listening on  " + addr)
		if err := http.ListenAndServe(addr, metricsMux); err != nil {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()
}
