package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(isDebug bool) {
	log.Logger = log.With().Caller().Logger()

	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(logLevel)
}
