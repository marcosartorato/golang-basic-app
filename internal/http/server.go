package hello

import (
	"context"
	"net/http"
	"time"

	"github.com/marcosartorato/myapp/internal/config"
	"github.com/marcosartorato/myapp/internal/metrics"
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
func CreateServer(cfg *config.ServerConfig) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/hello", metrics.Instrument(withRequestLogger(cfg.Logger, http.HandlerFunc(HelloHandler))))
	mux.Handle("/api/message", metrics.Instrument(withRequestLogger(cfg.Logger, http.HandlerFunc(MessageHandler))))

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}
	return server
}

// Start run the HTTP server on dedicated goroutine and return the shutdown function.
func RunServerWithShutdown(cfg *config.ServerConfig) func(context.Context) error {
	srv := CreateServer(cfg)

	go func() {
		addr := srv.Addr
		cfg.Logger.Info("App server listening on " + addr)
		if err := http.ListenAndServe(addr, srv.Handler); err != nil {
			cfg.Logger.Error("app server failed: %v", zap.Error(err))
		}
	}()

	return srv.Shutdown
}
