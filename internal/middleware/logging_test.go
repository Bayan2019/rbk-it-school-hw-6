package middleware_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLogger_AddsRequestIDAndLogsRequest(t *testing.T) {

	core, observedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := chi.NewRouter()
	router.Use(middleware.RequestLogger(logger))

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		// requestID, exists := c.Get("request_id")
		requestID, ok := r.Context().Value(middleware.RequestIDKey).(string)
		if !ok {
			requestID = "unknown" // Fallback if it wasn't set
		}

		require.True(t, ok)
		assert.Equal(t, "test-request-id", requestID)

		// c.JSON(http.StatusOK, gin.H{"message": "pong"})
		w.WriteHeader(http.StatusOK)
		dat, err := json.Marshal(map[string]string{"message": "pong"})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
		}
		_, err = w.Write(dat)
		if err != nil {
			log.Printf("Error setting Message is set: %s", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(middleware.RequestIDHeader, "test-request-id")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	logs := observedLogs.All()
	require.Len(t, logs, 1)

	assert.Equal(t, "http request completed", logs[0].Message)

	fields := logs[0].ContextMap()
	assert.Equal(t, "GET", fields["method"])
	assert.Equal(t, "/ping", fields["path"])
	assert.Equal(t, int64(200), fields["status"])
	assert.Equal(t, "test-request-id", fields["request_id"])
}

func TestRequestLogger_LogsWarnForClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, observedLogs := observer.New(zap.WarnLevel)
	logger := zap.New(core)

	// router := gin.New()

	router := chi.NewRouter()
	router.Use(middleware.RequestLogger(logger))

	router.Get("/bad-request", func(w http.ResponseWriter, r *http.Request) {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		w.WriteHeader(http.StatusBadRequest)
		dat, err := json.Marshal(map[string]string{"error": "bad request"})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
		}
		_, err = w.Write(dat)
		if err != nil {
			log.Printf("Error setting Message is set: %s", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/bad-request", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	logs := observedLogs.All()
	require.Len(t, logs, 1)

	assert.Equal(t, zap.WarnLevel, logs[0].Level)
}
