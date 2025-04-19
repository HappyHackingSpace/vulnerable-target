package logger

import (
	"os"
	"time"

	"github.com/happyhackingspace/vulnerable-target/pkg/options"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	options := options.GetOptions()

	level, err := zerolog.ParseLevel(options.VerbosityLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.TimeOnly,
		NoColor:    false,
	}

	log.Logger = zerolog.
		New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}
