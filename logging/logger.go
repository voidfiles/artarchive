package logging

import (
	"io"
	stdliblog "log"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLogger returns a constructed logger.
func NewLogger(verbose bool, out io.Writer) zerolog.Logger {
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Setup a logger.
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get hostname")
	}

	if out == nil {
		out = os.Stdout
	}

	logger := zerolog.New(out).With().
		Timestamp().
		Str("host", hostname).
		Logger()

	// Tell stdlib tools that use the default logger to use our logger.
	stdliblog.SetFlags(0)
	stdliblog.SetOutput(logger)

	return logger
}
