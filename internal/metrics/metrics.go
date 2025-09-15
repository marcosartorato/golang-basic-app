package metrics

import (
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
func CreateServer(cfg *config.ServerConfig) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}
	return server

}
