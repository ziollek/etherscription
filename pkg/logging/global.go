package logging

import (
	"os"

	"github.com/rs/zerolog"
)

var defaultLogger zerolog.Logger

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	defaultLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
}

func Logger() *zerolog.Logger {
	return &defaultLogger
}
