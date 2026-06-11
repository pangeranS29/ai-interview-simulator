package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

func Init() {
	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	
	// Check if in development mode
	if os.Getenv("ENV") == "development" {
		// Pretty print for development
		Log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		// JSON format for production
		Log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}