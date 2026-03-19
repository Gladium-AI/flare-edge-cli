package logging

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func New(stderr io.Writer) zerolog.Logger {
	logger := zerolog.New(stderr).With().Timestamp().Logger()
	level := zerolog.InfoLevel
	if os.Getenv("FLARE_EDGE_LOG_LEVEL") == "debug" {
		level = zerolog.DebugLevel
	}
	return logger.Level(level)
}
