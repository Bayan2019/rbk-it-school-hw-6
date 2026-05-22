package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
)

// Ch 2. Logging Lv 4. Global Logger vs. Dependency Injection
// Add a logger field to the server struct,
// and update server logging to use that injected logger.
type Server struct {
	httpServer *http.Server
	cancel     context.CancelFunc
	logger     *slog.Logger
}

func NewServer(
	// store store.Store,
	port int,
	cancel context.CancelFunc,
	accessLogger *slog.Logger,
) *Server {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middleware.RequestLogger(accessLogger)(mux),
	}

	s := &Server{
		httpServer: srv,
		// store:      store,
		cancel: cancel,
		logger: accessLogger,
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		dat, err := json.Marshal(map[string]string{"status": "ok"})
		if err != nil {
			s.logger.Error(fmt.Sprintf("Error marshalling JSON: %s", err))
			w.WriteHeader(500)
			return
		}
		_, err = w.Write(dat)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Error setting message: %s", err))
		}
	})
	// mux.Handle("POST /api/login", s.authMiddleware(http.HandlerFunc(s.handlerLogin)))
	// mux.Handle("POST /api/shorten", s.authMiddleware(http.HandlerFunc(s.handlerShortenLink)))
	// mux.Handle("GET /api/stats", s.authMiddleware(http.HandlerFunc(s.handlerStats)))
	// mux.Handle("GET /api/urls", s.authMiddleware(http.HandlerFunc(s.handlerListURLs)))
	// mux.HandleFunc("GET /{shortCode}", s.handlerRedirect)
	mux.HandleFunc("POST /admin/shutdown", s.HandlerShutdown)

	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	// Ch 1. Observability Lv 3. What Is Observability?
	// When the server starts, print the following message to the console,
	// where %d is the port number:
	// ln.Addr() returns a net.Addr interface.
	if addr, ok := ln.Addr().(*net.TCPAddr); ok {
		httpPort := addr.Port
		s.logger.Info(fmt.Sprintf("WeatherAPI is running on http://localhost:%d\n", httpPort))
	}

	if err := s.httpServer.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) HandlerShutdown(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENV") == "production" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	go s.cancel()
}
