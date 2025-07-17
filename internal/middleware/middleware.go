package middleware

import (
	"context"
	"github.com/Kofandr/API_Proxy.2/internal/logger"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

func LoggerMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := uuid.New().String()

		requestLogger := log.With("requestID", requestID)

		ctx := context.WithValue(r.Context(), "requestID", requestID)
		ctx = context.WithValue(ctx, logger.CtxLoggerKey{}, requestLogger)
		r = r.WithContext(ctx)

		requestLogger.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
		)

		rw := newResponseWriter(w)

		start := time.Now()
		next.ServeHTTP(rw, r)

		requestLogger.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}
