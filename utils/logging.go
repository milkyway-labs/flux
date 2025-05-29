package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/milkyway-labs/chain-indexer/types"
)

func NewLoggerFromConfig(cfg *types.LoggingConfig) (zerolog.Logger, error) {
	if cfg == nil {
		return zerolog.Logger{}, fmt.Errorf("got nil config")
	}

	var writer io.Writer
	if cfg.LogFormat == "text" {
		writer = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.RFC3339
		})
	} else {
		writer = os.Stdout
	}

	// Create the new logger
	logger := log.Output(writer)

	// Configure the log level
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parse log level: %w", err)
	}

	return logger.Level(logLevel), nil
}
