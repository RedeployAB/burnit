package log

import (
	"log/slog"
	"os"
)

// Logger wraps slog.Logger to provide a structured logger.
type Logger struct {
	stderr *slog.Logger
	stdout *slog.Logger
}

// New returns a new Logger.
func New() *Logger {
	return &Logger{
		stderr: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
		stdout: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// Debug logs at [LevelDebug].
func (l Logger) Debug(msg string, args ...any) {
	l.stdout.Debug(msg, args...)
}

// Info logs at [LevelInfo].
func (l Logger) Error(msg string, args ...any) {
	l.stderr.Error(msg, args...)
}

// Info logs at [LevelInfo].
func (l Logger) Info(msg string, args ...any) {
	l.stdout.Info(msg, args...)
}

// Warn logs at [LevelWarn].
func (l Logger) Warn(msg string, args ...any) {
	l.stderr.Warn(msg, args...)
}
