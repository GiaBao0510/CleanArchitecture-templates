package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (l *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		l.logger.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Int("status", ww.Status()),
			zap.Int("bytes_written", ww.BytesWritten()),
			zap.String("request_id", middleware.GetReqID(r.Context())),
			zap.Duration("duration", time.Since(start)),
		)
	})

}
