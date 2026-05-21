package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/config"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/server"
	"github.com/Bayan2019/rbk-it-school-hw-6/pkg/logger"
)

func main() {
	err := config.MustLoad("")
	if err != nil {
		fmt.Printf("warning: assuming default configuration: .env unreadable: %v\n", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// httpPort := flag.Int("port", 8899, "port to listen on")
	// dataDir := flag.String("data", "./data", "directory to store data")
	// flag.Parse()

	status := run(ctx, cancel, config.Cfg.App.Port)
	cancel()
	os.Exit(status)
}

func run(
	ctx context.Context,
	cancel context.CancelFunc,
	httpPort int,
	// dataDir string,
) int {

	// Ch 2. Logging Lv 5. Logger Configuration
	// Assume that in production,
	// Linko has a LINKO_LOG_FILE environment variable set.
	// In local development and staging, it is not set.
	var initializeLoggerFile = getEnv("LOG_FILE", "")

	// Ch 2. Logging Lv 5. Logger Configuration
	// Add an initializeLogger helper.
	logger, close, err := logger.InitializeLogger(initializeLoggerFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		return 1
	}
	// Ch 2. Logging Lv 8. Logger Cleanup
	// Call the close function before Linko exits.
	// defer a wrapper that calls it
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()

	// Ch 2. Logging Lv 4. Global Logger vs. Dependency Injection
	// Use the access logger for server/request logs
	s := server.NewServer(httpPort, cancel, logger)
	var serverErr error
	go func() {
		serverErr = s.Start()
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("failed to shutdown server: %v\n", err))
		return 1
	}

	// Ch 1. Observability Lv 3. What Is Observability?
	// When the server shuts down (before it exits), print:
	// Ch 2. Logging Lv 4. Global Logger vs. Dependency Injection
	// use the standard logger for your Store and shutdown messages
	logger.Debug("Server is shutting down")
	if serverErr != nil {
		logger.Error(fmt.Sprintf("server error: %v\n", serverErr))
		return 1
	}

	return 0
}

////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
