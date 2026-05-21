package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
)

func Test_requestLogger(t *testing.T) {
	logBuffer := &bytes.Buffer{}

	logger := slog.New(slog.NewTextHandler(logBuffer, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Time(slog.TimeKey, time.Date(2023, 10, 1, 12, 34, 57, 0, time.UTC))
			}
			if a.Key == "duration" {
				return slog.Duration("duration", time.Duration(0))
			}
			return a
		},
	}))

	requestLoggerMiddleware := middleware.RequestLogger(logger)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	loggedHandler := requestLoggerMiddleware(dummyHandler)

	req := httptest.NewRequest("GET", "http://localhost:8080/test?foo=bar", nil)
	req.Header.Set(middleware.RequestIDHeader, "test-request-id")

	rr := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	expectedLog := `time=2023-10-01T12:34:57.000Z level=INFO msg="Served request" method=GET path=/test status=0 duration=0s request_id=test-request-id
`
	if logBuffer.String() != expectedLog {
		t.Errorf("Expected log output:\n%s\nGot:\n%s", expectedLog, logBuffer.String())
	}
}
