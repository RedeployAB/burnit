package log

import (
	"log/slog"
	"os"
)

// Logger is a logger with different levels.
type Logger interface {
	// Debug logs at [LevelDebug].
	Debug(msg string, args ...any)
	// Error logs at [LevelError].
	Error(msg string, args ...any)
	// Info logs at [LevelInfo].
	Info(msg string, args ...any)
	// Warn logs at [LevelWarn].
	Warn(msg string, args ...any)
}

// Logger wraps slog.Logger to provide a structured logger.
type logger struct {
	stderr *slog.Logger
	stdout *slog.Logger
}

// New returns a new Logger.
func New() *logger {
	return &logger{
		stderr: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
		stdout: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// Debug logs at [LevelDebug].
func (l logger) Debug(msg string, args ...any) {
	l.stdout.Debug(msg, args...)
}

// Error logs at [LevelError].
func (l logger) Error(msg string, args ...any) {
	l.stderr.Error(msg, args...)
}

// Info logs at [LevelInfo].
func (l logger) Info(msg string, args ...any) {
	l.stdout.Info(msg, args...)
}

// Warn logs at [LevelWarn].
func (l logger) Warn(msg string, args ...any) {
	l.stderr.Warn(msg, args...)
}
