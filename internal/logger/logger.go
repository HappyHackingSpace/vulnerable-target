// Package logger provides logging functionality for the vulnerable target application.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds logger configuration options.
type Config struct {
	Level      string
	Output     io.Writer
	TimeFormat string
	NoColor    bool
}

// DefaultConfig returns the default logger configuration.
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Output:     os.Stdout,
		TimeFormat: time.TimeOnly,
		NoColor:    false,
	}
}

// New creates a new zerolog.Logger with the given configuration.
// This is a pure factory function that returns a configured logger
// without modifying any global state.
func New(cfg *Config) zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	output := zerolog.ConsoleWriter{
		Out:        cfg.Output,
		TimeFormat: cfg.TimeFormat,
		NoColor:    cfg.NoColor,
	}

	return zerolog.
		New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}

// NewWithLevel creates a new zerolog.Logger with the specified level
// and default settings for other options.
func NewWithLevel(level string) zerolog.Logger {
	cfg := DefaultConfig()
	cfg.Level = level
	return New(cfg)
}

// SetGlobal sets the global log.Logger to the provided logger.
// This is the only function that modifies global state, and should
// be called once at application startup.
func SetGlobal(logger zerolog.Logger) {
	log.Logger = logger
}

// InitWithLevel initializes the global logger with the specified verbosity level.
// Deprecated: Use New() or NewWithLevel() with SetGlobal() instead.
func InitWithLevel(verbosityLevel string) {
	logger := NewWithLevel(verbosityLevel)
	SetGlobal(logger)
}

// Init initializes the global logger with the default info level.
// Deprecated: Use New() or NewWithLevel() with SetGlobal() instead.
func Init() {
	InitWithLevel("info")
}
