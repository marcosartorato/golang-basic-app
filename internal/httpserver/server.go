package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	cfg "github.com/marcosartorato/myapp/internal/config"
	metrics "github.com/marcosartorato/myapp/internal/metricsserver"
	"go.uber.org/zap"
)

type ctxKey int

const loggerKey ctxKey = iota

// withRequestLogger attaches a request-scoped logger to the context
func withRequestLogger(base *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLog := base.With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.Time("ts", time.Now()),
		)
		ctx := context.WithValue(r.Context(), loggerKey, reqLog)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getLogger fetches the request-scoped logger, falling back to a no-op logger
func getLogger(r *http.Request) *zap.Logger {
	l, _ := r.Context().Value(loggerKey).(*zap.Logger)
	if l == nil {
		l = zap.NewNop()
	}
	return l
}

// Start run the HTTP on dedicated goroutine.
func CreateServer(logger *zap.Logger, opt cfg.Options) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/hello", metrics.Instrument(withRequestLogger(logger, http.HandlerFunc(HelloHandler))))
	mux.Handle("/api/message", metrics.Instrument(withRequestLogger(logger, http.HandlerFunc(MessageHandler))))

	addr := net.JoinHostPort(*opt.Host, *opt.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server
}

// Start run the HTTP server on dedicated goroutine and return the shutdown function.
func RunServerWithShutdown(logger *zap.Logger, opts ...cfg.Option) func(context.Context) error {
	var options cfg.Options
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			panic(err)
		}
	}
	srv := CreateServer(logger, options)

	go func() {
		addr := srv.Addr
		logger.Info("App server listening on " + addr)
		if err := http.ListenAndServe(addr, srv.Handler); err != nil {
			logger.Error("app server failed: %v", zap.Error(err))
		}
	}()

	return srv.Shutdown
}
