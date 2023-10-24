package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/dissipative/opabinia/internal/api/response"
)

type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	IsDebug() bool
}

type LoggerMiddleware struct {
	logger Logger
}

func NewLoggerMiddleware(logger Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: logger}
}

func (lm *LoggerMiddleware) DebugLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		defer func() {
			if errMsg := w.Header().Get(response.ErrHeader); errMsg != "" {
				lm.logger.ErrorContext(
					r.Context(),
					errMsg,
					slog.String("Method", r.Method),
					slog.String("URI", r.RequestURI),
					slog.String("User-Agent", r.Header.Get("User-Agent")),
					slog.Int64("Duration, ms", time.Since(t1).Milliseconds()),
				)
			} else if lm.logger.IsDebug() {
				lm.logger.DebugContext(
					r.Context(),
					"Request info",
					slog.String("Method", r.Method),
					slog.String("URI", r.RequestURI),
					slog.String("User-Agent", r.Header.Get("User-Agent")),
					slog.Int64("Duration, ms", time.Since(t1).Milliseconds()),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
