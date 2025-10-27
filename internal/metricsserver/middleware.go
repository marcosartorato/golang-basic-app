package metricsserver

import (
	"net/http"
	"strconv"
	"time"
)

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusCapturingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusCapturingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// Instrument wraps an http.Handler and records RED + extras.
func Instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.Path
		method := r.Method

		// Inflight
		InflightRequests.WithLabelValues(route).Inc()
		defer InflightRequests.WithLabelValues(route).Dec()

		// Request size (approx: ContentLength; could read body for exact size)
		if r.ContentLength > 0 {
			RequestSize.WithLabelValues(method, route).Observe(float64(r.ContentLength))
		}

		// Recoverer to count panics and avoid crashing the server
		defer func(start time.Time) {
			if rec := recover(); rec != nil {
				PanicsTotal.WithLabelValues(route).Inc()
				// Return 500 if nothing was written
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}(time.Now())

		start := time.Now()
		srw := &statusCapturingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(srw, r)

		status := strconv.Itoa(srw.status)
		elapsed := time.Since(start).Seconds()

		RequestsTotal.WithLabelValues(method, route, status).Inc()
		RequestDuration.WithLabelValues(method, route, status).Observe(elapsed)
		ResponseSize.WithLabelValues(method, route, status).Observe(float64(srw.bytes))
	})
}
