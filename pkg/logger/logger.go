package logger

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// func NewProductionLogger() (*zap.Logger, error) {
// 	return zap.NewProduction()
// }

// func NewDevelopmentLogger() (*zap.Logger, error) {
// 	return zap.NewDevelopment()
// }

type closeFunc func() error

func InitializeLogger(logFile string) (*slog.Logger, closeFunc, error) {
	// Ch 3. Structured Logging Lv 3. Log Levels
	handlers := []slog.Handler{
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelError,
		}),
	}
	// Ch 3. Structured Logging Lv 3. Log Levels
	closers := []closeFunc{}

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}
		// defer file.Close()
		// Ch 2. Logging Lv 7. Buffered Logging
		// wrap the file writer with bufio.NewWriterSize using an 8192 byte buffer.
		bufferedFile := bufio.NewWriterSize(file, 8192)
		// Ch 2. Logging Lv 8. Logger Cleanup
		// As you create your logger,
		// also create a "close" function
		// that cleans up any logger resources.
		close := func() error {
			// close function should .Flush the buffered writer
			if err := bufferedFile.Flush(); err != nil {
				return fmt.Errorf("failed to flush log file: %w", err)
			}
			// and .Close the file.
			if err := file.Close(); err != nil {
				return fmt.Errorf("failed to close log file: %w", err)
			}
			return nil
		}
		closers = append(closers, close)
		multiWriter := io.MultiWriter(os.Stdout, bufferedFile)
		handlers = append(handlers, slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	} else {
		handlers = append(handlers, slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	logger := slog.New(slog.NewMultiHandler(
		handlers...,
	))
	closer := func() error {
		var errs []error
		for _, close := range closers {
			if err := close(); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}

	// Ch 3. Structured Logging Lv 3. Log Levels
	// Use slog.Handlers
	// to configure your STDERR logs
	// to include DEBUG and above,
	// and your file logs to include INFO and above.
	// Use slog.NewMultiHandler to combine handlers into one logger
	// used throughout the app.
	logger = slog.New(slog.NewMultiHandler(
		handlers...,
	))
	return logger, closer, nil
}
