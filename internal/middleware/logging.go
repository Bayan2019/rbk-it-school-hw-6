package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const RequestIDHeader = "X-Request-ID"

func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()

			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = "unknown"
			}

			// c.Set("request_id", requestID)
			w.Header().Set("request_id", requestID)

			// c.Next()
			next.ServeHTTP(w, r)

			duration := time.Since(startedAt)
			// statusCode := c.Writer.Status()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			statusCode := ww.Status()

			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
				zap.String("request_id", requestID),
			}

			if len(GetErrors(r).Errors) > 0 {
				fields = append(fields, zap.String("gin_errors", GetErrors(r).String()))
			}

			if statusCode >= 500 {
				logger.Error("http request completed", fields...)
				return
			}

			if statusCode >= 400 {
				logger.Warn("http request completed", fields...)
				return
			}

			logger.Info("http request completed", fields...)
		})
	}
}
