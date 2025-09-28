package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Define your custom metrics here
var (
	reg = prometheus.NewRegistry()

	// RED + extras metrics
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			// Buckets tuned for typical API latencies; adjust for your service
			// 1ms .. ~16s : (-inf, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384]
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"method", "route", "status"},
	)

	InflightRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "http",
			Name:      "inflight_requests",
			Help:      "Number of inflight HTTP requests.",
		},
		[]string{"route"},
	)

	RequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "http",
			Name:      "request_size_bytes",
			Help:      "Approximate size of HTTP request in bytes.",
			Buckets:   prometheus.ExponentialBuckets(200, 2, 12), // 200B .. ~400KB (see RequestDuration for reasoning)
		},
		[]string{"method", "route"},
	)

	ResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "Size of HTTP responses in bytes.",
			Buckets:   prometheus.ExponentialBuckets(200, 2, 12),
		},
		[]string{"method", "route", "status"},
	)

	PanicsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "panics_total",
			Help:      "Number of panics recovered in HTTP handlers.",
		},
		[]string{"route"},
	)
)

// init registers all metrics
func init() {
	// App metrics, i.e. custom metrics defined above
	reg.MustRegister(RequestsTotal, RequestDuration, InflightRequests, RequestSize, ResponseSize, PanicsTotal)
	// Go/process runtime metrics (SRE staple)
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
}
