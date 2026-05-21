package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-6/pkg/spy"
)

const RequestIDHeader = "X-Request-ID"

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = "unknown"
			}
			startedAt := time.Now()
			spyReader := &spy.SpyReadCloser{ReadCloser: r.Body}
			r.Body = spyReader
			spyWriter := &spy.SpyResponseWriter{ResponseWriter: w}

			next.ServeHTTP(spyWriter, r)
			// next.ServeHTTP(w, r)
			// logger.Info(fmt.Sprintf("Served request: %s %s", r.Method, r.URL.Path))
			// Ch 3. Structured Logging Lv 5. Key-Value Pairs
			duration := time.Since(startedAt)
			logger.Info(
				"Served request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", spyWriter.StatusCode,
				"duration", duration,
				"request_id", requestID,
			)
		})
	}
}

// func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			startedAt := time.Now()

// 			requestID := r.Header.Get(RequestIDHeader)
// 			if requestID == "" {
// 				requestID = "unknown"
// 			}

// 			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

// 			// c.Set("request_id", requestID)
// 			w.Header().Set("request_id", requestID)

// 			// c.Next()
// 			// next.ServeHTTP(ww, r)

// 			duration := time.Since(startedAt)
// 			// statusCode := c.Writer.Status()
// 			statusCode := ww.Status()

// 			fields := []zap.Field{
// 				zap.String("method", r.Method),
// 				zap.String("path", r.URL.Path),
// 				zap.Int("status", statusCode),
// 				zap.Duration("duration", duration),
// 				zap.String("request_id", requestID),
// 			}

// 			if GetErrors(r) != nil {
// 				if len(GetErrors(r).Errors) > 0 {
// 					fields = append(fields, zap.String("gin_errors", GetErrors(r).String()))
// 				}
// 			}

// 			if statusCode >= 500 {
// 				logger.Error("http request completed", fields...)
// 				return
// 			}

// 			if statusCode >= 400 {
// 				logger.Warn("http request completed", fields...)
// 				return
// 			}

// 			logger.Info("http request completed", fields...)

// 			// c.Next()
// 			next.ServeHTTP(ww, r)
// 		})
// 	}
// }
